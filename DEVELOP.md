# PCF Dev Development

To develop PCF Dev you will need to have the following tools installed:

- [Packer](https://www.packer.io) v0.9.0+
- [Vagrant](https://www.vagrantup.com/) v1.8.1+
- [Virtualbox](https://www.virtualbox.org/) 5.0+ (to build Virtualbox boxes)
- [VMWare Fusion](https://www.vmware.com/products/fusion) 8+ or [VMWare Workstation](https://www.vmware.com/products/workstation) 11+ (to build VMWare boxes)
- [Vagrant AWS plugin](https://github.com/mitchellh/vagrant-aws) v0.6.0+ and an [AWS Account](https://aws.amazon.com/) (to build AWS boxes)

## Clone the PCF Dev source

```bash
git clone --recursive https://github.com/pivotal-cf/pcfdev.git
```

### Building a PCF Dev Box

To build OSS-only PCF Dev Vagrant boxes, run:

```bash
cd images
./build <build options> (see below)
vagrant box add --force output/pcfdev-virtualbox-v0.box --name pcfdev/pcfdev
vagrant box add --force output/pcfdev-vmware-v0.box --name pcfdev/pcfdev
vagrant box add --force output/pcfdev-aws-v0.box --name pcfdev/pcfdev
```

Build options:
* `-only=` with one or more of the following comma-separated builders: `virtualbox-iso`, `vmware-iso`, and/or `amazon-ebs`
* `-debug` to build all boxes in debug mode, pausing between each step with SSH login available

### Deploying a locally-built PCF Dev box

The Vagrantfile at the root of the repo is configured to run locally-built PCF Dev boxes.

```bash
cd ../..
vagrant up --provider=(virtualbox|vmware_fusion|vmware_workstation|aws)
```

## Contributing

If you are interested in contributing to PCF Dev, please refer to [CONTRIBUTING](CONTRIBUTING.md).

## Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
