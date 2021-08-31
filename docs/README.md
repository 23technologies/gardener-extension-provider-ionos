# [Gardener Extension for IONOS Cloud provider](https://gardener.cloud)

This controller implements Gardener's extension contract for the IONOS cloud provider.

## Controller implemented

- controlplane
- healthcheck
- infrastructure
- worker

## Supported features

- Generic controlplane actuator
- Generic healthcheck actuator
- Support for events reconcile and delete of infrastructure
- Worker actuator

### Infrastructure actions

- Supports creation of private networks in IONOS Cloud
- Adds Gardener Public Key for use in nodes

## Unsupported features

- Root volume customization (restricted to IONOS Cloud image sizes and type)
- Additional data volumes
- Mapping of Gardener Machine Profiles to IONOS Cloud image names
- Many more ... We highly appreciate any kind of patch.
