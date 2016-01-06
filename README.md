# microPCF (Pivotal Cloud Foundry)

microPCF is an open-source project that allows developers to easily run a single VM version of Pivotal Cloud Foundry.  It supports the [CF CLI](https://github.com/cloudfoundry/cli) and runs using [Vagrant](https://www.vagrantup.com/) on [VirtualBox](https://www.virtualbox.org/), [VMware Fusion for Mac](https://www.vmware.com/products/fusion), [VMware Workstation for Windows](https://www.vmware.com/products/workstation), or [Amazon Web Services](http://aws.amazon.com/).

[ [Website](http://micropcf.io) | [Latest Release](https://github.com/pivotal-cf/micropcf/releases/latest) | [Nightly Builds](https://micropcf.s3.amazonaws.com/nightly/index.html) ]

## Deploy microPCF

microPCF can be deployed with [Vagrant](https://vagrantup.com/). You will need:

* A microPCF Vagrantfile from the [latest release](https://github.com/pivotal-cf/micropcf/releases/latest) or [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html)
* [Vagrant](https://vagrantup.com/) 1.7+ installed

To deploy locally, you will need one of the following virtualizers
* Virtualbox 5.0+
* Vmware Fusion 8+ (for OSX)
* VMware Desktop 11+ (for Windows)

To deploy to AWS, you will need the vagrant-aws plugin (`vagrant plugin install vagrant-aws`) and credentials for an AWS account in your environment.

> Cloning this repo is not required to use microPCF.

##### Spin up a virtual environment

Download the Vagrantfile for your intended release version into a new local folder and spin up a virtual environment.
Vagrantfiles for each release version are available at https://github.com/pivotal-cf/micropcf/releases.

The following example assumes v0.3.0:

```bash
mkdir WORKSPACE && cd WORKSPACE
curl -L https://github.com/pivotal-cf/micropcf/releases/download/v0.3.0/Vagrantfile-v0.3.0.base -o Vagrantfile
vagrant up --provider <PROVIDER>
```
`<PROVIDER>` should be set to one of `virtualbox`, `vmware_fusion`, `vmware_desktop`, or `aws`.


##### Supported environment variables
You may set the following environment variables when you `vagrant up` to customize microPCF:

1. `MICROPCF_IP` - sets the IP address to bring up the VM on
  - defaults to 192.168.11.11 locally
  - defaults to AWS-assigned public IP on AWS
1. `MICROPCF_DOMAIN` - sets an alternate alias for the system routes to be defined on
  - defaults to `local.micropcf.io` when deploying locally
  - defaults to <MICROPCF_IP>.xip.io on AWS
1. `VM_CORES` (local only) - number of CPU cores to allocate on the Guest VM
  - defaults to host # of logical CPUs
1. `VM_MEMORY` (local only) - number of MB to allocate on the Guest VM 
  - defaults to 25% of host memory

##### Install `cf` - CF CLI

Please install the appropriate binary for your architecture from the [Cloud Foundry CLI README](https://github.com/cloudfoundry/cli#downloads) or the [Cloud Foundry CLI Releases page](https://github.com/cloudfoundry/cli/releases/latest).


##### Use CF CLI to interact with your microPCF

By default, you can point your CLI to your microPCF deployment by running `cf api api.<MICROPCF_DOMAIN> --skip-ssl-validation`. If you're deploying microPCF locally, the default value for `<MICROPCF_DOMAIN>` is `local.micropcf.io`.

To stage a simple app on microPCF, `cd` into the app directory and run `cf push <APP_NAME>`.

See cf documentation for information on [deploying apps](http://docs.cloudfoundry.org/devguide/deploy-apps/) and [attaching services](http://docs.cloudfoundry.org/devguide/services/).

## Troubleshooting

See our [troubleshooting](TROUBLESHOOTING.md) page.

## Contributing

If you are interested in contributing to microPCF, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
