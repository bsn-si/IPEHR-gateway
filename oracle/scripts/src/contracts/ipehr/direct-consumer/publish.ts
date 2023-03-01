import * as ethers from "ethers"

import {
  getDirectConsumerContractFactory,
  getProviderWithAccountSigner,
  getLinkTokenContract,
  config,
  log,
} from "../../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()
  const gasPrice = await signer.getGasPrice()

  log.title("Publish IPEHR direct consumer contract")
  await log.signer(signer)

  const factory = getDirectConsumerContractFactory(signer)

  log.directConsumerOptions(
    Object.assign({}, config.contracts.directConsumer, {
      link: config.chainlink.token.address,
      oracle: config.chainlink.oracle.address,
    }),
  )

  const contract = await factory.deploy(
    config.chainlink.token.address,
    config.chainlink.oracle.address,
    config.contracts.directConsumer.jobId,
    config.contracts.directConsumer.apiHost,
    {
      gasPrice,
    },
  )

  log.transactionResponse("Contract published, wait confirmations...", contract.deployTransaction)
  const txReceipt = await contract.deployTransaction.wait(1)
  log.transactionReceipt(`Contract published & confirmed: ${contract.address}`, txReceipt)

  // prettier-ignore
  log.title(`Transfer '${config.amount.link}' link tokens \nto contract '${contract.address}', wait confirmations`)
  const linkToken = getLinkTokenContract(signer, config.chainlink.token.address)
  const amount = ethers.utils.parseEther(config.amount.link)

  const transferTokensTx = await linkToken.transfer(contract.address, amount)
  const transferReceipt = await transferTokensTx.wait(1)
  log.transactionReceipt("Link tokens transferred to contract", transferReceipt)
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
