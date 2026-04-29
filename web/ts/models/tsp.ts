import { ValidateNested, IsArray, IsString, IsNumber } from "class-validator"
import { Type } from "class-transformer"
import { Graph } from "./graph.js"

export class ParseTSPResponse {
  @ValidateNested()
  @Type(() => Graph)
  graph!: Graph
}

export class GetTSPResponse {
  @ValidateNested()
  @Type(() => Graph)
  graph!: Graph
}

export class ListTSPResponse {
  @IsArray()
  @IsString()
  files!: string[]
}

export class GraphStatsRequest {
  @ValidateNested()
  @Type(() => Graph)
  graph!: Graph
}

export class GraphStatsResponse {
  @IsNumber()
  nodes_count!: number

  @IsNumber()
  edges_count!: number

  @IsNumber()
  possible_solutions!: number

  @IsNumber()
  avg_edge_weight!: number

  @IsNumber()
  max_edge_weight!: number

  @IsNumber()
  min_edge_weight!: number

  @IsNumber()
  expected_path_length!: number

  @IsNumber()
  recommended_pheromone_multiplier!: number

  @IsNumber()
  recommended_evaporation_rate!: number
}