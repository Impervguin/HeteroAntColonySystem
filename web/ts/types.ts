// =======================
// GRAPH (общий)
// =======================

export type MetadataType =
  | "manhattan_2d"
  | "euclidean_2d"
  | "euclidean_3d"
  | "geo"

// --- Metadata variants ---

export type Metadata2D = {
  x: number
  y: number
}

export type Metadata3D = {
  x: number
  y: number
  z: number
}

export type MetadataGeo = {
  lat: number
  lon: number
}

export type NodeMetadata =
  | Metadata2D
  | Metadata3D
  | MetadataGeo
  | Record<string, any>

export function is2D(m: any): m is Metadata2D {
  return m && "x" in m && "y" in m && !("z" in m)
}

export function is3D(m: any): m is Metadata3D {
  return m && "x" in m && "y" in m && "z" in m
}

export function isGeo(m: any): m is MetadataGeo {
  return m && "lat" in m && "lon" in m
}

// --- Node / Edge ---

export type GraphNode = {
  id: string
  name: string
  metadata: NodeMetadata
}

export type GraphEdge = {
  source: string
  target: string
  weight: number
}

// --- Graph ---

export type Graph = {
  nodes: GraphNode[]
  edges: GraphEdge[]
  metadata_type: MetadataType
}


export type ParseTSPResponse = {
  graph: Graph
}

export type GetTSPResponse = {
  graph: Graph
}

// HACO Response

export type AvgCoeffs = {
  alpha: number
  beta: number
  run: number
}

export type PheromoneItem = {
  source_id: string
  target_id: string
  pheromone: number
}

export type PheromoneMap = {
  items: PheromoneItem[]
}

export type Path = {
  run: number
  path: string[]
  score: number
}

export type HacoRunDetailsResponse = {
  best_score: number
  best_path: string[]

  avg_coeffs: AvgCoeffs[]
  final_pheromone_map: PheromoneMap
  best_paths: Path[]
}

// HACO Request

export type HacoSelectionStrategy =
  | { type: "best" }
  | { type: "tournament"; k: number }

export type HacoCrossoverStrategy = {
  type: "arithmetic"
}

export type HacoMutationStrategy =
  | { type: "uniform"; min: number; max: number }
  | { type: "gauss"; mean: number; std: number }

export type HacoLocalOptimisationStrategy = {
  type: "noop" | "2opt"
}

export type HacoRunRequest = {
  graph: Graph

  default_alpha: number
  default_beta: number
  pheromone_multiplier: number
  evaporation_rate: number
  initial_pheromone: number

  generation_count: number
  colony_size: number
  generation_period: number
  parent_count: number

  selection: HacoSelectionStrategy
  crossover: HacoCrossoverStrategy
  mutation: HacoMutationStrategy
  local_optimisation: HacoLocalOptimisationStrategy
}