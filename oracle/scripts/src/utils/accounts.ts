import { readFileSync } from "fs"
import * as ethers from "ethers"
import * as path from "path"

import config from "./config"

const ACCOUNTS_DIR = path.join(__dirname, "..", "assets/accounts")

// prettier-ignore
export const ACCOUNT_WALLET_PATH = path.join(
  ACCOUNTS_DIR,
  
  config.account.file 
    ? config.account.file
    : "Account.json",
)

export const getBaseWallet = async () =>
  ethers.Wallet.fromEncryptedJson(
    readFileSync(ACCOUNT_WALLET_PATH, "utf-8"),
    config.account.password,
  )
