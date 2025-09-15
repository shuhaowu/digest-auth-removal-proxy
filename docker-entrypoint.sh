#!/bin/bash
set -e

if [ $# -gt 0 ]; then
    exec /usr/local/bin/digest-auth-removal-proxy "$@"
else
    ARGS=()

    [ -n "$DIGEST_AUTH_REMOVAL_PROXY_USERNAME" ] && ARGS+=("--username" "$DIGEST_AUTH_REMOVAL_PROXY_USERNAME")
    [ -n "$DIGEST_AUTH_REMOVAL_PROXY_PASSWORD" ] && ARGS+=("--password" "$DIGEST_AUTH_REMOVAL_PROXY_PASSWORD")
    [ -n "$DIGEST_AUTH_REMOVAL_PROXY_BACKEND" ] && ARGS+=("--backend" "$DIGEST_AUTH_REMOVAL_PROXY_BACKEND")
    [ -n "$DIGEST_AUTH_REMOVAL_PROXY_LISTEN_HOST" ] && ARGS+=("--listen-host" "$DIGEST_AUTH_REMOVAL_PROXY_LISTEN_HOST")
    [ -n "$DIGEST_AUTH_REMOVAL_PROXY_LISTEN_PORT" ] && ARGS+=("--listen-port" "$DIGEST_AUTH_REMOVAL_PROXY_LISTEN_PORT")
    [ "$DIGEST_AUTH_REMOVAL_PROXY_DEBUG" = "1" ] && ARGS+=("--debug")

    exec /usr/local/bin/digest-auth-removal-proxy "${ARGS[@]}"
fi
