import {
  IsNumber,
  ValidateNested
} from "class-validator"
import { Type } from "class-transformer"

import { Graph } from "../graph.js"

import {
  BaseSelection,
  BestSelection,
  TournamentSelection
} from "./selection.js"

import {
  BaseMutation,
  UniformMutation,
  GaussMutation
} from "./mutation.js"

import { ArithmeticCrossover } from "./crossover.js"
import { BaseOptimisation } from "./optimisation.js"

export class HacoRunRequest {
  @ValidateNested()
  @Type(() => Graph)
  graph!: Graph

  @IsNumber()
  default_alpha!: number

  @IsNumber()
  default_beta!: number

  @IsNumber()
  pheromone_multiplier!: number

  @IsNumber()
  evaporation_rate!: number

  @IsNumber()
  initial_pheromone!: number

  @IsNumber()
  generation_count!: number

  @IsNumber()
  colony_size!: number

  @IsNumber()
  generation_period!: number

  @IsNumber()
  parent_count!: number

  @ValidateNested()
  @Type(() => BaseSelection, {
    discriminator: {
      property: "type",
      subTypes: [
        { value: BestSelection, name: "best" },
        { value: TournamentSelection, name: "tournament" }
      ]
    },
    keepDiscriminatorProperty: true
  })
  selection!: BaseSelection

  @ValidateNested()
  @Type(() => ArithmeticCrossover)
  crossover!: BaseSelection

  @ValidateNested()
  @Type(() => BaseMutation, {
    discriminator: {
      property: "type",
      subTypes: [
        { value: UniformMutation, name: "uniform" },
        { value: GaussMutation, name: "gauss" }
      ]
    },
    keepDiscriminatorProperty: true
  })
  mutation!: BaseMutation

  @ValidateNested()
  @Type(() => BaseOptimisation)
  local_optimisation!: BaseOptimisation
}