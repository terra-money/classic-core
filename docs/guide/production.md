# Running for production
In this section, we will talk about how to run a node for production. At the moment we will focus on Linux (RPM based) only.

`terrad` does not require super user account. We **strongly** recommend to run `terrad` as a normal user. However, during the setup process, we need super user permission for creating/modifying files.

This guide is for general purpose only. We recommend to read [Validator](https://docs.terra.money/features/overview) section for operating a validator node.

## Firewall configuration
`terrad` uses several TCP ports for each purposes.

* `26656` is the default port for P2P protocol. This port has to be opened to communicate with other nodes. This port must be opened for joining the network. **However,** it does not have to be opened to the public. For validator nodes, we recommend to configure `persistent_peers` and close this port to the public.
* `26657` is the default port for RPC protocol. This port is used for querying / sending transactions. In other words, this port needs to be opened for serving queries from `terracli`. It is safe _NOT_ to open this port to the public unless you are planning to run public node.
* `1317` is the default port for Lite Client Daemon(LCD), which can be executed by `terracli rest-server`. LCD provides HTTP RESTful API layer to interact with `terrad` node(RPC). You can check `https://lcd.terra.dev/swagger-ui/` out for examples. Again, you don't need to open this port unless you have use of it.
* `26660` is the default port for interacting with the [Prometheus](https://prometheus.io) database which can be used for monitoring the environment. This port is not opened in the default configuration.

## Raise the maximum number of opened files for one process
`terrad` can open more than 1024 files (which is default maximum) concurrently.
We should increase this limit.
Modify `/etc/security/limits.conf` to raise the `nofile` capability.
```
*                soft    nofile          65535
*                hard    nofile          65535
```

## Running server as a daemon
There are several ways to run a node, we recommend to register `terrad` as a `systemd` service.

### Register terrad service
We have to make a service definition file in `/etc/systemd/system` directory.

#### Sample file: `/etc/systemd/system/terrad.service`
```
[Unit]
Description=Terra Daemon
After=network.target

[Service]
Type=simple
User=terra
ExecStart=/data/terra/go/bin/terrad start
Restart=on-abort

[Install]
WantedBy=multi-user.target

[Service]
LimitNOFILE=65535
```
Modify the `Service` Section from the above given sample to suit your settings.
Note that even if we raised the number of open files for a process, we still need the `LimitNOFILE` section.

After creating a service definition file, you need to execute `systemctl daemon-reload`

### Controlling service
Use systemctl to control (start, stop, restart)

* Start: `systemctl start terrad`
* Stop: `systemctl stop terrad`
* Restart: `systemctl restart terrad`

#### Accessing log file
* Entire log: `journalctl -t terrad`
* Entire log reversed: `journalctl -t terrad -r`
* Latest and continuous: `journalctl -t terrad -f`
