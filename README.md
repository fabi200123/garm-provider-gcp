# Garm External Provider For GCP

The GCP external provider allows [garm](https://github.com/cloudbase/garm) to create Linux and Windows runners on top of AWS virtual machines.

## Build

Clone the repo:

```bash
git clone https://github.com/cloudbase/garm-provider-gcp
```

Build the binary:

```bash
cd garm-provider-gcp
go build .
```

Copy the binary on the same system where garm is running, and [point to it in the config](https://github.com/cloudbase/garm/blob/main/doc/providers.md#the-external-provider).

## Configure
