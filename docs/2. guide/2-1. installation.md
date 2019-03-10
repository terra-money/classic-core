## Install Terra

This guide will explain how to install the `terrad` and `terracli` entrypoints onto your system. With these installed on a server, you can participate in the latest testnet as either a [Full Node](./join-network.md#run-a-full-node) or a [Validator](./setup-validator.md).

### Install Go

Install `go` by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
```

::: tip
**Go 1.11.5+** is required for Terra.
:::

### Install the binaries

Next, let's install the latest version of Terra. Here we'll use the `master` branch, which contains the latest stable release.
If necessary, make sure you `git checkout` the correct
[released version](https://github.com/terra-project/terra/releases).

```bash
mkdir -p $GOPATH/src/github.com/terra-project/terra
cd $GOPATH/src/github.com/terra-project/terra
git clone https://github.com/terra-project/terra
git checkout master
make tools install
```

> *NOTE*: If you have issues at this step, please check that you have the latest stable version of GO installed.

That will install the `terrad` and `terracli` binaries. Verify that everything is OK:

```bash
$ terrad version --long
$ terracli version --long
```

`terracli` for instance should output something similar to:

```
terra-money: 0.1
git commit: 1fba7308fa226e971964cd6baad9527d4b51d9fc
vendor hash: 1aec7edfad9888a967b3e9063e42f66b28f447e6
build tags: netgo ledger
go version go1.11.5 linux/amd64
```

##### Build Tags

Build tags indicate special features that have been enabled in the binary.

| Build Tag | Description                                     |
| --------- | ----------------------------------------------- |
| netgo     | Name resolution will use pure Go code           |
| ledger    | Ledger devices are supported (hardware wallets) |


### Next

Now you can [join the public testnet](./join-network.md) or [create you own  testnet](./deploy-testnet.md)
