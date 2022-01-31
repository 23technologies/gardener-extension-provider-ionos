############# builder
FROM eu.gcr.io/gardener-project/3rd/golang:1.17.5 AS builder

ENV BINARY_PATH=/go/bin
WORKDIR /go/src/github.com/23technologies/gardener-extension-provider-ionos

COPY . .
RUN make build

############# base
FROM eu.gcr.io/gardener-project/3rd/alpine:3.13 as base

############# gardener-extension-provider-ionos
FROM base AS gardener-extension-provider-ionos
LABEL org.opencontainers.image.source="https://github.com/23technologies/gardener-extension-provider-ionos"

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-provider-ionos /gardener-extension-provider-ionos
COPY --from=builder /go/bin/gardener-extension-admission-ionos /gardener-extension-admission-ionos
ENTRYPOINT ["/gardener-extension-provider-ionos"]
