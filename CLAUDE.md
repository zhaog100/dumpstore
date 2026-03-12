# dumpstore — Claude context

## Workflow Orchestration
### 1. Plan Node Default
- Enter plan mode for ANY non-trivial task (3+ steps or architectural decisions)
- If something goes sideways, STOP and re-plan immediately - don't keep pushing
- Use plan mode for verification steps, not just building
- Write detailed specs upfront to reduce ambiguity
### 2. Subagent Strategy
- Use subagents liberally to keep main contect window clean
- Offload research, exploration, and parallel analysis to subagents
- For complex problens, throw more compute at it via subagents
- One tack per subagent for focused execution
### 3. Self-Improvement Loop
- After ANY correction from the user: update tasks/lessons.md with the pattern
- Write rules for yourself that prevent the same mistake
- Ruthlessly iterate on these lessons until mistake rate drops
- Review lessons at session start for relevant project
### 4. Verification Before Done
- Never mark a task complete without proving it works
- Diff behavior between main and your changes when relevant
- Ask yourself: "Would a staff engineer approve this?"
- Run tests, check logs, demonstrate correctness
### 5. Demand Elegance (Balanced)
- For non-trivial changes: pause and ask "is there a more elegant way?"
- If a fix feels hacky: "Knowing everything I know now, implement the elegant solution"
- Skip this for simple, chvious fixes - don't over-engineer
- Challenge your own work before presenting it
### 6. Autonomous Bug Fixing
- When given a bug report: just fix it. Don't ask for hand-holding
- Point at logs, errors, failing tests - then resolve them
- Zero context switching required from the user
- Go fix failing CI tests without being told how
## Task Management
## Core Principles
- **Simplicity First**: Make every change as simple as possible. Inpact minimal code.
- **No Laziness**: Find root causes. No temporary fixes. Senior developer standards.
- **Minimal Impact**: Changes should only touch what's necessary. Avoid introducing bugs.


## Build & check

```bash
go build ./...   # must always pass before committing
go vet ./...     # must always pass
```

No external Go dependencies. Standard library only.

## Architecture rules

**Read ops → `internal/zfs` package (direct exec)**
`zpool`/`zfs` CLI commands are called directly via `exec.Command` for low latency. No Ansible for reads.

**Write ops → Ansible playbooks**
Create/destroy operations go through `ansible-playbook` with `ANSIBLE_STDOUT_CALLBACK=ndjson` (custom callback plugin at `playbooks/callback_plugins/ndjson.py`). The runner in `internal/ansible/runner.go` parses the ndjson output to extract task results.

Do not change this split without a good reason — it exists to avoid Ansible's Python startup overhead on every API call.

## Adding a new operation

1. If it's a **read**: add a function to `internal/zfs/zfs.go` and a handler in `internal/api/handlers.go`.
2. If it's a **write**: add a playbook in `playbooks/`, wire it up in `handlers.go` using `h.runner.Run(...)`.
3. Register the route in `handlers.go:RegisterRoutes`.
4. Add the UI in `static/app.js` (render function + fetch call) and `static/index.html` (button + dialog if needed).

## Playbook conventions

- All playbooks target `localhost` with `gather_facts: false`.
- Required extra vars must be documented in a header comment.
- Always include an `assert` task before any mutating command.
- Task names must be stable — `RunAndGetStdout` looks them up by name.

## Frontend

- Vanilla JS, no framework, no build step.
- All data lives in the `state` object.
- Render functions are pure: they read from `state` and write innerHTML.
- Always escape user-controlled strings through `esc()` before inserting into HTML.
- The `api()` helper throws on non-2xx responses with the server's `error` field.

## Security

- Input to Ansible extra-vars is checked for shell-special characters (`@;|&$\``) in handlers before the playbook call.
- The `static/` directory is served with `http.FileServer` — do not put secrets there.
- The service runs as root (required for ZFS). Do not expose it on a public interface without authentication in front of it.

## File map

| File | Responsibility |
|------|---------------|
| `main.go` | Server setup, flag parsing, dependency check (`ansible-playbook` in PATH, `playbooks/` and `static/` dirs exist) |
| `internal/zfs/zfs.go` | `ListPools`, `ListDatasets`, `ListSnapshots`, `IOStats` — all direct CLI calls |
| `internal/ansible/runner.go` | `Runner.Run` — executes a playbook and returns parsed `PlaybookOutput`; `RunAndGetStdout` — convenience wrapper |
| `internal/api/handlers.go` | All HTTP handlers + input validation + `writeJSON` / `writeError` helpers |
| `playbooks/zfs_dataset_create.yml` | Creates filesystem or volume; vars: `name`, `type`, `volsize`, `compression`, `quota`, `mountpoint` |
| `playbooks/zfs_snapshot_create.yml` | Creates snapshot; vars: `dataset`, `snapname`, `recursive` |
| `playbooks/zfs_snapshot_destroy.yml` | Destroys snapshot; vars: `snapshot`, `recursive` |
| `playbooks/inventory/localhost` | `ansible_connection=local`, `ansible_python_interpreter=auto_silent` |
| `static/index.html` | Page shell, dialogs (new dataset, new snapshot) |
| `static/app.js` | State, fetch, render functions for pools/datasets/snapshots/iostat, dialog wiring |
| `static/style.css` | Dark monospace theme, CSS variables in `:root` |
| `contrib/dumpstore.service` | systemd unit; binary at `/usr/local/lib/dumpstore/dumpstore` |
| `contrib/dumpstore.rc` | FreeBSD rc.d script |
| `install.sh` | Standalone install/uninstall script (wraps build + service setup) |
| `Makefile` | `build`, `install` (copies binary + playbooks + static, enables service), `uninstall`, `clean` |
