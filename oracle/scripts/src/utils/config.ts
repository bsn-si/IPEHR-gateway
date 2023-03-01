import { readFileSync } from "fs"
import * as yargs from "yargs"
import { jsonc } from "jsonc"
import * as path from "path"

// override options arguments
const options = yargs(process.argv.slice(2)).argv

// common config file
const json = readFileSync(path.join(__dirname, "../..", "config.jsonc"), "utf-8")
const config = jsonc.parse(json)

// // crutch replacements
// config.chainlink.token.address = options["link-token"] || config.chainlink.token.address
// config.LINK_AMOUNT = options["link-amount"] || config.LINK_AMOUNT
// config.ETH_AMOUNT = options["eth-amount"] || config.ETH_AMOUNT

// // account.file?: string
// // ACCOUNT_PASSWORD: string

interface Config {
  node: {
    url: string
  }

  account: {
    password: string
    file: string
  }

  chainlink: {
    address: string
    token: {
      address: string
    }
    oracle: {
      address: string
    }
  }

  contracts: {
    directConsumer: {
      address: string
      apiHost: string
      jobId: string
    }

    statistics: {
      consumerAddress: string
      storageAddress: string
    }
  }

  amount: {
    link: string
    eth: string
  }
}

export default config as Config
