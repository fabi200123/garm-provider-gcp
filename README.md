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

The config file for this external provider is a simple toml used to configure the AWS credentials it needs to spin up virtual machines.

```bash
project_id
zone = "eu-central-1"
network_id = "sample_network_id"
CredentialsFile = "/home/ubuntu/credentials.json"
```

## Creating a pool

After you [add it to garm as an external provider](https://github.com/cloudbase/garm/blob/main/doc/providers.md#the-external-provider), you need to create a pool that uses it. Assuming you named your external provider as ```gcp``` in the garm config, the following command should create a new pool:

```bash
garm-cli pool create \
    --os-type linux \
    --os-arch amd64 \
    --enabled=true \
    --flavor e2-medium \
    --image  debian\
    --min-idle-runners 0 \
    --repo e0207029-e3bf-493a-9caa-64e80610c5ee \
    --tags gcp,linux \
    --provider-name gcp
```