# Digest Auth Removal Proxy

This project provides a reverse proxy that transparently removes HTTP Digest Authentication from backend servers, allowing clients to interact without needing authentication. It is especially useful for cases like the Prusa Core One's PrusaLink server, which enforces Digest Auth and does not allow disabling or changing authentication methods ([issue #3559](https://github.com/prusa3d/Prusa-Firmware-Buddy/issues/3559)). This makes integration with standard reverse proxies or SSO solutions difficult. The proxy handles Digest Auth itself and forwards requests to the backend, so clients do not need to be aware of Digest Auth, making it possible to integrate with reverse proxies like Nginx.

**Security Warning:**
This proxy removes authentication from the backend. You should use it only in trusted environments or in conjunction with another authentication layer (such with auth request based SSO on nginx). Running this proxy without additional authentication exposes your backend to anyone who can reach the proxy.

## Running with GitHub Container Registry (GHCR)

A pre-built container image is available at [ghcr.io/shuhaowu/digest-auth-removal-proxy:v1](https://github.com/shuhaowu/digest-auth-removal-proxy/pkgs/container/digest-auth-removal-proxy). It works the same as the locally built Docker image documented below.

### Example

```sh
docker pull ghcr.io/shuhaowu/digest-auth-removal-proxy:v1

docker run \
  -e DIGEST_AUTH_REMOVAL_PROXY_USERNAME=<username> \
  -e DIGEST_AUTH_REMOVAL_PROXY_PASSWORD=<password> \
  -e DIGEST_AUTH_REMOVAL_PROXY_BACKEND=http://prusalink-coreone \
  -p 8080:8080 \
  ghcr.io/shuhaowu/digest-auth-removal-proxy:v1
```

## Building and Running (Bare Metal)

### Build

```sh
go build -o digest-auth-removal-proxy .
```

### Run

```sh
./digest-auth-removal-proxy \
  --username <username> \
  --password <password> \
  --backend http://prusalink-core-one \
```

Other flags:

- `--username`: Digest Auth username for backend
- `--password`: Digest Auth password for backend
- `--backend`: Backend URL (e.g., `http://prusalink-coreone:80`)
- `--listen-host`: Host to listen on (default: all interfaces)
- `--listen-port`: Port to listen on (default: 8080)
- `--debug`: Enable debug logging (optional)

Then go to https://localhost:8080 and you should see it working.

## Running with Docker

A minimal Docker image is provided. The entrypoint script converts environment variables to CLI flags.

### Example

```sh
docker run --rm \
  -e DIGEST_AUTH_REMOVAL_PROXY_USERNAME=<username> \
  -e DIGEST_AUTH_REMOVAL_PROXY_PASSWORD=<password> \
  -e DIGEST_AUTH_REMOVAL_PROXY_BACKEND=http://prusalink-coreone \
  -e DIGEST_AUTH_REMOVAL_PROXY_LISTEN_HOST=0.0.0.0 \
  -e DIGEST_AUTH_REMOVAL_PROXY_LISTEN_PORT=8080 \
  -e DIGEST_AUTH_REMOVAL_PROXY_DEBUG=1 \
  -p 8080:8080 \
  digest-auth-removal-proxy
```

Then go to https://localhost:8080 and you should see it working.

- All CLI flags can be set via environment variables:
  - `DIGEST_AUTH_REMOVAL_PROXY_USERNAME`
  - `DIGEST_AUTH_REMOVAL_PROXY_PASSWORD`
  - `DIGEST_AUTH_REMOVAL_PROXY_BACKEND`
  - `DIGEST_AUTH_REMOVAL_PROXY_LISTEN_HOST`
  - `DIGEST_AUTH_REMOVAL_PROXY_LISTEN_PORT`
  - `DIGEST_AUTH_REMOVAL_PROXY_DEBUG` (set to `1` for debug logging)
