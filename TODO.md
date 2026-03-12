# dumpstore — Code Review TODO

## Critical

- [ ] **`runner.go`: Unbounded scanner buffer** — Add `scanner.Buffer(buf, maxSize)` to cap NDJSON output reads; unbounded reads can exhaust memory on runaway Ansible output.
- [ ] **Passwords in plaintext without TLS enforcement** — No HTTPS enforcement at server layer. Document or enforce TLS; passwords for user/SMB endpoints are exposed without it.
- [ ] **No rate limiting on password endpoints** — `/api/users` password operations have no throttling, enabling brute-force attacks.

## High

- [ ] **`handlers.go`: Dataset existence not checked before playbook** — delete, props, ACL handlers validate name format but don't verify dataset exists; playbook fails with cryptic errors.
- [ ] **`handlers.go`: ACL principal regex too loose** — `aclSafeRe` allows ambiguous NFSv4 principals (e.g. multiple `@`). Tighten to enforce strict `user@domain` format.
- [ ] **`acl.go`: `DatasetHasACL` silently returns false when tools missing** — If `getfacl`/`nfs4_getfacl` aren't installed, UI shows "no ACLs" with no indication of the real problem. Surface an error instead.
- [ ] **`playbooks/user_create.yml`: No rollback on partial failure** — If group creation succeeds but user creation fails, group is orphaned. Add cleanup/rollback logic.
- [ ] **`handlers.go`: Group members not verified to exist** — `modifyGroup` accepts comma-separated member list without checking usernames exist in `/etc/passwd`.

## Medium

- [ ] **`handlers.go`: JSON decoder accepts trailing garbage** — `{"key":"val"}junk` is silently accepted. Check for `io.EOF` after decode or use a stricter decoder.
- [ ] **`handlers.go`: `reZFSName` allows numeric-only path components** — e.g. `pool/123` passes validation but ZFS rejects it. Require first char of each component to be `[a-zA-Z]`.
- [ ] **`zfs.go`: `GetMountpointOwnership` symlink behavior undocumented** — `stat` follows symlinks silently; document or pick explicit behavior (`-L` vs. not).
- [ ] **`handlers.go`: Insufficient error context on playbook failure** — When playbook exits non-zero with no task-level failure detected, error doesn't name the last-executed task.

## Low

- [ ] **No audit logging** — User/group/dataset mutations aren't logged with any operator identity. Compliance gap.
- [ ] **ACL remove doesn't verify entry exists first** — Calling `setfacl -x` on a nonexistent entry fails silently. Pre-check with `GetDatasetACL()`.
- [ ] **`handlers.go`: Confusing error messages for missing datasets** — Many handlers let the playbook surface the error; a pre-existence check would give cleaner UX.
