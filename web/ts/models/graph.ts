import { IsArray, IsNumber, IsString, ValidateNested, IsIn } from "class-validator"
import { Type } from "class-transformer"

export class Metadata2D {
  @IsNumber()
  x!: number

  @IsNumber()
  y!: number
}

export class Metadata3D {
  @IsNumber()
  x!: number

  @IsNumber()
  y!: number

  @IsNumber()
  z!: number
}

export class MetadataGeo {
  @IsNumber()
  lat!: number

  @IsNumber()
  lon!: number
}

export class GraphNode {
  @IsString()
  id!: string

  @IsString()
  name!: string

  metadata!: any
}

export class GraphEdge {
  @IsString()
  source!: string

  @IsString()
  target!: string

  @IsNumber()
  weight!: number
}

export class Graph {
  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => GraphNode)
  nodes!: GraphNode[]

  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => GraphEdge)
  edges!: GraphEdge[]

  @IsIn(["manhattan_2d", "euclidean_2d", "euclidean_3d", "geo"])
  metadata_type!: string
}