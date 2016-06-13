# PCF Dev

PCF Dev is a new distribution of Cloud Foundry designed to run on a developer’s laptop or workstation.  PCF Dev gives application developers the full Cloud Foundry experience in a lightweight, easy to install package.  PCF Dev is intended for application developers who wish to develop and debug their application locally on a full-featured Cloud Foundry.  PCF Dev is also an excellent getting started environment for developers interested in learning and exploring Cloud Foundry.

> More information about the project can be found on the [FAQ](FAQ.md#general-questions).

## Open Source

This repository contains source code that allows developers to build an open source version of PCF Dev that only contains the Elastic Runtime. The binary distribution of PCF Dev that is available on the [Pivotal Network](https://network.pivotal.io/) contains other PCF components (such as the MySQL, Redis, and RabbitMQ marketplace services) that are not available in this repository.

However, we encourage you to leave any feedback or issues you may encounter regarding the full, binary distribution of PCF Dev in [this repository's Github issues](https://github.com/pivotal-cf/pcfdev/issues).

## Install

1. Download the latest `pcfdev-VERSION-PLATFORM.zip` from the [Pivotal Network](https://network.pivotal.io/).
1. Unzip the zip file and navigate to its containing folder.
1. Run `cf install-plugin pcfdev-VERSION-PLATFORM` at your command line
1. Run `cf dev start`

> Check out the [documentation](https://docs.pivotal.io/pcf-dev/) for more information. Running `cf dev help` will display an overview of PCF Dev VM management commands.

### Prerequisites

* [CF CLI](https://github.com/cloudfoundry/cli)
* [VirtualBox](https://www.virtualbox.org/): 5.0+
* Internet connection (or [Dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) or [Acrylic](http://mayakron.altervista.org/wikibase/show.php?id=AcrylicHome)) required for wildcard DNS resolution

### Using the Cloud Foundry CLI Plugin

Follow the instructions provided at the end of `cf dev start` to connect to PCF Dev:

```
Downloading VM...
Progress: |====================>| 100%
VM downloaded
Importing VM...
Starting VM...
Provisioning VM...
Waiting for services to start...
49 out of 49 running
 _______  _______  _______    ______   _______  __   __
|       ||       ||       |  |      | |       ||  | |  |
|    _  ||       ||    ___|  |  _    ||    ___||  |_|  |
|   |_| ||       ||   |___   | | |   ||   |___ |       |
|    ___||      _||    ___|  | |_|   ||    ___||       |
|   |    |     |_ |   |      |       ||   |___  |     |
|___|    |_______||___|      |______| |_______|  |___|
is now running.
To begin using PCF Dev, please run:
	cf login -a api.local.pcfdev.io --skip-ssl-validation
Email: admin
Password: admin
```

> The `local.pcfdev.io` domain may differ slightly for your PCF Dev instance.

To stage a simple app on PCF Dev, `cd` into the app directory and run `cf push <APP_NAME>`.

See cf documentation for information on [deploying apps](http://docs.cloudfoundry.org/devguide/deploy-apps/) and [attaching services](http://docs.cloudfoundry.org/devguide/services/).

## Uninstall

To temporarily stop PCF Dev run `cf dev stop`.

To destroy your PCF Dev VM run `cf dev destroy`.

To uninstall the PCF Dev cf CLI plugin run `cf uninstall-plugin pcfdev`

## Contributing

If you are interested in contributing to PCF Dev, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2016 [Pivotal Software, Inc](http://www.pivotal.io/).

PCF Dev uses a version of Monit that can be found [here](https://github.com/pivotal-cf/pcfdev-monit), under the GPLv3 license.
