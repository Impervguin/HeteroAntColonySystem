import { ValidateNested } from "class-validator"
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