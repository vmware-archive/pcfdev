# microPCF (Pivotal Cloud Foundry)

MicroPCF is an open source project for running a local version of Pivotal Cloud Foundry.  It supports the [CF CLI](https://github.com/cloudfoundry/cli) and runs using [Vagrant](https://www.vagrantup.com/) on [VirtualBox](https://www.virtualbox.org/), [VMware Fusion for Mac](https://www.vmware.com/products/fusion), [VMware Workstation for Windows](https://www.vmware.com/products/workstation), and [Amazon Web Services](http://aws.amazon.com/).

[ [Website](http://micropcf.io) | [Latest Release](https://github.com/pivotal-cf/micropcf/releases/latest) | [Nightly Builds](https://micropcf.s3.amazonaws.com/nightly/index.html) ]

## Deploy microPCF with Vagrant

A colocated deployment of microPCF can be launched locally with [Vagrant](https://vagrantup.com/). You will need:

* A microPCF Vagrantfile from the [latest release](https://github.com/pivotal-cf/micropcf/releases/latest) or [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html)
* [Vagrant](https://vagrantup.com/) 1.7+ installed

> NOTE: Ubuntu 14.04 LTS does not install a compatible version of Vagrant by default. You can upgrade the version that you get out of the box by downloading the `.deb` file from [Vagrant](http://www.vagrantup.com/downloads.html).

##### Spin up a virtual environment

Download the Vagrantfile into a new local folder, and open a prompt to that folder:

```bash
vagrant up --provider virtualbox
```

This spins up a virtual environment that is accessible at `local.micropcf.io`

> Use an Administrator shell to deploy using VMware Workstation on Windows.

##### Supported environment variables

1. `MICROPCF_IP` - sets the IP address to bring up the VM on
1. `MICROPCF_DOMAIN` - sets an alternate alias for the system routes to be defined on
  - defaults to `local.micropcf.io`, then `$MICROPCF_IP.xip.io`
1. `VM_CORES` - number of CPU cores to allocate on the Guest VM (defaults to host # of logical CPUs)
1. `VM_MEMORY` - number of MB to allocate on the Guest VM (defaults to 25% of host memory)

##### Install `cf` - CF CLI

More information is available on the [Cloud Foundry CLI README](https://github.com/cloudfoundry/cli#downloads) or the [Cloud Foundry CLI Releases](https://github.com/cloudfoundry/cli/releases/latest) page.  Please install the appropriate binary for your architecture.

## Contributing

If you are interested in contributing to microPCF, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
