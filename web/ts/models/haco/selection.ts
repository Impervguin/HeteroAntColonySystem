import { IsIn, IsNumber, IsString, ValidateIf } from "class-validator"

export abstract class BaseSelection {
  @IsString()
  @IsIn(["best", "tournament", "roulette"])
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

export class RouletteSelection extends BaseSelection {
  type: "roulette" = "roulette"
}