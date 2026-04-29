import { IsArray, IsNumber, IsString, ValidateNested } from "class-validator"
import { Type } from "class-transformer"

export class AvgCoeffs {
  @IsNumber()
  alpha!: number

  @IsNumber()
  beta!: number

  @IsNumber()
  run!: number
}

export class PheromoneItem {
  @IsString()
  source_id!: string

  @IsString()
  target_id!: string

  @IsNumber()
  pheromone!: number
}

export class PheromoneMap {
  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => PheromoneItem)
  items!: PheromoneItem[]
}

export class Path {
  @IsNumber()
  run!: number

  @IsArray()
  path!: string[]

  @IsNumber()
  score!: number
}

export class HacoRunDetailsResponse {
  @IsNumber()
  best_score!: number

  @IsArray()
  best_path!: string[]

  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => AvgCoeffs)
  avg_coeffs!: AvgCoeffs[]

  @ValidateNested()
  @Type(() => PheromoneMap)
  final_pheromone_map!: PheromoneMap

  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => Path)
  best_paths!: Path[]
}