# ACFG

Configuration-driven Assetto Corsa server manager that makes it easy to start multiple instances from a preset. Everything
is saved in memory only, there's no database.

Ports will be allocated and configured automatically for you. If you have a new config and want to reuse the existing
server just drag and drop the preset on the running server.

## How to use

This assumes that you have a working installation of an Assetto Corsa server via the [official manual](http://www.assettocorsa.net/forum/index.php?faq/dedicated-server-manual.28/) or [steamcmd (App ID: 302550)](https://developer.valvesoftware.com/wiki/SteamCMD).

Download the latest release binary and upload it to your server where you have a working AC Server instance.
Refer to the env vars below and configure them accordingly.

Example run script:
```bash
#!/bin/bash
export ACFG_SERVER_IP=x.x.x.x
export ACFG_PORT=13337
export ACFG_SERVER_CFGS_DIR=cfgs
export ACFG_SERVER_LOGS_DIR=logs
export ACFG_ACSERVER_DIR=../srv-ac
export ACFG_JWT_SECRET=_top_secret_change_me_
export ACFG_USERS=user1:pass1,user2:pass2
export ACFG_STRACKER_DIR=/path/to/stracker/stracker_linux_x86

./acfg-linux-amd64
```

## Environment Config

#### Required
| Variable | Description |
|----------------------------|------------------------------|
| ACFG_SERVER_IP | The server IP to connect to |
| ACFG_ACSERVER_DIR | The directory path to the acServer binary (without the binary) |
| ACFG_JWT_SECRET | A randomly generated JSON Web Token secret |
| ACFG_USERS | The user credentials, e.g. `user1:password1,user2:password2` |

#### Optional
| Variable | Default | Description |
|----------------------------|-----------------------|------------------------------|
| ACFG_PORT | 1337 | The port which ACFG should run on |
| ACFG_SERVER_CFGS_DIR | cfgs | The directory where temporary configs are saved |
| ACFG_SERVER_LOGS_DIR | logs | The directory where instance logs are saved |
| ACFG_ACSERVER_BINARY | acServer |  |
| ACFG_STRACKER_DIR | - | The directory to the stracker binary |
| ACFG_IS_PROD | true | Whether to disable some dev tools, you can probably ignore this |

For details, see the [config.go](server/app/config.go) file.

## Stracker
You can start stracker with each instance if you configured the basics and ACFG knows about it via the env config.
The following `stracker.ini` values will be set automatically for each instance:

```ini
[STRACKER_CONFIG]
ac_server_cfg_ini
listening_port

[HTTP_CONFIG]
listen_port
```

All other options like database config should be done by you.

## Development

#### Build Dependencies
- [dep](https://github.com/golang/dep) package manager to install Go project dependencies
- [fileb0x](https://github.com/UnnoTed/fileb0x) to create a single binary including static client code

#### Dependencies
Install the JS dependencies using Yarn in the `client` directory, then install all Go dependencies for the dummy
servers and the server itself using `dep ensure` in each directory.

#### Symlink dummy servers
Add symlinks for the dummy servers of your platform (darwin or linux) in `dummysrv/{acserver,stracker}`.
- `ln -s stracker-darwin-amd64 stracker`
- `ln -s acServer-darwin-amd64 acServer`

#### Run API Server
This assumes that there is a basic directory structure like `_acfg/{cfgs,logs}`.

```
$ cd server
$ ACFG_ACSERVER_DIR=../dummysrv/acserver/ ACFG_SERVER_IP=localhost ACFG_SERVER_CFGS_DIR=../_acfg/cfgs ACFG_SERVER_LOGS_DIR=../_acfg/logs ACFG_STRACKER_DIR=../dummysrv/stracker/ ACFG_IS_PROD=0 ACFG_JWT_SECRET="__secret__" ACFG_USERS="admin:pass" go run main/main.go
```

#### Start & watch client
```
$ cd client
$ yarn start
```

## Docker
Not currently supported due to dynamic port allocation. Maybe via port ranges?
```
docker run --name acfg -itd -e ACFG_SERVER_IP=x.x.x.x -v /path/srv-ac:/ac -v /path/docker-images/acfg/cfgs:/cfgs -v /path/docker-images/acfg/logs:/logs -p 8090:1337 acfg
```
