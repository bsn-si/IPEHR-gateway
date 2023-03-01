import * as ethers from "ethers"

import { getBaseWallet } from "../utils"
import config from "./config"

export async function getProviderWithSigner(wallet: ethers.Wallet) {
  const provider = new ethers.providers.WebSocketProvider(config.node.url)
  const signer = wallet.connect(provider)

  return {
    provider,
    signer,
  }
}

export async function getProviderWithAccountSigner() {
  const wallet = await getBaseWallet()
  return getProviderWithSigner(wallet)
}
