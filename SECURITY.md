# Security Notes

dumpstore is designed for **trusted, private networks** (home lab, local LAN). It has no built-in authentication and runs as root. The notes below describe known risks and the recommended mitigations.

---

## TLS / plaintext passwords

dumpstore does **not** enforce HTTPS at the application layer. Several API endpoints accept passwords in the request body:

- `POST /api/users` — Unix user password (`password` field)
- `POST /api/users/{name}` — Unix user password change
- `POST /api/smb/users` — Samba password (`smb_password` field)

Without TLS these credentials travel in plaintext over the network.

**Recommended mitigations (pick one):**

1. **Reverse proxy with TLS termination** — put dumpstore behind nginx, Caddy, or Traefik and only expose the proxy over HTTPS. The proxy handles certificate management; dumpstore stays on `127.0.0.1`.
2. **SSH tunnel** — `ssh -L 8080:localhost:8080 nas-host` and access via `http://localhost:8080`. No certificate required.
3. **VPN** — restrict access to a WireGuard or OpenVPN network you already trust.

A future release may add a built-in TLS flag (`-tls-cert` / `-tls-key`) to terminate HTTPS directly, but that is not currently implemented.

---

## No rate limiting on password endpoints

The password-accepting endpoints listed above have no rate limiting or lockout logic. An attacker with network access can attempt passwords at line rate.

**Recommended mitigations:**

- Apply the network controls above so the service is not reachable from untrusted hosts.
- If exposing via a reverse proxy, configure the proxy's rate-limiting module (e.g. `limit_req` in nginx, `rate_limit` in Caddy) on the relevant paths.

---

## General advice

- Bind to `127.0.0.1` (the default) and access via a proxy or tunnel rather than binding to `0.0.0.0`.
- Firewall the port at the OS level if a proxy is not in use.
- The service runs as root (required for ZFS). Treat it with the same access controls you would apply to a root shell.
