import { IsIn, IsNumber, IsString, ValidateIf } from "class-validator"

export abstract class BaseSelection {
  @IsString()
  @IsIn(["best", "tournament"])
  type!: string
}

export class BestSelection extends BaseSelection {
  type: "best" = "best"
}

export class TournamentSelection extends BaseSelection {
  type: "tournament" = "tournament"

  constructor(k: number) {
    super()
    this.k = k
  }

  @ValidateIf(o => o.type === "tournament")
  @IsNumber()
  k!: number
}