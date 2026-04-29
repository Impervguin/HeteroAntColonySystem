import { IsIn, IsString } from "class-validator"

export abstract class BaseOptimisation {
  @IsString()
  @IsIn(["noop", "2opt"])
  type!: "noop" | "2opt"
}

export class NoOpLocalOptimisation extends BaseOptimisation {
  type: "noop" = "noop"
}

export class TwoOptLocalOptimisation extends BaseOptimisation {
  type: "2opt" = "2opt"
}