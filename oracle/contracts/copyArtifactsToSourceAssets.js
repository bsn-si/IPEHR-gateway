const fs = require("fs/promises");
const path = require("path");

async function main() {
  const sources = [
    "artifacts/contracts/StatisticsConsumer.sol/StatisticsConsumerContract.json",
    "artifacts/contracts/StatisticsContract.sol/StatisticsContract.json",
    "artifacts/contracts/DirectConsumer.sol/DirectConsumerContract.json",
  ];

  const outdir = path.join(__dirname, "..", "scripts/src/assets/contracts");

  const tasks = sources.map((source) =>
    fs.copyFile(
      path.join(__dirname, source),
      path.join(outdir, path.basename(source))
    )
  );

  await Promise.all(tasks);
}

main()
  .then(() => process.exit())
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
