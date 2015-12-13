# microPCF Development

To develop microPCF you will need to have the following tools installed:

- [Packer](https://www.packer.io) v0.8.6+
- [Vagrant](https://www.vagrantup.com/) v1.7.4+
- [Direnv](https://github.com/direnv/direnv) v2.7.0+
- [Virtualbox](https://www.virtualbox.org/) 5.0+ (to build Virtualbox boxes)
- [VMWare Fusion](https://www.vmware.com/products/fusion) 8+ or [VMWare Workstation](https://www.vmware.com/products/workstation) 11+ (to build VMWare boxes)
- [Vagrant AWS plugin](https://github.com/mitchellh/vagrant-aws) v0.6.0+ and an [AWS Account](https://aws.amazon.com/) (to build AWS boxes)

## Clone the microPCF source

```bash
git clone --recursive https://github.com/pivotal-cf/micropcf.git
```

## Build microPCF

Setup your shell for building microPCF:

```bash
cd micropcf
direnv allow
```

### Building a microPCF Box

To build microPCF Vagrant boxes, run:

```bash
cd images/base
./build <build options> (see below)
vagrant box add --force output/micropcf-virtualbox-v0.box --name micropcf/base
vagrant box add --force output/micropcf-vmware-v0.box --name micropcf/base
vagrant box add --force output/micropcf-aws-v0.box --name micropcf/base
```

Build options:
* `-only=` with one or more of the following comma-separated builders: `virtualbox-iso`, `vmware-iso`, and/or `amazon-ebs`
* `-var `dev=true` to leave the `bosh-provisioner` binary and BOSH releases inside of the box for re-deployment
* `-debug` to build all boxes in debug mode, pausing between each step with SSH login available

### Deploying a locally-built microPCF box

The Vagrantfile at the root of the repo is configured to run locally-built microPCF boxes.

```bash
cd ../..
vagrant up --provider=(virtualbox|vmware_fusion|vmware_workstation|aws)
```

## Contributing

If you are interested in contributing to microPCF, please refer to [CONTRIBUTING](CONTRIBUTING.md).

## Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
