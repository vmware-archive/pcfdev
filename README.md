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

## Development

> NOTE: These instructions are for people contributing code to microPCF.  If you only want to deploy microPCF, see above.
> These instructions cover Vagrant/Virtualbox development.
> A similar process can be followed for Vagrant/VMWare and Vagrant/AWS development.  More documentation is forthcoming.

To develop microPCF you will need to have the following tools installed:

- Packer
- Vagrant
- Virtualbox
- Direnv _(optional)_

### Clone the microPCF source

```bash
git clone --recursive https://github.com/pivotal-cf/micropcf.git
```

### Build microPCF

Setup your shell for building microPCF:

```bash
cd micropcf
direnv allow # or: source .envrc
```

#### Building a microPCF Box

If you change any Diego components, you'll need to build a local Vagrant box with your changes.  If you don't plan to change any Diego components, you can update the `box_version` property in the `Vagrantfile` to point to a [pre-built microPCF box on Atlas](https://atlas.hashicorp.com/micropcf/boxes/base) and skip this step.

```bash
./images/base/build -only=virtualbox-iso
vagrant box add --force micropcf-virtualbox-v0.box --name micropcf/base
```

### Deploy microPCF

```bash
# in micropcf/vagrant
vagrant up --provider=virtualbox
```

## Contributing

If you are interested in contributing to microPCF, please refer to [CONTRIBUTING](CONTRIBUTING.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
