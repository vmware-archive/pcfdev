# Frequently Asked Questions

## General Questions

### What is MicroPCF?

MicroPCF is a new distribution of Cloud Foundry designed to run on a developer’s laptop or workstation.  MicroPCF gives application developers the full Cloud Foundry experience in a lightweight, easy to install package.

### Who should use MicroPCF?

MicroPCF is intended for application developers who wish to develop and debug their application locally on a full-featured Cloud Foundry.  MicroPCF is also an excellent getting started environment for developers interested in learning and exploring Cloud Foundry.

### If my application runs on MicroPCF, will it run on PCF?

Yes.  MicroPCF is designed to mirror PCF exactly.  If your application runs on MicroPCF, it will run on PCF with no modification in almost all cases.

### Why do I need Vagrant?

Vagrant is a product from Hashicorp that "provides easy to configure, reproducible, and portable work environments", allowing us to perform the task of provisioning a Cloud Foundry environment for you.  In conjunction with Atlas (formerly Vagrant Cloud), we can distribute the Linux VM required to run Cloud Foundry and then provision and configure the server using the virtualization platform of your choosing.  More information about Vagrant can be found [here](https://docs.vagrantup.com/v2/why-vagrant/index.html).

## Troubleshooting

### Why does `vagrant up` say it has "no available version" ?

Cloning the repository and running `vagrant up` from the root will result in the error below.  Please follow the [install instructions](README.md#install) and use one of the published `Vagrantfile`s.

```
The box you're attempting to add has no available version that
matches the constraints you requested. Please double-check your
settings. Also verify that if you specified version constraints,
that the provider you wish to use is available for these constriants.

Box: micropcf/base
Address: https://atlas.hashicorp.com/micropcf/base
Constraints: 0
Available versions: 0.0.1, 0.1.0, 0.2.0, .... etc
```

### Why does `cf api` fail with "Invalid SSL Cert" error ?

MicroPCF comes with a self-signed SSL certificate for its API and requires the `--skip-ssl-validation` option.  This also applies to the Spring Boot Dashboard, which requires the checkbox "Self-signed" in order to connect.

```
○ → cf api api.local.micropcf.io
Setting api endpoint to api.local.micropcf.io...
FAILED
Invalid SSL Cert for api.local.micropcf.io
TIP: Use 'cf api --skip-ssl-validation' to continue with an insecure API endpoint
```

### Why does the `cf ssh` handshake fail?

`cf ssh` requires using the `-k` option to skip host key validation since it uses self-signed certifcates.

```
○ → cf ssh app
FAILED
Error opening SSH connection: ssh: handshake failed: Unable to verify identity of host.

The fingerprint of the received key was "4d:2b:ff:a4:97:8e:25:36:a0:cc:04:bc:9d:71:c7:6c".
```

### My box download failed and I can't resume the download.  What do I do?

Prior to Vagrant 1.8.0, it is necessary to manually delete temporary files in `~/.vagrant.d/tmp` prior to running `vagrant up` again.  Newer versions of Vagrant support resuming box downloads properly.

```
==> default: Adding box 'micropcf/base' (v0.20.0) for provider: virtualbox
    default: Downloading: https://atlas.hashicorp.com/micropcf/boxes/base/versions/0.20.0/providers/virtualbox.box
==> default: Box download is resuming from prior download progress
An error occurred while downloading the remote file. The error
message, if any, is reproduced below. Please fix this error and try
again.

HTTP server doesn't seem to support byte ranges. Cannot resume.
```

### Why does `vagrant up` say my network collides with another device?

By default, MicroPCF will attempt to reserve `192.168.11.11` as its address.  If this network is already in use (for example, if you try using VMware after VirtualBox), you'll see one of the below errors.  Please see the [configuration](README.md#configuration) section to set `MICROPCF_IP` to a valid address.

We recommend trying one of the following first:

```bash
MICROPCF_IP=192.168.22.22 MICROPCF_DOMAIN=2.micropcf.io vagrant up --provider=<provider>
MICROPCF_IP=192.168.33.33 MICROPCF_DOMAIN=3.micropcf.io vagrant up --provider=<provider>
MICROPCF_IP=192.168.44.44 MICROPCF_DOMAIN=4.micropcf.io vagrant up --provider=<provider>
```

```
The specified host network collides with a non-hostonly network!
This will cause your specified IP to be inaccessible. Please change
the IP or name of your host only network so that it no longer matches that of
a bridged or non-hostonly network.
```

```
The host only network with the IP '192.168.11.11' would collide with
another device 'vboxnet'. This means that VMware cannot create
a proper networking device to route to your VM. Please choose
another IP or shut down the existing device.
```

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
