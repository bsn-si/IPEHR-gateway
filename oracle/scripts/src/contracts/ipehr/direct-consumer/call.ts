import {
  getProviderWithAccountSigner,
  getDirectConsumerContract,
  config,
  log,
} from "../../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()
  await log.signer(signer)

  const contract = getDirectConsumerContract(signer, config.contracts.directConsumer.address)
  log.title(`Show current data from direct consumer "${contract.address}", wait results...`)

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
