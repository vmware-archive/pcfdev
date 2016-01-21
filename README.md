# MicroPCF (Pivotal Cloud Foundry)

MicroPCF gives developers an easy way to download and install Cloud Foundry as a [Vagrant](https://www.vagrantup.com/) box.  It currently runs on [VirtualBox](https://www.virtualbox.org/), [VMware Fusion](https://www.vmware.com/products/fusion) (Mac), [VMware Workstation](https://www.vmware.com/products/workstation) (Windows/Linux), or [Amazon Web Services](http://aws.amazon.com/).

## Deploy MicroPCF

MicroPCF is distributed by versioned `Vagrantfile` downloads from either [Github releases](https://github.com/pivotal-cf/micropcf/releases) or the [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html) page.

1. Download a versioned `Vagrantfile-<VERSION>` into a new folder
1. Rename the downloaded file to `Vagrantfile`
1. `vagrant up`

#### On Mac or Linux

```bash
$ mkdir <WORKSPACE> && cd <WORKSPACE>
$ curl -L https://github.com/pivotal-cf/micropcf/releases/download/v0.4.0/Vagrantfile-v0.4.0.base -o Vagrantfile
$ vagrant up --provider <PROVIDER>
```

- `<WORKSPACE>` should be an empty local folder
- `<PROVIDER>` should be set to one of `virtualbox`, `vmware_fusion`, `vmware_workstation`, or `aws`.

#### On Windows

1. Download a versioned `Vagrantfile-<VERSION>`
1. Move it to a new folder
1. Rename the file to `Vagrantfile`
1. Open a command shell and `cd` to the newly-created folder
1. `vagrant up`

### Prerequisites

#### Local deployments 

* [Vagrant](https://vagrantup.com/) 1.7+
* [CF CLI](https://github.com/cloudfoundry/cli)
* One of the following virtualization platforms:
  - [VirtualBox](https://www.virtualbox.org/): 5.0+
  - [VMware Fusion](https://www.vmware.com/products/fusion): 8+ (for OSX)
  - [VMware Workstation](https://www.vmware.com/products/workstation): 11+ (for Windows)

> VMware requires the [Vagrant VMware](https://www.vagrantup.com/vmware#buy-now) plugin that is sold by [Hashicorp](https://hashicorp.com/).

#### Remote deployments (Amazon Web Services)

* [Vagrant](https://vagrantup.com/) 1.7+
* [CF CLI](https://github.com/cloudfoundry/cli)
* [Vagrant AWS](https://github.com/mitchellh/vagrant-aws) plugin: `vagrant plugin install vagrant-aws`
* Credentials for an AWS account
* Default Security Group configured to allow ingress traffic to MicroPCF:
  - HTTP (tcp/80) from 0.0.0.0/0
  - SSH (tcp/22) from 0.0.0.0/0
  - HTTPS (tcp/443) from 0.0.0.0/0
  - Custom (tcp/2222) from 0.0.0.0/0  (for `cf ssh` support)
* Environment Variables set as follows:

```bash
export AWS_ACCESS_KEY_ID=<...>
export AWS_SECRET_ACCESS_KEY=<...>
export AWS_SSH_PRIVATE_KEY_NAME=<...> # name of the remote SSH key in AWS
export AWS_SSH_PRIVATE_KEY_PATH=<...> # path to the local SSH key
export AWS_INSTANCE_NAME=<...> # optional
export AWS_REGION=<...> # optional, defaults to us-east-1
```

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

##### Use CF CLI to interact with MicroPCF

By default, you can connect the `cf` utility with the following:

```bash
$ cf api api.local.micropcf.io --skip-ssl-validation
$ cf login -u admin -p admin -o micropcf-org -s micropcf-space
```

> Replace `local.micropcf.io` above with the output from the `vagrant up` command giving the target

To stage a simple app on MicroPCF, `cd` into the app directory and run `cf push <APP_NAME>`.

See cf documentation for information on [deploying apps](http://docs.cloudfoundry.org/devguide/deploy-apps/) and [attaching services](http://docs.cloudfoundry.org/devguide/services/).

## FAQ

See our [FAQ](FAQ.md).

## Contributing

If you are interested in contributing to MicroPCF, please refer to the [contributing guidelines](CONTRIBUTING.md) and [development instructions](DEVELOP.md).

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
