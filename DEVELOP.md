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

### Customizing PCF Dev

Our build tool has the ability to build compiled releases or releases from source. By default, it will try to build releases that have been compiled by the PCF Dev team. If you have a *non-compiled* release that present on your workstation, you can configure the build to use it using the **path:** key. Simply edit the versions.json file at the root of this repo like `"cf" :` is done below:

```json
{
  "releases": {
    "cf" : {
      "path": /Users/pivotal/[path-to-release-folder]
    },
    "diego" : {
      "version": "v0.1480.0",
      "sha1": "bfd87d6ef08458e19e2abc6fc6888ba9ac29fde6",
      "source_location": "https://github.com/cloudfoundry/diego-release",
      "compiled_release_url" : "https://s3.amazonaws.com/pcfdev/compiled-releases/diego-8d1450da393eae98d565b9e0e7154c742e75e513.tgz"
    },
```

If you would like to a different *compiled* release than is offered in the versions.json, simply make sure that the appropriate keys are modified.

> Note: any necessary manifest changes can be done to the manifest.yml file at the root of this repo for a successful build.

## Contributing

If you are interested in contributing to PCF Dev, please refer to [CONTRIBUTING](CONTRIBUTING.md).

## Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
