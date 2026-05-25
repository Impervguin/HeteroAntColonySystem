import { IsIn, IsNumber, IsString, ValidateIf } from "class-validator"

export abstract class BaseCrossover {
  @IsString()
  @IsIn(["arithmetic", "sbx", "blx"])
  type!: string
}


export class ArithmeticCrossover extends BaseCrossover {
  type: "arithmetic" = "arithmetic"
}

export class SBXCrossover extends BaseCrossover {
  type: "sbx" = "sbx"

  constructor(eta: number) {
    super()
    this.eta = eta
  }

  @ValidateIf(o => o.type === "sbx")
  @IsNumber()
  eta!: number
}

export class BLXCrossover extends BaseCrossover {
  type: "blx" = "blx"

  constructor(gamma: number) {
    super()
    this.gamma = gamma
  }

  @ValidateIf(o => o.type === "blx")
  @IsNumber()
  gamma!: number
}