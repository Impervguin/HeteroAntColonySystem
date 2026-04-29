import { IsIn, IsString } from "class-validator"

export abstract class BaseCrossover {
  @IsString()
  @IsIn(["arithmetic"])
  type!: string
}


export class ArithmeticCrossover extends BaseCrossover {
  type: "arithmetic" = "arithmetic"
}