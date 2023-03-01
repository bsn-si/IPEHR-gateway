import * as ethers from "ethers"
import { getProviderWithAccountSigner, getLinkTokenContract, config } from "../utils"

async function main() {
  const { provider, signer } = await getProviderWithAccountSigner()
  const linkToken = await getLinkTokenContract(signer, config.chainlink.token.address)

  console.log("Token Contract: ", {
    address: linkToken.address,

    data: {
      totalSupply: ethers.utils.formatUnits(await linkToken.totalSupply()),
      name: await linkToken.name(),
    },

    balances: {
      chainlink: ethers.utils.formatUnits(await linkToken.balanceOf(config.chainlink.address)),
      account: ethers.utils.formatUnits(await linkToken.balanceOf(signer.address)),
    },
  })

  console.log("Native Tokens: ", {
    chainlink: ethers.utils.formatEther(await provider.getBalance(config.chainlink.address)),
    account: ethers.utils.formatEther(await signer.getBalance()),
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
