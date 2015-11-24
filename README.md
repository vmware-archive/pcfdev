# Cloud Foundry Forge

<table width="100%" border="0">
  <tr>
    <td>
      <a href="http://forge.cf"><img src="https://raw.githubusercontent.com/cloudfoundry-incubator/forge/master/forge.png" align="left" width="300px" ></a>
    </td>
    <td>Forge is an open source project for running containerized workloads on a cluster. Forge bundles up http load-balancing, a cluster scheduler, log aggregation/streaming and health management into an easy-to-deploy and easy-to-use package.
    </td>
  </tr>
</table>

[ [Website](http://forge.cf) | [Latest Release](https://github.com/cloudfoundry-incubator/forge/releases/latest) | [Nightly Builds](https://forge.s3.amazonaws.com/nightly/index.html) ]

## Deploy Forge with Vagrant

A colocated deployment of Forge can be launched locally with [Vagrant](https://vagrantup.com/). You will need:

* A Forge bundle from the [latest release](https://github.com/cloudfoundry-incubator/forge/releases/latest) or [nightly builds](https://forge.s3.amazonaws.com/nightly/index.html)
* [Vagrant](https://vagrantup.com/) 1.6+ installed

> NOTE: Ubuntu 14.04 LTS does not install a compatible version of Vagrant by default. You can upgrade the version that you get out of the box by downloading the `.deb` file from [Vagrant](http://www.vagrantup.com/downloads.html).

##### Spin up a virtual environment

Unzip the Forge bundle, and switch to the vagrant directory

```bash
unzip forge-bundle-VERSION.zip
cd forge-bundle-VERSION/vagrant
vagrant up --provider virtualbox
```

This spins up a virtual environment that is accessible at `local.forge.cf`

##### Install `ltc` (the Forge CLI)

If you're running Linux: `curl -O http://receptor.local.forge.cf/v1/sync/linux/ltc`

If you're running OS X: `curl -O http://receptor.local.forge.cf/v1/sync/osx/ltc`

Finally: `chmod +x ltc`

##### Use the Forge CLI to target Forge

```bash
./ltc target local.forge.cf
```

## Deploy Forge with Terraform

A scalable cluster deployment of Forge can be launched on Amazon Web Services with [Terraform](https://www.terraform.io). You will need:

* An [Amazon Web Services account](http://aws.amazon.com/)
* [AWS Access and Secret Access Keys](http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSGettingStartedGuide/AWSCredentials.html)
* [AWS EC2 Key Pairs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html)
* [Terraform 0.6.2+](https://www.terraform.io/intro/getting-started/install.html)
* A Forge bundle from the [latest release](https://github.com/cloudfoundry-incubator/forge/releases/latest) or the [nightly builds](https://forge.s3.amazonaws.com/nightly/index.html) page

##### Configure your virtual environment

Unzip the Forge bundle, and switch to the terraform/aws directory

```bash
unzip forge-bundle-VERSION.zip
cd forge-bundle-VERSION/terraform/aws
```

Update the `terraform.tfvars` file with your AWS credentials and desired cluster configuration.

##### Deploy the cluster to AWS

```bash
terraform apply
```

Terraform will generate a `terraform.tfstate` file.  This file describes the cluster that was built - keep it around in order to modify/tear down the cluster.

##### Install `ltc` (the Forge CLI)

After a successful deployment Terraform will print the Forge target and Forge user information. Refer to the `target = <forge target>` output line to find the address of your cluster.

If you're running Linux: `curl -O http://receptor.<forge target>/v1/sync/linux/ltc`

If you're running OS X: `curl -O http://receptor.<forge target>/v1/sync/osx/ltc`

Finally: `chmod +x ltc`

##### Use the Forge CLI to target Forge

```bash
./ltc target <forge target>
```

## Development

> NOTE: These instructions are for people contributing code to Forge. If you only want to deploy Forge, see above.
> These instructions cover Vagrant/Virtualbox development.
> A similar process can be followed for Vagrant/VMWare, Vagrant/AWS, and Terraform/AWS development. More documentation is forthcoming.

To develop Forge you will need to have the following tools installed:

- Packer
- Vagrant
- Virtualbox
- Direnv _(optional)_

### Clone the Forge source

```bash
git clone --recursive https://github.com/cloudfoundry-incubator/forge.git
```

### Build Forge

Setup your shell for building Forge:

```bash
cd forge
direnv allow # or: source .envrc
```

#### Building a Forge Box

If you change any Diego components, you'll need to build a local Vagrant box with your changes.
If you don't plan to change any Diego components, you can update the `box_version` property
in `vagrant/Vagrantfile` to point to a [pre-built Forge box on Atlas](https://atlas.hashicorp.com/forge/boxes/colocated)
and skip this step.

```bash
bundle
cd vagrant
./build -only=virtualbox-iso
vagrant box add --force forge-virtualbox-v0.box --name forge/colocated
```

#### Building a release of Forge

```bash
# in forge/vagrant
../release/build forge.tgz
```

### Deploy Forge

Once you have a Forge tarball, use `vagrant` to deploy Forge:

```bash
# in forge/vagrant
vagrant up --provider=virtualbox
```

### Install ltc

Compiling ltc is as simple as using `go install`:

```bash
go install github.com/cloudfoundry-incubator/ltc
```

### Test the running Forge Cluster

```bash
ltc target local.forge.cf
ltc test -v
```

## Contributing

If you are interested in contributing to Forge, please refer to [CONTRIBUTING](CONTRIBUTING.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
