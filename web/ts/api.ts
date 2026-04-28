import {
  Graph,
  ParseTSPResponse,
  GetTSPResponse,
  HacoRunDetailsResponse,
  HacoRunRequest
} from "./types.js"

// =======================
// Helpers
// =======================

async function safeJson(res: Response): Promise<any> {
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`HTTP ${res.status}: ${text}`)
  }
  return await res.json()
}

function getApiBase(): string {
  return window.APP_CONFIG.apiBase ?? "http://localhost:8080/api/v1"
}

// =======================
// Runtime validation (минимальная)
// =======================

function assertGraph(obj: any): Graph {
  if (!obj || !Array.isArray(obj.nodes) || !Array.isArray(obj.edges)) {
    throw new Error("Invalid Graph shape")
  }
  return obj as Graph
}

function assertParseTSPResponse(obj: any): ParseTSPResponse {
  if (!obj || !obj.graph) {
    throw new Error("Invalid ParseTSPResponse")
  }
  return {
    graph: assertGraph(obj.graph)
  }
}

function assertGetTSPResponse(obj: any): GetTSPResponse {
  if (!obj || !obj.graph) {
    throw new Error("Invalid GetTSPResponse")
  }
  return {
    graph: assertGraph(obj.graph)
  }
}

function assertHacoRunDetailsResponse(obj: any): HacoRunDetailsResponse {
  if (
    !obj ||
    typeof obj.best_score !== "number" ||
    !Array.isArray(obj.best_path) ||
    !Array.isArray(obj.avg_coeffs) ||
    !obj.final_pheromone_map ||
    !Array.isArray(obj.best_paths)
  ) {
    throw new Error("Invalid HacoRunDetailsResponse")
  }

  return obj as HacoRunDetailsResponse
}

// =======================
// 1. TSP API
// =======================

// POST /api/v1/tsp/parse
export async function parseTSP(file: File): Promise<ParseTSPResponse> {
  const formData = new FormData()
  formData.append("file", file)

  

  const res = await fetch(`${getApiBase()}/tsp/parse`, {
    method: "POST",
    body: formData
  })

  const json = await safeJson(res)
  return assertParseTSPResponse(json)
}

// GET /api/v1/tsp/:file
export async function getTSP(filename: string): Promise<GetTSPResponse> {
  const res = await fetch(`${getApiBase()}/tsp/${encodeURIComponent(filename)}`, {
    method: "GET"
  })

  const json = await safeJson(res)
  return assertGetTSPResponse(json)
}

// =======================
// 2. HACO REQUEST BUILDER
// =======================

export type HacoFormInput = {
  default_alpha: number
  default_beta: number
  pheromone_multiplier: number
  evaporation_rate: number
  initial_pheromone: number

  generation_count: number
  colony_size: number
  generation_period: number
  parent_count: number

  selection: HacoRunRequest["selection"]
  crossover: HacoRunRequest["crossover"]
  mutation: HacoRunRequest["mutation"]
  local_optimisation: HacoRunRequest["local_optimisation"]
}

// сборка строго под Go DTO
export function buildHacoRequest(
  graph: Graph,
  input: HacoFormInput
): HacoRunRequest {
  return {
    graph,

    default_alpha: Number(input.default_alpha),
    default_beta: Number(input.default_beta),
    pheromone_multiplier: Number(input.pheromone_multiplier),
    evaporation_rate: Number(input.evaporation_rate),
    initial_pheromone: Number(input.initial_pheromone),

    generation_count: Number(input.generation_count),
    colony_size: Number(input.colony_size),
    generation_period: Number(input.generation_period),
    parent_count: Number(input.parent_count),

    selection: input.selection,
    crossover: input.crossover,
    mutation: input.mutation,
    local_optimisation: input.local_optimisation
  }
}

// =======================
// 3. HACO API
// =======================

// POST /api/v1/haco/run/details
export async function runHacoDetails(
  request: HacoRunRequest
): Promise<HacoRunDetailsResponse> {
  const res = await fetch(`${getApiBase()}/haco/run/details`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(request)
  })

  const json = await safeJson(res)
  return assertHacoRunDetailsResponse(json)
}