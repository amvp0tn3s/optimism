FROM --platform=linux/amd64 us-west1-docker.pkg.dev/blockchaintestsglobaltestnet/dev-images/op-geth@sha256:c8c823b13ad01d1a528a61094c1d6fda7308013c692787692c06045d8bec9f8a

RUN apk add --no-cache jq

COPY l2-op-geth-entrypoint.sh /entrypoint.sh

VOLUME ["/db"]

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
