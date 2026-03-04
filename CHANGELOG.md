# Changelog

All notable changes to this project will be documented here.

## [v0.0.3] — 2026-03-04

### Changed
- Replace placeholder text logo with SVG lockup in the UI header
- Add SVG favicon (`dumpstore-blue-dark-icon48.svg`)
- Update README to use dark/light-mode-aware SVG logos via `<picture>`
- Add `images/` directory with full set of logo variants (blue/mono, dark/light, icon48/icon80/lockup)

## [v0.0.2] — 2026-02-xx

### Added
- **Live updates via SSE** — Server-Sent Events endpoint (`GET /api/events`) pushes pool, dataset, snapshot, and I/O changes every 10 s; browser falls back to 30 s REST polling if the connection is lost
- Subscription broker (`internal/broker`) with per-topic pub/sub and change detection (JSON equality check)
- Background ZFS poller goroutine that publishes only on data change
- Dark/light mode logo variants in README
- Screenshots in README

## [v0.0.1] — 2026-01-xx

### Added
- Initial release
- Go HTTP server (stdlib only, no external dependencies)
- **System info** — hostname, OS, kernel, CPU, uptime, load averages, process stats
- **Pool overview** — health badges, usage bars, fragmentation, deduplication ratio, vdev tree
- **I/O statistics** — live read/write IOPS and bandwidth per pool
- **Disk health** — S.M.A.R.T. data per drive via `smartctl`
- **Dataset browser** — collapsible tree, compression, quota, mountpoint
- **Dataset management** — create, edit properties, and delete (with confirm-by-typing dialog) filesystems and volumes
- **Snapshot management** — list, create (recursive), and delete snapshots
- Ansible playbook runner for write operations with structured JSON output
- Prometheus metrics endpoint (`GET /metrics`)
- systemd unit file (Linux) and rc.d script (FreeBSD)
- `make install` with OS-aware service registration
