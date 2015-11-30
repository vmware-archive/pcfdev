# MicroPCF

<table width="100%" border="0">
  <tr>
    <td>
      <a href="http://micropcf.cf"><img src="https://raw.githubusercontent.com/cloudfoundry-incubator/micropcf/master/micropcf.png" align="left" width="300px" ></a>
    </td>
    <td>MicroPCF is an open source project for running containerized workloads on a cluster. MicroPCF bundles up http load-balancing, a cluster scheduler, log aggregation/streaming and health management into an easy-to-deploy and easy-to-use package.
    </td>
  </tr>
</table>

[ [Website](http://micropcf.cf) | [Latest Release](https://github.com/cloudfoundry-incubator/micropcf/releases/latest) | [Nightly Builds](https://micropcf.s3.amazonaws.com/nightly/index.html) ]

## Deploy MicroPCF with Vagrant

A colocated deployment of MicroPCF can be launched locally with [Vagrant](https://vagrantup.com/). You will need:

* A MicroPCF bundle from the [latest release](https://github.com/cloudfoundry-incubator/micropcf/releases/latest) or [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html)
* [Vagrant](https://vagrantup.com/) 1.6+ installed

> NOTE: Ubuntu 14.04 LTS does not install a compatible version of Vagrant by default. You can upgrade the version that you get out of the box by downloading the `.deb` file from [Vagrant](http://www.vagrantup.com/downloads.html).

##### Spin up a virtual environment

Unzip the MicroPCF bundle, and switch to the vagrant directory

```bash
unzip micropcf-bundle-VERSION.zip
cd micropcf-bundle-VERSION/vagrant
vagrant up --provider virtualbox
```

This spins up a virtual environment that is accessible at `local.micropcf.cf`

##### Install `ltc` (the MicroPCF CLI)

If you're running Linux: `curl -O http://receptor.local.micropcf.cf/v1/sync/linux/ltc`

If you're running OS X: `curl -O http://receptor.local.micropcf.cf/v1/sync/osx/ltc`

Finally: `chmod +x ltc`

##### Use the MicroPCF CLI to target MicroPCF

```bash
./ltc target local.micropcf.cf
```

## Deploy MicroPCF with Terraform

A scalable cluster deployment of MicroPCF can be launched on Amazon Web Services with [Terraform](https://www.terraform.io). You will need:

* An [Amazon Web Services account](http://aws.amazon.com/)
* [AWS Access and Secret Access Keys](http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSGettingStartedGuide/AWSCredentials.html)
* [AWS EC2 Key Pairs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html)
* [Terraform 0.6.2+](https://www.terraform.io/intro/getting-started/install.html)
* A MicroPCF bundle from the [latest release](https://github.com/cloudfoundry-incubator/micropcf/releases/latest) or the [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html) page

##### Configure your virtual environment

Unzip the MicroPCF bundle, and switch to the terraform/aws directory

```bash
unzip micropcf-bundle-VERSION.zip
cd micropcf-bundle-VERSION/terraform/aws
```

Update the `terraform.tfvars` file with your AWS credentials and desired cluster configuration.

##### Deploy the cluster to AWS

```bash
terraform apply
```

Terraform will generate a `terraform.tfstate` file.  This file describes the cluster that was built - keep it around in order to modify/tear down the cluster.

##### Install `ltc` (the MicroPCF CLI)

After a successful deployment Terraform will print the MicroPCF target and MicroPCF user information. Refer to the `target = <micropcf target>` output line to find the address of your cluster.

If you're running Linux: `curl -O http://receptor.<micropcf target>/v1/sync/linux/ltc`

If you're running OS X: `curl -O http://receptor.<micropcf target>/v1/sync/osx/ltc`

Finally: `chmod +x ltc`

##### Use the MicroPCF CLI to target MicroPCF

```bash
./ltc target <micropcf target>
```

## Development

> NOTE: These instructions are for people contributing code to MicroPCF. If you only want to deploy MicroPCF, see above.
> These instructions cover Vagrant/Virtualbox development.
> A similar process can be followed for Vagrant/VMWare, Vagrant/AWS, and Terraform/AWS development. More documentation is forthcoming.

To develop MicroPCF you will need to have the following tools installed:

- Packer
- Vagrant
- Virtualbox
- Direnv _(optional)_

### Clone the MicroPCF source

```bash
git clone --recursive https://github.com/cloudfoundry-incubator/micropcf.git
```

### Build MicroPCF

Setup your shell for building MicroPCF:

```bash
cd micropcf
direnv allow # or: source .envrc
```

#### Building a MicroPCF Box

If you change any Diego components, you'll need to build a local Vagrant box with your changes.
If you don't plan to change any Diego components, you can update the `box_version` property
in `vagrant/Vagrantfile` to point to a [pre-built MicroPCF box on Atlas](https://atlas.hashicorp.com/micropcf/boxes/colocated)
and skip this step.

```bash
bundle
cd vagrant
./build -only=virtualbox-iso
vagrant box add --force micropcf-virtualbox-v0.box --name micropcf/colocated
```

#### Building a release of MicroPCF

```bash
# in micropcf/vagrant
../release/build micropcf.tgz
```

### Deploy MicroPCF

Once you have a MicroPCF tarball, use `vagrant` to deploy MicroPCF:

```bash
# in micropcf/vagrant
vagrant up --provider=virtualbox
```

### Install ltc

Compiling ltc is as simple as using `go install`:

```bash
go install github.com/cloudfoundry-incubator/ltc
```

### Test the running MicroPCF Cluster

```bash
ltc target local.micropcf.cf
ltc test -v
```

## Contributing

If you are interested in contributing to MicroPCF, please refer to [CONTRIBUTING](CONTRIBUTING.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
