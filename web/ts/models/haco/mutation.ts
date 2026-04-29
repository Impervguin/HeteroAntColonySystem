import { IsIn, IsNumber, IsString, ValidateIf } from "class-validator"

export abstract class BaseMutation {
  @IsString()
  @IsIn(["uniform", "gauss"])
  type!: string
}

export class UniformMutation extends BaseMutation {
  type: "uniform" = "uniform"

  constructor(min: number, max: number) {
    super()
    this.min = min
    this.max = max
  }

  @ValidateIf(o => o.type === "uniform")
  @IsNumber()
  min!: number

  @ValidateIf(o => o.type === "uniform")
  @IsNumber()
  max!: number
}

export class GaussMutation extends BaseMutation {
  type: "gauss" = "gauss"

  constructor(mean: number, std: number) {
    super()
    this.mean = mean
    this.std = std
  }

  @ValidateIf(o => o.type === "gauss")
  @IsNumber()
  mean!: number

  @ValidateIf(o => o.type === "gauss")
  @IsNumber()
  std!: number
}