#! /bin/sh

set -e
set -x

NAMESPACE=$1
if [ -z "$NAMESPACE" ]; then
    exit 1
fi

TOKEN=$2
if [ -z "$TOKEN" ]; then
    echo "No token provided: continuing without auth token"
fi

REMOTE=
# try getting the remote address for 1 minute
for _ in $(seq 1 60); do
    REMOTE=$(kubectl -n "$NAMESPACE" get ingress lepton-api-server-ingress -o jsonpath='{.metadata.annotations.external-dns\.alpha\.kubernetes\.io/hostname}')
    if [ -n "$REMOTE" ]; then
        break
    fi
    sleep 1
done
if [ -z "$REMOTE" ]; then
    exit 1
fi

REMOTE=https://$REMOTE/api/v1
# run up to 60-minute as we add more e2e tests for example models
if ! COLUMNS=2000 go test -timeout 3600s -v ./e2e-tests/... --remote-url "$REMOTE" --auth-token "$TOKEN"; then
    exit 1
fi
