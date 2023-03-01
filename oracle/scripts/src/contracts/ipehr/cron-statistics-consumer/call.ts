import {
  getStatisticsConsumerContract,
  getProviderWithAccountSigner,
  config,
  log,
} from "../../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()

  log.title("Get latest data from Statistics consumer")
  await log.signer(signer)

  const contract = getStatisticsConsumerContract(
    signer,
    config.contracts.statistics.consumerAddress,
  )

  log.title(`Show current data from statistics storage "${contract.address}", wait results...`)
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
