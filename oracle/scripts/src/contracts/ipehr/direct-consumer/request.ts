import * as ethers from "ethers"

import {
  getProviderWithAccountSigner,
  getDirectConsumerContract,
  getLinkTokenContract,
  config,
  log,
} from "../../../utils"

const TARGET_PERIOD = "202212"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()

  log.title("Request data via Oracle from Chainlink")
  await log.signer(signer)

  const contract = getDirectConsumerContract(signer, config.contracts.directConsumer.address)
  const linkToken = getLinkTokenContract(signer, config.chainlink.token.address)

  const prevTokenBalance: ethers.BigNumber = await linkToken.balanceOf(
    config.contracts.directConsumer.address,
  )

  const requestOptions = {
    gasPrice: await signer.getGasPrice(),
    gasLimit: 6000000,
  }

  log.title(`Request latest stats for contract "${contract.address}", wait confirmations...`)
  const latestTx = await contract.requestLatestStatistics(requestOptions)
  const latestReceipt = await latestTx.wait(1)
  log.transactionReceipt("Latest stats requested successfully", latestReceipt)

  log.title(`Request stats for period "${TARGET_PERIOD}", wait confirmations...`)
  const getPeriodDataTx = await contract.requestStatisticsByPeriod("202212", requestOptions)
  const periodReceipt = await getPeriodDataTx.wait(1)
  log.transactionReceipt("Latest stats requested successfully", periodReceipt)

  // prettier-ignore
  const linkTokensSpent = prevTokenBalance.sub(await linkToken.balanceOf(config.contracts.directConsumer.address))
  log.title(`For this request spent ${linkTokensSpent.toString()} tokens`)
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
