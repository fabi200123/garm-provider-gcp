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
project_id = "garm-testing"
zone = "europe-west1-d"
network_id = "https://www.googleapis.com/compute/v1/projects/garm-testing/global/networks/garm"
subnetwork_id = "projects/garm-testing/regions/europe-west1/subnetworks/garm"
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
    --image  projects/debian-cloud/global/images/debian-11-bullseye-v20240110\
    --min-idle-runners 0 \
    --repo eb3f78b6-d667-4717-97c4-7aa1f3852138 \
    --tags gcp,linux \
    --provider-name gcp
```

Always find a recent image to use.

## Tweaking the provider

Garm supports sending opaque json encoded configs to the IaaS providers it hooks into. This allows the providers to implement some very provider specific functionality that doesn't necessarily translate well to other providers. Features that may exists on GCP, may not exist on Azure or AWS and vice versa.

To this end, this provider supports the following extra specs schema:

```bash
{
    "$schema": "http://cloudbase.it/garm-provider-aws/schemas/extra_specs#",
    "type": "object",
    "description": "Schema defining supported extra specs for the Garm AWS Provider",
    "properties": {
        "disksize": {
            "type": "integer"
            "description": "The size of the root disk in GB. Default is 127 GB."
        }
        "network_id": {
            "type": "string"
            "description": "The name of the network attached to the instance."
        }
        "subnet_id": {
            "type": "string"
            "description": "The name of the subnetwork attached to the instance."
        }
        "nic_type": {
            "type": "string"
            "description": "The type of NIC attached to the instance. Default is VIRTIO_NET."
        }
    }
}
```

An example of extra specs json would look like this:

```bash
{
    "disksize": 255,
    "network_id": "https://www.googleapis.com/compute/v1/projects/garm-testing/global/networks/garm",
    "subnet_id": "projects/garm-testing/regions/europe-west1/subnetworks/garm",
    "nic_type": "VIRTIO_NET"
}
```

To set it on an existing pool, simply run:

```bash
garm-cli pool update --extra-specs='{"disksize" : 100}' <POOL_ID>
```

You can also set a spec when creating a new pool, using the same flag.

Workers in that pool will be created taking into account the specs you set on the pool.