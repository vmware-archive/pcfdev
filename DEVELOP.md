# PCF Dev Development

To develop PCF Dev you will need to have the following tools installed:

- [Packer](https://www.packer.io) v0.9.0+
- [Vagrant](https://www.vagrantup.com/) v1.8.1+
- [Virtualbox](https://www.virtualbox.org/) 5.0+
- [Go](https://golang.org) 1.6.1+
- [jq](https://stedolan.github.io/jq/) 1.5+

## Clone the PCF Dev source

```bash
git clone --recursive https://github.com/pivotal-cf/pcfdev.git
```

### Building a PCF Dev Box

To build an OSS-only PCF Dev Vagrant box, run:

```bash
./bin/build -only=virtualbox-iso # pass -debug for more output
vagrant box add --force output/pcfdev-virtualbox-v0.box --name pcfdev/pcfdev
```

> Note: Support for VMware Fusion/Workstation has been discontinued. Support for AWS is temporarily suspended until a commercial version of PCF Dev becomes available from the AWS Marketplace.

### Deploying a locally-built PCF Dev box

The Vagrantfile at the root of the repo is configured to run a locally-built PCF Dev box.

```bash
cd ..
vagrant up
```

## Contributing

If you are interested in contributing to PCF Dev, please refer to [CONTRIBUTING](CONTRIBUTING.md).

## Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
