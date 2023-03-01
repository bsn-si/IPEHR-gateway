import {
  getStatisticsConsumerContractFactory,
  getProviderWithAccountSigner,
  config,
  log,
} from "../../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()

  log.title("Publish contract Statistics consumer")
  await log.signer(signer)

  const factory = getStatisticsConsumerContractFactory(signer)
  const contract = await factory.deploy(config.contracts.statistics.storageAddress)

  log.transactionResponse("Contract published, wait confirmations...", contract.deployTransaction)
  const txReceipt = await contract.deployTransaction.wait(1)
  log.transactionReceipt(`Contract published & confirmed: ${contract.address}`, txReceipt)

  // request initial data from storage contract
  const getTx = await contract.requestLatestStats()
  await getTx.wait(1)

  log.stats({
    total: await contract.getTotal(),
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
