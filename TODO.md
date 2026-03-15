# dumpstore — Code Review TODO

## Critical

- [x] **`runner.go`: Unbounded scanner buffer** — Added `scanner.Buffer(buf, 4 MB)` cap on NDJSON line reads.
- [x] **Passwords in plaintext without TLS enforcement** — Documented in `SECURITY.md` with recommended mitigations (reverse proxy, SSH tunnel, VPN).
- [x] **No rate limiting on password endpoints** — Documented in `SECURITY.md`; mitigation is proxy-layer rate limiting.

## High

- [x] **`handlers.go`: Dataset existence not checked before playbook** — `datasetExists` helper added; `deleteDataset`, `setDatasetProps`, `setDatasetOwnership`, `setACLEntry`, `removeACLEntry` all return 404 before running the playbook.
- [x] **`handlers.go`: ACL principal regex too loose** — `aclSafeRe` now enforces per-field `@` rules: at most one, not leading, trailing only for all-uppercase NFSv4 well-known principals (OWNER@, GROUP@, EVERYONE@).
- [x] **`acl.go`: `DatasetHasACL` silently returns false when tools missing** — Signature changed to `(bool, error)`; both tool errors are now propagated. Call site in `getACLStatus` logs `slog.Warn` and falls back to `acltype` property.
- [ ] **`playbooks/user_create.yml`: No rollback on partial failure** — skipped; Ansible lacks native rollback and the orphaned-group risk is low enough to defer.
- [x] **`handlers.go`: Group members not verified to exist** — `modifyGroup` now calls `system.ListUsers()` and rejects unknown member names with a 400 listing the offenders.

## Medium

- [x] **`handlers.go`: JSON decoder accepts trailing garbage** — `decodeJSON` helper added; does a second `Decode` and rejects unless it returns `io.EOF`. All 14 call sites updated.
- [x] **`handlers.go`: `reZFSName` allows numeric-only path components** — Regex first-char class tightened from `[a-zA-Z0-9]` to `[a-zA-Z]` for each path component.
- [x] **`zfs.go`: `GetMountpointOwnership` symlink behavior undocumented** — Added `-L` flag explicitly and documented the follow-symlink choice in the function comment.
- [x] **`handlers.go`: Insufficient error context on playbook failure** — `lastTaskName` tracked during NDJSON scan; included in the fallback non-zero-exit error message.

## Low

- [ ] **No audit logging** — User/group/dataset mutations aren't logged with any operator identity. Compliance gap.
- [ ] **ACL remove doesn't verify entry exists first** — Calling `setfacl -x` on a nonexistent entry fails silently. Pre-check with `GetDatasetACL()`.
- [ ] **`handlers.go`: Confusing error messages for missing datasets** — Many handlers let the playbook surface the error; a pre-existence check would give cleaner UX.
