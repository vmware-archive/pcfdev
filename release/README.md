# Forge Release - `forge.tgz`

Forge base images are baked with Diego and CF components that Forge depends on.
A `forge.tgz` file is used to provision a Forge base image (for any infrastructure)
with Forge-specific components. A `forge.tgz` file is a built release of the
[Forge repository](https://github.com/cloudfoundry-incubator/forge).

## Building

To build a `forge.tgz` file, execute the `release/build` script. You must have the
following dependencies:
- Go 1.4+ with Linux/AMD64 cross compilation support
- All submodules up-to-date (use `git submodule update --init --recursive`)

## Installation

Installation is usually handled by a Terraform or Vagrant deployment configuration. All
Forge deployment configurations should adhere to the following guidelines.

To install a `forge.tgz` file, first extract the `install` folder at the root of the tarball.
- The `install/common` script must be run on every Forge instance.
- The `install/cell` script must be run on every Forge cell
  with the `forge.tgz` file as the first argument.
- The `install/brain` script must be run on every Forge brain
  with the `forge.tgz` file as the first argument.
- These may be run in any order. All of them are necessary
  to provision a colocated Forge instance.

Other directories inside of the `install` directory contain infrastructure-specific patch scripts.
- These scripts must be executed after the platform-independent install scripts described above.
- These scripts must be named `common`, `cell`, or `brain`, and should be executable in any order.
- These scripts should not depend on the `forge.tgz` file.

Finally, the `install/start` script enables all of the services and waits for them to start.
It must be run last. No other scripts should start or stop services besides `install/start`.
