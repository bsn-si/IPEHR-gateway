import {
  getStatisticsContractFactory,
  getProviderWithAccountSigner,
  config,
  log,
} from "../../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()

  log.title("Publish contract Statistics storage")
  await log.signer(signer)

  const factory = getStatisticsContractFactory(signer)
  const contract = await factory.deploy(config.chainlink.address)

  log.transactionResponse("Contract published, wait confirmations...", contract.deployTransaction)
  const txReceipt = await contract.deployTransaction.wait(1)
  log.transactionReceipt(`Contract published & confirmed: ${contract.address}`, txReceipt)

  // Initial data for test if needed
  const setTx = await contract.setTotal(1, 2, 3)
  const setReceipt = await setTx.wait(1)
  log.transactionReceipt("Complete set test data", setReceipt)

  log.stats({
    period: await contract.getByPeriod(202212),
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
