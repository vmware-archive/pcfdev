# MicroPCF

MicroPCF is a new distribution of Cloud Foundry designed to run on a developer’s laptop or workstation.  MicroPCF gives application developers the full Cloud Foundry experience in a lightweight, easy to install package.  MicroPCF is intended for application developers who wish to develop and debug their application locally on a full-featured Cloud Foundry.  MicroPCF is also an excellent getting started environment for developers interested in learning and exploring Cloud Foundry.

> More information about the project can be found on the [FAQ](FAQ.md#general-questions).

## Install

1. Download the `micropcf-<VERSION>.zip` from [Github releases](https://github.com/pivotal-cf/micropcf/releases) or [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html).
1. Unzip the `micropcf-<VERSION>.zip`.
1. Open a terminal or command prompt and navigate to the `micropcf-<VERSION>` folder.
1. Run `vagrant up --provider=<provider>` at a command prompt
  - Where `<provider>` is `virtualbox`, `vmware_fusion` or `vmware_workstation`
  - See [Configuration](#configuration) for additional options for `vagrant up`

> Check out the [troubleshooting guide](FAQ.md#troubleshooting) for more information.

### Prerequisites

* [Vagrant](https://vagrantup.com/) 1.7+
* [CF CLI](https://github.com/cloudfoundry/cli)
* Internet connection required (for DNS)
* One of the following:
  - [VirtualBox](https://www.virtualbox.org/): 5.0+
  - [VMware Fusion](https://www.vmware.com/products/fusion): 8+ (for OSX)
  - [VMware Workstation](https://www.vmware.com/products/workstation): 11+ (for Windows/Linux)

> VMware requires the [Vagrant VMware](https://www.vagrantup.com/vmware#buy-now) plugin that is sold by [Hashicorp](https://hashicorp.com/).

### Configuration

The following environment variables can be set during `vagrant up` to customize the MicroPCF deployment:

1. `MICROPCF_IP` - sets the IP address to bring up the VM on
  - defaults to 192.168.11.11 locally
  - defaults to AWS-assigned public IP on AWS
1. `MICROPCF_DOMAIN` - sets an alternate alias for the system routes to be defined on
  - defaults to `local.micropcf.io` when deploying locally
  - defaults to <MICROPCF_IP>.xip.io on AWS or when `MICROPCF_IP` is set
1. `VM_CORES` (local only) - number of CPU cores to allocate on the Guest VM
  - defaults to host # of logical CPUs
1. `VM_MEMORY` (local only) - number of MB to allocate on the Guest VM 
  - defaults to 25% of host memory

### Using the Cloud Foundry CLI

Follow the instructions provided at the end of `vagrant up` to connect to MicroPCF:

```
==> default: MicroPCF is now running.
==> default: To begin using MicroPCF, please run:
==> default: 	cf api api.local.micropcf.io --skip-ssl-validation
==> default: 	cf login
==> default: Email: admin
==> default: Password: admin
```

> `local.micropcf.io` above will show the domain configured for your MicroPCF instance.

To stage a simple app on MicroPCF, `cd` into the app directory and run `cf push <APP_NAME>`.

See cf documentation for information on [deploying apps](http://docs.cloudfoundry.org/devguide/deploy-apps/) and [attaching services](http://docs.cloudfoundry.org/devguide/services/).

## Contributing

If you are interested in contributing to MicroPCF, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
