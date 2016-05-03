# PCF Dev Development for CloudStack(vSphere hypervisor)

> This is only for CloudStack environment, please refer to the official [development instructions] (DEVELOP.md) for other environment.

To develop PCF Dev you will need to have the following tools installed:

- [Packer](https://www.packer.io) v0.9.0+
- [Vagrant](https://www.vagrantup.com/) v1.8.1+
- [Virtualbox](https://www.virtualbox.org/) 5.0+ (to build Virtualbox boxes)
- [Vagrant CloudStack plugin](https://github.com/schubergphilis/vagrant-cloudstack) v1.3.0+ and an CloudStack Account (to build CloudStack boxes)

## Install the CloudStack plugin to vagrant

```bash
vagrant plugin install vagrant-cloudstack
```

## Clone the PCF Dev source

```bash
git clone --recursive https://github.com/pivotal-cf/pcfdev.git
```

> pcfdev v0.11.0 is the tested environment.

### Building a PCF Dev Box

To build OSS-only PCF Dev Vagrant boxes, run:

```bash
cd images
./build -only=vmware-iso
vagrant box add --force output/pcfdev-vmware-v0.box --name pcfdev/pcfdev
```
> go version should be 1.5.3 or later.

Build options:
* `-only=` with one or more of the following comma-separated builders: `virtualbox-iso`, `vmware-iso`, and/or `amazon-ebs`
* `-debug` to build all boxes in debug mode, pausing between each step with SSH login available

### Upload the image as template

It totally depends on which cloudstack environment you are using. However you need to upload the image to create the template. Then you can get the template_id.

Using ovftool to convert the image to ova.

```
$ cd image
$ ovftool -st=vmx -tt=ova oss-2016-03-28_2339.vmx oss-2016-03-28_2339.ova
Opening VMX source: oss-2016-03-28_2339.vmx
Opening OVA target: oss-2016-03-28_2339.ova
Writing OVA package: oss-2016-03-28_2339.ova
Transfer Completed                    
Completed successfully
$ mv oss-2016-03-28_2339.ova ../
```

Then upload it to your cloudstack environment.

See the details in [Working with Templates](http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/4.8/templates.html#vsphere-templates-and-isos) at Cloud Stack Docs.


### Deploying a CloudStack-vSphere PCF Dev box

The Vagrantfile for CloudStack-vSphere at the root of the repo is configured to run locally-built PCF Dev boxes.

```bash
cd ../..
vagrant up --provider=CloudStack
```

## Contributing

If you are interested in contributing to PCF Dev, please refer to [CONTRIBUTING](CONTRIBUTING.md).

## Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
