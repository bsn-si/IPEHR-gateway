import { getOracleContractFactory, getProviderWithAccountSigner, config, log } from "../../utils"

async function main() {
  log.title("Connect to node & get signer")
  const { signer } = await getProviderWithAccountSigner()

  log.title("Publish contract 'Operator.sol' as oracle")
  await log.signer(signer)

  const oracleFactory = getOracleContractFactory(signer)
  const contract = await oracleFactory.deploy(config.chainlink.token.address, signer.address)

  log.transactionResponse("Contract published, wait confirmations...", contract.deployTransaction)
  const txReceipt = await contract.deployTransaction.wait(1)
  log.transactionReceipt(`Contract published & confirmed: ${contract.address}`, txReceipt)

  log.title(`Add permissions for chainlink "${config.chainlink.address}". Wait confirmations...`)
  const addPermissionTx = await contract.setAuthorizedSenders([config.chainlink.address])
  const addPermissionReceipt = await addPermissionTx.wait(1)

  log.transactionReceipt("Permissions granted", addPermissionReceipt)
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
