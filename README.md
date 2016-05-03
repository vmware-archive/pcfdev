# PCF Dev
> Japanese version is [here](README_ja.md)

> oss pcfdev for cloudstack-vsphere is [here](DEVELOP_CS.md)

PCF Dev is a new distribution of Cloud Foundry designed to run on a developer’s laptop or workstation.  PCF Dev gives application developers the full Cloud Foundry experience in a lightweight, easy to install package.  PCF Dev is intended for application developers who wish to develop and debug their application locally on a full-featured Cloud Foundry.  PCF Dev is also an excellent getting started environment for developers interested in learning and exploring Cloud Foundry.

> More information about the project can be found on the [FAQ](FAQ.md#general-questions).

## Open Source

This repository contains source code that allows developers to build an open source version of PCF Dev that only contains the Elastic Runtime. The binary distribution of PCF Dev that is available on the [Pivotal Network](https://network.pivotal.io/) contains other PCF components (such as the MySQL, Redis, and RabbitMQ marketplace services) that are not available in this repository.

However, we encourage you to leave any feedback or issues you may encounter regarding the full, binary distribution of PCF Dev in [this repository's Github issues](https://github.com/pivotal-cf/pcfdev/issues).

## Install

1. Download the latest `pcfdev-<VERSION>.zip` from the [Pivotal Network](https://network.pivotal.io/).
1. Unzip the `pcfdev-<VERSION>.zip`.
1. Open a terminal or command prompt and navigate to the `pcfdev-<VERSION>` folder.
1. Run `./start-osx` at a command prompt
  - See [Configuration](#configuration) for additional options

> Check out the [troubleshooting guide](FAQ.md#troubleshooting) for more information.

### Prerequisites

* [Vagrant](https://vagrantup.com/) 1.8+
* [CF CLI](https://github.com/cloudfoundry/cli)
* Internet connection required (for DNS)
* [VirtualBox](https://www.virtualbox.org/): 5.0+

### Configuration

The following environment variables can be set during `start-osx` to customize the PCF Dev deployment:

1. `PCFDEV_IP` - sets the IP address to bring up the VM on
  - defaults to 192.168.11.11 locally
  - defaults to AWS-assigned public IP on AWS
1. `PCFDEV_DOMAIN` - sets an alternate alias for the system routes to be defined on
  - defaults to `local.pcfdev.io` when deploying locally
  - defaults to `<PCFDEV_IP>.xip.io` on AWS or when `PCFDEV_IP` is set
1. `VM_CORES` (local only) - number of CPU cores to allocate on the Guest VM
  - defaults to host # of logical CPUs
1. `VM_MEMORY` (local only) - number of MB to allocate on the Guest VM
  - defaults to 25% of host memory

### Using the Cloud Foundry CLI

Follow the instructions provided at the end of `start-osx` to connect to PCF Dev:

```
==> default: PCF Dev is now running.
==> default: To begin using PCF Dev, please run:
==> default: 	cf api api.local.pcfdev.io --skip-ssl-validation
==> default: 	cf login
==> default: Email: admin
==> default: Password: admin
```

> `local.pcfdev.io` above will show the domain configured for your PCF Dev instance.

To stage a simple app on PCF Dev, `cd` into the app directory and run `cf push <APP_NAME>`.

See cf documentation for information on [deploying apps](http://docs.cloudfoundry.org/devguide/deploy-apps/) and [attaching services](http://docs.cloudfoundry.org/devguide/services/).

## Uninstall

To temporarily stop PCF Dev:

1. Open a terminal or command prompt and navigate to the `pcfdev-<VERSION>` folder.
1. Run `./stop-osx` at a command prompt
  - You can use the `start-osx` script to resume the stopped PCF Dev instance

To permanently destroy PCF Dev:

1. Open a terminal or command prompt and navigate to the `pcfdev-<VERSION>` folder.
1. Run `./destroy-osx` at a command prompt

## Contributing

If you are interested in contributing to PCF Dev, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2016 [Pivotal Software, Inc](http://www.pivotal.io/).

PCF Dev uses a version of Monit that can be found [here](https://github.com/pivotal-cf/pcfdev-monit), under the GPLv3 license.
