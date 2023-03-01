import * as ethers from "ethers"

import { getProviderWithAccountSigner, getLinkTokenContract, config } from "../utils"

async function main() {
  const { provider, signer } = await getProviderWithAccountSigner()

  const contract = getLinkTokenContract(signer, config.chainlink.token.address)
  console.log("Transfer funds & link tokens to chainlink: ", config.chainlink.token.address)

  const transferTx = await signer.sendTransaction({
    value: ethers.utils.parseEther(config.amount.eth),
    to: config.chainlink.address,
  })

  console.log("Transfer transaction sended, wait confirmations...")

  const receipt = await transferTx.wait(1)
  console.log("Transfer confirmed: ", receipt)

  console.log("Balances after transfer: ", {
    account: ethers.utils.formatEther(await signer.getBalance()),
    chainlink: ethers.utils.formatEther(await provider.getBalance(config.chainlink.address)),
  })

  console.log("\n\n---------------------------------------------------\n\n")

  const transferTokensTx = await contract.transfer(
    config.chainlink.address,
    ethers.utils.parseUnits(config.amount.link),
  )

  console.log("Transfer tokens transaction sended, wait confirmations...")

  const tokensReceipt = await transferTokensTx.wait()
  console.log("Transfer confirmed: ", tokensReceipt)

  console.log("Token balances after transfer: ", {
    account: ethers.utils.formatUnits(await contract.balanceOf(signer.address)),
    chainlink: ethers.utils.formatUnits(await contract.balanceOf(config.chainlink.address)),
  })
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
