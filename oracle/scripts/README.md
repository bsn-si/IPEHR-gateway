# Helper Scripts
Common scripts for routine actions.

- `link-token:publish` - Publish new instance of link token contract to blockchain
- `chainlink:fill` - For fill chainlink node account balance from caller, for transactions payment
- `mock-server:start` - For start test mock server of IPEHR stats
- `stat:balances` - Show balances for caller & chainlink node account
- `oracle:publish` - For publish `Operator.sol` contract & add permissions for chainlink account
- `oracle:grant` - For add permissions for chainlink account in existing `Operator.sol` contract
- `direct-consumer:publish` - For publish example of direct consumer contract 
- `direct-consumer:request` - For make new request to operator for request latest data
- `direct-consumer:call` - For get current saved stats
- `cron-statistics:publish` - For publish storage contract for data
- `cron-statistics:call` - For get current saved stats in storage contract
- `cron-statistics-consumer:publish` - For publish example contract that request data from storage contract
- `cron-statistics-consumer:call` - For get current saved stats in consumer contract

# Usage
Install all dependencies & run as common npm scripts.

``` bash
npm install
npm run link-token:publish
# ...etc
```

# Config
At now scripts support only values from config file. Config file is `config.jsonc` in root of this folder.
By default required options is `account` & `node`. 

#### Account
For make requests you need setup signer account who will make requests.

`account.file` - required file name, placed in `src/assets/accounts/` folder, is encrypted account in json format.
`account.password` - required password for encrypted account file.

#### Node
`node.url` - Websocket or https address to RPC server, for make transactions. (example - infura).

## Required options for scripts
- `link-token:publish` - only common `account` & `node`
- `chainlink:fill` - also required `chainlink.address`, `chainlink.token.address`, and enough funds in the account
- `mock-server:start` - has no dependencies
- `stat:balances` - `account`, `chainlink.address`, `chainlink.token.address`
- `oracle:publish` - `account`, `chainlink.address`, `chainlink.token.address`
- `oracle:grant` - `account`, `chainlink.address`, `chainlink.token.address`
- `direct-consumer:publish` - `account`, `chainlink.token.address`, `chainlink.oracle.address`, `contracts.directConsumer.jobId`, `contracts.directConsumer.apiHost`, `amount.link` (amount of tokens for payments from consumer contract to oracle)
- `direct-consumer:request` - `account`, `contracts.directConsumer.address`, `chainlink.token.address`
- `direct-consumer:call` - `account`, `contracts.directConsumer.address`
- `cron-statistics:publish` - `account`, `chainlink.address`
- `cron-statistics:call` - `account`, `contracts.statistics.storageAddress`
- `cron-statistics-consumer:publish` - `account`, `contracts.statistics.storageAddress`
- `cron-statistics-consumer:call` - `account`, `contracts.statistics.consumerAddress`
