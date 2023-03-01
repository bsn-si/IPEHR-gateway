import { getLinkTokenContractFactory, getProviderWithAccountSigner, log } from "../../utils"

async function main() {
  log.title("Connect to node & get signer")

  const { signer } = await getProviderWithAccountSigner()
  log.title("Publish contract 'LinkToken'")
  await log.signer(signer)

  const factory = getLinkTokenContractFactory(signer)
  const contract = await factory.deploy()

  log.transactionResponse("Contract published, wait confirmations...", contract.deployTransaction)
  const txReceipt = await contract.deployTransaction.wait(1)
  log.transactionReceipt(`Contract published & confirmed: ${contract.address}`, txReceipt)
}

main()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(() => {
    process.exit()
  })
