#!/usr/bin/env bash

docker rm -f vsummary-vcsim | 2>/dev/null || true
docker run -d --name vsummary-vcsim \
  -p 8989:8989 \
  -v $(pwd)/testdata/tls/:/data/tls \
  -u root \
  gbolo/vcsim vcsim \
    -l 0.0.0.0:8989 \
    -tls \
    -tlscert /data/tls/server_vcenter-simulator-chain.pem \
    -tlskey /data/tls/server_vcenter-simulator-key.pem \
    -pg 10 -dc 5 -app 0 \
    -folder 0 -ds 3 -pool 2 \
    -pod 0 -cluster 3 -vm 10
