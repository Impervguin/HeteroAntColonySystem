import { Graph, GraphNode } from "./models/graph.js"
import { PheromoneMap, Path, AvgCoeffs } from "./models/haco/response.js"
import { HacoRunDetailsResponse } from "./models/haco/response.js"

// =======================
// Metadata utils
// =======================

function getCoords(node: GraphNode): [number, number] {
  const m: any = node.metadata

  if ("x" in m && "y" in m) {
    return [m.x, m.y]
  }

  if ("lat" in m && "lon" in m) {
    return [m.lon, m.lat] // lon = x, lat = y
  }

  throw new Error("Unsupported metadata format")
}

function is3D(node: GraphNode): boolean {
  const m: any = node.metadata
  return "z" in m
}

// =======================
// 1. GRAPH (nodes + path)
// =======================

export function renderGraph(
  graph: Graph,
  bestPath?: string[]
) {
  const nodes = graph.nodes

  const xs: number[] = []
  const ys: number[] = []

  nodes.forEach(n => {
    const [x, y] = getCoords(n)
    xs.push(x)
    ys.push(y)
  })

  const traces: any[] = [
    {
      x: xs,
      y: ys,
      mode: "markers",
      type: "scatter",
      name: "nodes"
    }
  ]

  if (bestPath && bestPath.length > 0) {
    const pathX: number[] = []
    const pathY: number[] = []

    bestPath.forEach(id => {
      const node = nodes.find(n => n.id === id)!
      const [x, y] = getCoords(node)
      pathX.push(x)
      pathY.push(y)
    })

    traces.push({
      x: pathX,
      y: pathY,
      mode: "lines",
      type: "scatter",
      name: "best path"
    })
  }

  Plotly.newPlot("graph", traces, {
    margin: { t: 20 }
  })
}

// =======================
// 2. PHEROMONE GRAPH
// =======================

export function renderPheromones(
  graph: Graph,
  pheromones: PheromoneMap
) {
  const nodes = graph.nodes

  const max = Math.max(
    ...pheromones.items.map(i => i.pheromone),
    1
  )

  const traces: any[] = []

  pheromones.items.forEach(e => {
    const s = nodes.find(n => n.id === e.source_id)!
    const t = nodes.find(n => n.id === e.target_id)!

    const [x1, y1] = getCoords(s)
    const [x2, y2] = getCoords(t)

    traces.push({
      x: [x1, x2],
      y: [y1, y2],
      mode: "lines",
      type: "scatter",
      line: {
        width: 2
      },
      opacity: e.pheromone / max
    })
  })

  Plotly.newPlot("pheromone-graph", traces, {
    margin: { t: 20 }
  })
}

// =======================
// 3. HEATMAP
// =======================

export function renderHeatmap(
  graph: Graph,
  pheromones: PheromoneMap
) {
  const nodes = graph.nodes
  const size = nodes.length

  const index = new Map<string, number>()
  nodes.forEach((n, i) => index.set(n.id, i))

  const matrix: number[][] = Array.from({ length: size }, () =>
    Array(size).fill(0)
  )

  pheromones.items.forEach(e => {
    const i = index.get(e.source_id)!
    const j = index.get(e.target_id)!
    matrix[i][j] = e.pheromone
  })

  Plotly.newPlot("heatmap", [
    {
      z: matrix,
      type: "heatmap"
    }
  ])
}

// =======================
// 4. ALPHA / BETA
// =======================

export function renderCoeffs(coeffs: AvgCoeffs[]) {
  const x = coeffs.map(c => c.run)

  Plotly.newPlot("coeffs", [
    {
      x,
      y: coeffs.map(c => c.alpha),
      name: "alpha",
      type: "scatter"
    },
    {
      x,
      y: coeffs.map(c => c.beta),
      name: "beta",
      type: "scatter"
    }
  ])
}

// =======================
// 5. BEST PATH SCORE
// =======================

export function renderBestScore(paths: Path[]) {
  Plotly.newPlot("score", [
    {
      x: paths.map(p => p.run),
      y: paths.map(p => p.score),
      type: "scatter",
      name: "best score"
    }
  ])
}

// =======================
// 6. FULL PIPELINE
// =======================

export function renderAll(
  graph: Graph,
  result: HacoRunDetailsResponse
) {
  renderGraph(graph, result.best_path)

  renderPheromones(graph, result.final_pheromone_map)

  renderHeatmap(graph, result.final_pheromone_map)

  renderCoeffs(result.avg_coeffs)

  renderBestScore(result.best_paths)
}