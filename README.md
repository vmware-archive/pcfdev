# microPCF (Pivotal Cloud Foundry)

MicroPCF is an open source project for running a local version of Pivotal Cloud Foundry.  It supports the [CF CLI](https://github.com/cloudfoundry/cli) and runs using [Vagrant](https://www.vagrantup.com/) on [VirtualBox](https://www.virtualbox.org/), [VMware Fusion for Mac](https://www.vmware.com/products/fusion), [VMware Workstation for Windows](https://www.vmware.com/products/workstation), and [Amazon Web Services](http://aws.amazon.com/).

[ [Website](http://micropcf.io) | [Latest Release](https://github.com/pivotal-cf/micropcf/releases/latest) | [Nightly Builds](https://micropcf.s3.amazonaws.com/nightly/index.html) ]

## Deploy microPCF with Vagrant

A colocated deployment of microPCF can be launched locally with [Vagrant](https://vagrantup.com/). You will need:

* A microPCF Vagrantfile from the [latest release](https://github.com/pivotal-cf/micropcf/releases/latest) or [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html)
* [Vagrant](https://vagrantup.com/) 1.7+ installed

##### Spin up a virtual environment

Download the Vagrantfile into a new local folder, and open a prompt to that folder:

```bash
# download https://github.com/pivotal-cf/micropcf/releases/download/<VERSION>/Vagrantfile-<VERSION>.base
# mv Vagrantfile-<VERSION>.base Vagrantfile
vagrant up --provider virtualbox
```

By default, the cluster can be targeted at `cf api api.local.micropcf.io --skip-ssl-validation`.

> Unless you're attempting to develop microPCF itself, please download the Vagrantfile from [Github Releases](https://github.com/pivotal-cf/micropcf/releases/latest) or [Nightly Builds](https://micropcf.s3.amazonaws.com/nightly/index.html).  There's no need to clone the repository in order to use microPCF.

##### Supported environment variables

These variables must be set during vagrant up.

1. `MICROPCF_IP` - sets the IP address to bring up the VM on
1. `MICROPCF_DOMAIN` - sets an alternate alias for the system routes to be defined on
  - defaults to `local.micropcf.io`, then `$MICROPCF_IP.xip.io`
1. `VM_CORES` - number of CPU cores to allocate on the Guest VM (defaults to host # of logical CPUs)
1. `VM_MEMORY` - number of MB to allocate on the Guest VM (defaults to 25% of host memory)

##### Install `cf` - CF CLI

More information is available on the [Cloud Foundry CLI README](https://github.com/cloudfoundry/cli#downloads) or the [Cloud Foundry CLI Releases](https://github.com/cloudfoundry/cli/releases/latest) page.  Please install the appropriate binary for your architecture.

## Troubleshooting

1. Ubuntu 14.04 LTS does not install a compatible version of Vagrant by default.  A compatible version can be found on the [Vagrant Downloads](http://www.vagrantup.com/downloads.html) page.
1. Use an Administrator shell to deploy using VMware Workstation on Windows.

## Contributing

If you are interested in contributing to microPCF, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
