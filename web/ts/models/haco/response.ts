import { IsArray, IsNumber, IsOptional, IsString, ValidateNested } from "class-validator"
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


export class MemStat {
  @IsNumber()
  run!: number

  @IsNumber()
  heap!: number

  @IsNumber()
  sys!: number
}

export class TimeStat {
  @IsNumber()
  run!: number

  @IsOptional()
  @IsNumber()
  moment?: number

  @IsNumber()
  time!: number
}

export class MemData {
  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => MemStat)
  stats!: MemStat[]

  @ValidateNested()
  @Type(() => MemStat)
  start!: MemStat

  @ValidateNested()
  @Type(() => MemStat)
  end!: MemStat
}

export class TimeData {
  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => TimeStat)
  runs!: TimeStat[]

  @ValidateNested()
  @Type(() => TimeStat)
  start!: TimeStat

  @ValidateNested()
  @Type(() => TimeStat)
  end!: TimeStat
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

  @ValidateNested()
  @Type(() => MemData)
  memory!: MemData

  @ValidateNested()
  @Type(() => TimeData)
  time!: TimeData
}