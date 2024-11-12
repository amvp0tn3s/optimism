FROM --platform=linux/amd64 us-west1-docker.pkg.dev/blockchaintestsglobaltestnet/dev-images/op-geth@sha256:93694cd5e2e94ac6c1b60539e185eeb6584c3400e43c239c47bd26be06dd656a

RUN apk add --no-cache jq

COPY l2-op-geth-entrypoint.sh /entrypoint.sh

VOLUME ["/db"]

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
