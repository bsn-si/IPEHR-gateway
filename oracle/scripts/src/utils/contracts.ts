import { readFileSync } from "fs"
import * as ethers from "ethers"
import * as path from "path"

const CONTRACTS_DIR = path.join(__dirname, "..", "assets/contracts")

const LINK_TOKEN_CONTRACT_JSON = readFileSync(path.join(CONTRACTS_DIR, "LinkToken.json"), "utf-8")
const ORACLE_CONTRACT_JSON = readFileSync(path.join(CONTRACTS_DIR, "Operator.json"), "utf-8")

const DIRECT_CONSUMER_CONTRACT_JSON = readFileSync(
  path.join(CONTRACTS_DIR, "DirectConsumerContract.json"),
  "utf-8",
)

const STATISTICS_CONSUMER_CONTRACT_JSON = readFileSync(
  path.join(CONTRACTS_DIR, "StatisticsConsumerContract.json"),
  "utf-8",
)

const STATISTICS_CONTRACT_JSON = readFileSync(
  path.join(CONTRACTS_DIR, "StatisticsContract.json"),
  "utf-8",
)

export function getLinkTokenContractFactory(signer: ethers.Signer) {
  return ethers.ContractFactory.fromSolidity(LINK_TOKEN_CONTRACT_JSON, signer)
}

export function getLinkTokenContract(signer: ethers.Signer, address: string) {
  const factory = getLinkTokenContractFactory(signer)
  return factory.attach(address)
}

export function getOracleContractFactory(signer: ethers.Signer) {
  return ethers.ContractFactory.fromSolidity(ORACLE_CONTRACT_JSON, signer)
}

export function getOracleContract(signer: ethers.Signer, address: string) {
  const factory = getOracleContractFactory(signer)
  return factory.attach(address)
}

export function getDirectConsumerContractFactory(signer: ethers.Signer) {
  return ethers.ContractFactory.fromSolidity(DIRECT_CONSUMER_CONTRACT_JSON, signer)
}

export function getDirectConsumerContract(signer: ethers.Signer, address: string) {
  const factory = getDirectConsumerContractFactory(signer)
  return factory.attach(address)
}

export function getStatisticsContractFactory(signer: ethers.Signer) {
  return ethers.ContractFactory.fromSolidity(STATISTICS_CONTRACT_JSON, signer)
}

export function getStatisticsContract(signer: ethers.Signer, address: string) {
  const factory = getStatisticsContractFactory(signer)
  return factory.attach(address)
}

export function getStatisticsConsumerContractFactory(signer: ethers.Signer) {
  return ethers.ContractFactory.fromSolidity(STATISTICS_CONSUMER_CONTRACT_JSON, signer)
}

export function getStatisticsConsumerContract(signer: ethers.Signer, address: string) {
  const factory = getStatisticsConsumerContractFactory(signer)
  return factory.attach(address)
}
