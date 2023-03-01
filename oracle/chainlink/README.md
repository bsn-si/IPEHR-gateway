# Chainlink Setup
This is a simple quick guide on how to deploy a local chainlink node for development and testing. You can see the full guide on the [official website](https://docs.chain.link/chainlink-nodes/v1/running-a-chainlink-node).

### Build & Install
*For build chainlink you need - go-lang compiler with configured GOPATH, and nodejs v16+ with pnpm installed*

``` bash
git clone https://github.com/smartcontractkit/chainlink
cd chainlink/
make install
```

### Database
You need create new empty database for chainlink in postgresql.
For example in ubuntu after installing and configuring postgresql.

``` bash
sudo -u postgres psql

create user chainlink_admin with encrypted password 'your_password';
create database chainlink_dev;
grant all privileges on database chainlink_dev to chainlink_admin;
```

### Config & Migrations
For set options for chainlink you need run process with env variables. You can check sample configuration for chainlink in [env.example](env.example), 
please set `DATABASE_URL` to your empty database for chainlink. Also you need choose network and set `LINK_CONTRACT_ADDRESS`, `ETH_URL`, `ETH_CHAIN_ID` for node.

> If need publish your own instance of Link Contract token, please check command `link-token:publish` in [helper scripts](../scripts/README.md).

After write config file with env variables you need migrate all tables & schemes. You can make that with `run_with_env.sh` helper.

``` bash
./run_with_env.sh ./.env chainlink node db migrate
```

### Run Node
For run node you need start process with env config options. For first time you need enter password for master key, and create account for access to Operator UI.

``` bash
./run_with_env.sh ./.env chainlink node start
```

After that you can open Operator UI: http://localhost:6688/

### Run process with pm2
If you need run node in daemon in pm2, you can update options in `ecosystem.config.js` and run it with:

``` bash
pm2 start ecosystem.config.js
```
