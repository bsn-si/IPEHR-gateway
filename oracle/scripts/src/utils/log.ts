import { TransactionResponse, TransactionReceipt } from "@ethersproject/abstract-provider"
import * as ethers from "ethers"
import * as chalk from "chalk"

import { getLinkTokenContract } from "./contracts"
import config from "./config"

const delimiter = chalk.bold.whiteBright("——————————————————————————————————————————————————")
const filter = (data: (string | false | undefined)[]) => data.filter(d => !!d).join("\n")

export function title(title: string) {
  console.log(delimiter)
  console.log(chalk.bgGrey.bold.white(title))
}

export function transactionResponse(
  title: string,
  { hash, blockHash, blockNumber }: TransactionResponse,
) {
  console.log(delimiter)
  console.log(chalk.bgGray.bold.white(title))

  console.log(
    filter([
      `${chalk.bold("Transaction")}`,
      `${chalk.bold("Hash")}:         ${hash}`,
      blockHash && `${chalk.bold("Block Hash")}:   ${blockHash}`,
      blockNumber && `${chalk.bold("Block Number")}: ${blockNumber}`,
    ]),
  )
}

export function transactionReceipt(title: string, receipt: TransactionReceipt) {
  console.log(delimiter)
  console.log(chalk.bgGray.bold.white(title))

  console.log(
    filter([
      `${chalk.bold("Transaction Receipt")}`,
      `${chalk.bold("Hash")}:         ${receipt.transactionHash}`,
      `${chalk.bold("Block Hash")}:   ${receipt.blockHash}`,
      `${chalk.bold("Block Number")}: ${receipt.blockNumber}`,
      `${chalk.bold("Gas Used")}:     ${receipt.gasUsed.toString()}`,
    ]),
  )
}

export function directConsumerOptions({
  jobId,
  apiHost,
  link,
  oracle,
}: Omit<typeof config.contracts.directConsumer, "address"> & { link: string; oracle: string }) {
  console.log(delimiter)
  console.log("Direct Consumer Options")

  console.log(
    filter([
      `${chalk.bold("Oracle Address")}:   ${oracle}`,
      `${chalk.bold("Link Address")}:     ${link}`,
      `${chalk.bold("API Host")}:         ${apiHost}`,
      `${chalk.bold("Chainlink Job Id")}: ${jobId}`,
    ]),
  )
}

export function stats({
  total,
  period,
}: {
  total: ethers.BigNumber[]
  period?: ethers.BigNumber[]
}) {
  console.log(delimiter)
  console.log("Current statistics")

  const details = ([documents, patients, time]: ethers.BigNumber[]) =>
    filter([
      `${chalk.bold("Documents")}:   ${documents.toString()}`,
      `${chalk.bold("Patients")}:    ${patients.toString()}`,
      `${chalk.bold("Time/Period")}: ${time.toString()}`,
    ])

  console.log(filter([`${chalk.bold("Total")}:`, details(total)]))
  period && console.log(filter([`${chalk.bold("Period")}:`, details(period)]))
}

export async function signer(signer: ethers.Wallet) {
  console.log(delimiter)
  console.log(chalk.bgGray.bold.white("Signer"))

  const balances: any = {
    eth: await signer.getBalance(),
  }

  if (config.chainlink.token.address) {
    const contract = await getLinkTokenContract(signer, config.chainlink.token.address)
    balances.link = await contract.balanceOf(signer.address)
  }

  console.log(
    filter([
      `${chalk.bold("Address")}: ${signer.address}`,
      `${chalk.bold("Balances")}:`,
      `    ${chalk.bold("ETH")}: ${ethers.utils.formatEther(balances.eth)}`,
      `    ${chalk.bold("Link")}: ${
        balances.link ? ethers.utils.formatUnits(balances.link) : "Unknown"
      }`,
    ]),
  )
}