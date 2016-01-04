# Troubleshooting microPCF

####"The box you're attempting to add has no available version..."

If you're getting the following error on `vagrant up`, chances are you've used the wrong Vagrantfile.
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

Please make sure that you downloaded your Vagrantfile from the [MicroPCF releases page](https://github.com/pivotal-cf/micropcf/releases) or from our [nightly builds](https://micropcf.s3.amazonaws.com/nightly/index.html). It will have a name like Vagrantfile-v0.3.0.base.

There's no need to clone the repository in order to use microPCF unless you're attempting to develop microPCF itself.

####"Invalid SSL Cert for api.local.micropcf.io"

This occurs when you have forgotten to append `--skip-ssl-validation` to your `cf api` command. Instead enter:

`cf api api.local.micropcf.io --skip-ssl-validation`

####"Error opening SSH connection: ssh: handshake failed: Unable to verify identity of host."

This occurs when you attempt to `cf ssh` without passing `-k` to skip host key validation. Instead try:

`cf ssh -k <YOUR-APP-NAME>`

####"HTTP server doesn't seem to support byte ranges"

If you get the following error on `vagrant up` after a previously aborted `vagrant up`, you need to initiate a new download by going to ~/.vagrant.d/tmp and deleting the partial download file.

```
==> default: Box download is resuming from prior download progress
....
HTTP server doesn't seem to support byte ranges. Cannot resume.
```

####Getting a 502 on `cf api`

If you've successfully run `vagrant up` and can `vagrant ssh` into the machine and see that monit services are running properly, it may be a problem with your DNS settings. Try changing your network settings to use the Google DNS at 8.8.8.8 and run `cf api api.local.micropcf.io --skip-ssl-validation` again.

####What external ports are unavailable to bind as TCP routes?

The following ports are reserved for use by microPCF:

|     |     |     |     |     |
|-----|-----|-----|-----|-----|
|22   |53   |80   |1169 |1234 |
|1700 |2222 |2380 |4001 |4222 |
|4223 |7001 |7777 |8070 |8080 |
|8081 |8082 |8090 |8300 |8301 |
|8302 |8400 |8444 |8500 |8888 |
|8889 |9016 |9999 |17009|17014|
|17110|17111|17222|44445|     |

#### Running VirtualBox and VMware

If you're using both VirtualBox and VMware on the same machine, you may see this error:

```bash
The specified host network collides with a non-hostonly network!
This will cause your specified IP to be inaccessible. Please change
the IP or name of your host only network so that it no longer matches that of
a bridged or non-hostonly network.
```

In this case, one of your hypervisors has grabbed the 192.168.11.\* IP range and is preventing the other from accessing them. Use `ifconfig` to figure out which owns the network:

```bash
$ ifconfig
...
vboxnet1: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	ether 0a:00:27:00:00:01
	inet 192.168.11.1 netmask 0xffffff00 broadcast 192.168.11.255
...
```

In this case the VirtualBox interface `vboxnet1` has the network, so you can bring it down to free up the network:

```bash
sudo ifconfig vboxnet1 down
```

If VMware owns the network you'll see something like this:

```bash
$ ifconfig
...
vmnet9: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	ether 00:50:56:c0:00:09
	inet 192.168.11.1 netmask 0xffffff00 broadcast 192.168.11.255
...
```

> You can configure MicroPCF to run on an alternate IP address by setting, for example, `MICROPCF_IP=192.168.22.22` during `vagrant up`.


####Other tips:

* Ubuntu 14.04 LTS does not install a compatible version of Vagrant by default.  A compatible version can be found on the [Vagrant Downloads](http://www.vagrantup.com/downloads.html) page.
* Use an Administrator shell to deploy using VMware Workstation on Windows.

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
