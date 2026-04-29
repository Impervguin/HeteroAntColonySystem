import { IsNumber } from "class-validator";
import { GraphStatsResponse } from "./tsp";
import { HacoRunDetailsResponse } from "./haco/response";

export type GraphStatsRenderRequest = GraphStatsResponse;

export class RuntimeStatsRenderRequest {
  @IsNumber()
  score!: number

  @IsNumber()
  seen_paths!: number

  @IsNumber()
  total_time!: number

  @IsNumber()
  min_time!: number

  @IsNumber()
  avg_time!: number

  @IsNumber()
  max_time!: number

  @IsNumber()
  min_memory!: number

  @IsNumber()
  avg_memory!: number

  @IsNumber()
  max_memory!: number
}

export function fromRunDetailsResponse(response: HacoRunDetailsResponse, seenPaths: number): RuntimeStatsRenderRequest {
    let avgMemory: number = 0
    let maxMemory: number = response.memory.stats[0].sys
    let minMemory: number = response.memory.stats[0].sys
    const startMemory = response.memory.start.sys
    for (const stat of response.memory.stats) {
        const mem: number = stat.sys - startMemory
        avgMemory += mem
        if (mem > maxMemory) {
            maxMemory = mem
        }
        if (mem < minMemory) {
            minMemory = mem
        }
    }
    avgMemory /= response.memory.stats.length
    avgMemory = Math.round(avgMemory)

    let avgTime: number = 0
    let maxTime: number = response.time.runs[0].time
    let minTime: number = response.time.runs[0].time
    for (const stat of response.time.runs) {
        avgTime += stat.time
        if (stat.time > maxTime) {
            maxTime = stat.time
        }
        if (stat.time < minTime) {
            minTime = stat.time
        }
    }
    avgTime /= response.time.runs.length
    const totalTime: number = response.time.end.moment! - response.time.start.moment!
 

    return {
        score: response.best_score,
        seen_paths: seenPaths,
        total_time: totalTime,
        min_time: minTime,
        avg_time: avgTime,
        max_time: maxTime,
        min_memory: minMemory,
        avg_memory: avgMemory,
        max_memory: maxMemory
    }
}