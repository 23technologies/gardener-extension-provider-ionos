############# builder
FROM eu.gcr.io/gardener-project/3rd/golang:1.16.2 AS builder

ENV BINARY_PATH=/go/bin
WORKDIR /go/src/github.com/23technologies/gardener-extension-provider-ionos

COPY . .
RUN make build

############# base
FROM eu.gcr.io/gardener-project/3rd/alpine:3.13.2 as base

############# gardener-extension-provider-ionos
FROM base AS gardener-extension-provider-ionos
LABEL org.opencontainers.image.source="https://github.com/23technologies/gardener-extension-provider-ionos"

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-provider-ionos /gardener-extension-provider-ionos
ENTRYPOINT ["/gardener-extension-provider-ionos"]

############# gardener-extension-validator-ionos
#FROM base AS gardener-extension-validator-ionos
#
#COPY --from=builder /go/bin/gardener-extension-validator-ionos /gardener-extension-validator-ionos
#ENTRYPOINT ["/gardener-extension-validator-ionos"]
