import * as express from "express"
import { Request } from "express"
import * as morgan from "morgan"

const app = express()
const port = 3000

interface Params {
  period: string
}

function randomUint(_min: number, _max: number): number {
  const min = Math.ceil(_min),
    max = Math.floor(_max)
  return Math.floor(Math.random() * (max - min + 1)) + min
}

app.use(morgan("combined"))

app.get("/", (req, res) => {
  const now = new Date()

  return res.status(200).json({
    type: "LATEST",

    data: {
      documents: randomUint(100, 1000),
      patients: randomUint(100, 1000),
      time: now.getTime(),
    },

    month: {
      documents: randomUint(100, 1000),
      patients: randomUint(100, 1000),
      time: parseInt(`${now.getFullYear()}${now.getMonth() + 1}`),
    },
  })
})

app.get("/:period", ({ params: { period } }: Request<Params>, res) => {
  const time = parseInt(period)

  if (isNaN(time)) {
    return res.status(403).json({ error: "Invalid period" })
  }

  return res.status(200).json({
    type: "PERIOD",
    data: {
      documents: randomUint(100, 1000),
      patients: randomUint(100, 1000),
      time,
    },
  })
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
