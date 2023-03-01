import { getOracleContract, getProviderWithAccountSigner, config, log } from "../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()
  await log.signer(signer)

  const contract = getOracleContract(signer, config.chainlink.oracle.address)

  log.title(`Add permissions for chainlink "${config.chainlink.address}". Wait confirmations...`)
  const tx = await contract.setAuthorizedSenders([config.chainlink.address])
  const receipt = await tx.wait(1)
  log.transactionReceipt("Permissions granted", receipt)
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
