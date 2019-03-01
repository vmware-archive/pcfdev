# Frequently Asked Questions

## General Questions

### What is PCF Dev?

PCF Dev is a new distribution of Cloud Foundry designed to run on a developer’s laptop or workstation.  PCF Dev gives application developers the full Cloud Foundry experience in a lightweight, easy to install package.

### Who should use PCF Dev?

PCF Dev is intended for application developers who wish to develop and debug their application locally on a full-featured Cloud Foundry.  PCF Dev is also an excellent getting started environment for developers interested in learning and exploring Cloud Foundry.

### If my application runs on PCF Dev, will it run on PCF?

Yes.  PCF Dev is designed to mirror PCF exactly.  If your application runs on PCF Dev, it will run on PCF with no modification in almost all cases.

## Troubleshooting

### Why does `cf api` and/or `cf login` fail with an "Invalid SSL Cert" error?

PCF Dev comes with a self-signed SSL certificate for its API and requires the `--skip-ssl-validation` option.  This also applies to the Spring Boot Dashboard, which requires the checkbox "Self-signed" in order to connect.

```
○ → cf api api.local.pcfdev.io
Setting api endpoint to api.local.pcfdev.io...
FAILED
Invalid SSL Cert for api.local.pcfdev.io
TIP: Use 'cf api --skip-ssl-validation' to continue with an insecure API endpoint
```

## Networking

### Container-to-router

This is traffic from the app container to the gorouter. It is enabled by default. This allows apps to communicate with each other by using the routes published by gorouter.

### Container-to-guest

This is traffic from the app container to the virtual machine in which PCF Dev is running. It is enabled by default. This may be useful if you want to run other services inside of the guest virtual machine for your applications to use, but doing so is not encouraged. Instead, extra services should be run on the host (see below). The IP address of the guest is `192.168.11.11` or `local.pcfdev.io` (unless this address is already in use).

### Container-to-host

This is traffic from the app container to the host on which the virtual machine is running. It is enabled by default. This can be used to run services on your host that are available to your apps in PCF Dev.  The IP address of the host accessible to the app is `192.168.11.1` or `host.cfdev.sh` (unless this address is already in use). For example, in order to connect your app to a MongoDB instance running on the host on port `27017`, run the following commands:

```bash
cf create-user-provided-service my-mongo-db -p '{ "uri": "mongodb://<username>:<password>@host.pcfdev.io:27017/<database>" }'
cf bind-service <app> my-mongo-db
cf restage <app>
```

### Container-to-external

This is traffic from the app container to a destination external to the host. It allows your application to reach the internet. Traffic to public and private IP addresses is enabled by default in PCF Dev. You may remove the `all_pcfdev` security group to restrict access to only public IP addresses, as a default PCF installation would be configured.

### Container-to-container

This is traffic directly between two containers in the same PCF Dev deployment. It is useful for running applications that must communicate with each other but do not need or want a publicly-accessible route. It is not enabled and will not be available until it is supported in Pivotal Cloud Foundry.

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2015 [Pivotal Software, Inc](http://www.pivotal.io/).
