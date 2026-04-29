import { Graph, GraphNode, GraphEdge } from "./models/graph.js"
import { PheromoneMap, Path, AvgCoeffs } from "./models/haco/response.js"
import { HacoRunDetailsResponse } from "./models/haco/response.js"
import { t } from "./i18n.js"

import './shims/global-shim.js'
import * as Plotly from 'plotly.js-dist-min';

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
  const edges = graph.edges

  const xs: number[] = []
  const ys: number[] = []
  const nodeNames: string[] = []

  nodes.forEach(n => {
    const [x, y] = getCoords(n)
    xs.push(x)
    ys.push(y)
    nodeNames.push(n.name || n.id)
  })

  const traces: any[] = []

  // Draw all edges with high transparency (not in path)
  if (edges && edges.length > 0) {
    const edgeXs: number[] = []
    const edgeYs: number[] = []
    
    edges.forEach(edge => {
      const sourceNode = nodes.find(n => n.id === edge.source)
      const targetNode = nodes.find(n => n.id === edge.target)
      
      if (sourceNode && targetNode) {
        const [x1, y1] = getCoords(sourceNode)
        const [x2, y2] = getCoords(targetNode)
        
        edgeXs.push(x1, x2, null as any)
        edgeYs.push(y1, y2, null as any)
      }
    })
    
    traces.push({
      x: edgeXs,
      y: edgeYs,
      mode: "lines",
      type: "scatter",
      name: t("otherEdges"),
      line: {
        color: "rgba(115, 115, 115, 0.3)", // neutral-500 with low opacity
        width: 1.5,
        dash: "solid"
      },
      hoverinfo: "none",
      showlegend: true
    })
  }

  // Draw best path edges (highlighted and cycled)
  if (bestPath && bestPath.length > 0) {
    const pathEdgeXs: number[] = []
    const pathEdgeYs: number[] = []
    
    // Draw edges between consecutive nodes in bestPath
    for (let i = 0; i < bestPath.length - 1; i++) {
      const sourceNode = nodes.find(n => n.id === bestPath[i])
      const targetNode = nodes.find(n => n.id === bestPath[i + 1])
      
      if (sourceNode && targetNode) {
        const [x1, y1] = getCoords(sourceNode)
        const [x2, y2] = getCoords(targetNode)
        
        pathEdgeXs.push(x1, x2, null as any)
        pathEdgeYs.push(y1, y2, null as any)
      }
    }
    
    // Close the cycle: connect last node to first node
    if (bestPath.length > 1) {
      const lastNode = nodes.find(n => n.id === bestPath[bestPath.length - 1])
      const firstNode = nodes.find(n => n.id === bestPath[0])
      
      if (lastNode && firstNode) {
        const [x1, y1] = getCoords(lastNode)
        const [x2, y2] = getCoords(firstNode)
        pathEdgeXs.push(x1, x2, null as any)
        pathEdgeYs.push(y1, y2, null as any)
      }
    }
    
    traces.push({
      x: pathEdgeXs,
      y: pathEdgeYs,
      mode: "lines",
      type: "scatter",
      name: t("bestPathTSP"),
      line: {
        color: "#f59e0b", // amber-500
        width: 3,
        dash: "solid"
      },
      hoverinfo: "none",
      showlegend: true
    })
  }

  // Draw all nodes (all are in path for TSP)
  traces.push({
    x: xs,
    y: ys,
    mode: "markers+text",
    type: "scatter",
    name: t("cities"),
    text: nodeNames,
    textposition: "top center",
    textfont: {
      size: 10,
      color: "#d4d4d8", // neutral-300
      family: "Arial, sans-serif"
    },
    marker: {
      size: 12,
      color: "#3b82f6", // blue-500
      symbol: "circle",
      line: {
        color: "#171717", // neutral-900 (dark background)
        width: 2
      }
    },
    hoverinfo: "text",
    hovertext: nodeNames.map((name, idx) =>
      `<b style="color:#f59e0b">${name}</b><br>ID: ${nodes[idx].id}<br>Coordinates: (${xs[idx].toFixed(2)}, ${ys[idx].toFixed(2)})`
    ),
    showlegend: true
  })

  // Calculate layout with neutral theme colors
  const xRange = Math.max(...xs) - Math.min(...xs)
  const yRange = Math.max(...ys) - Math.min(...ys)
  const padding = Math.max(xRange, yRange) * 0.1

  const layout = {
    title: {
      text: t("tspGraphTitle"),
      font: {
        size: 16,
        family: "Arial, sans-serif",
        weight: "bold" as "bold",
        color: "#fef3c7" // amber-50
      },
      x: 0.5,
      xanchor: "center" as "center"
    },
    xaxis: {
      title: {
        text: t("xCoordinate"),
        font: { color: "#a3a3a3" } // neutral-400
      },
      showgrid: true,
      gridcolor: "#404040", // neutral-700
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252", // neutral-600
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" }, // neutral-400
      range: [Math.min(...xs) - padding, Math.max(...xs) + padding]
    },
    yaxis: {
      title: {
        text: t("yCoordinate"),
        font: { color: "#a3a3a3" } // neutral-400
      },
      showgrid: true,
      gridcolor: "#404040", // neutral-700
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252", // neutral-600
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" }, // neutral-400
      range: [Math.min(...ys) - padding, Math.max(...ys) + padding],
      scaleanchor: "x" as "x",
      scaleratio: 1
    },
    hovermode: "closest" as "closest",
    plot_bgcolor: "#262626", // neutral-800
    paper_bgcolor: "#171717", // neutral-900
    margin: { 
      t: 50,
      r: 80, // Increased to make room for legend
      b: 50,
      l: 50
    },
    legend: {
      x: 1.02,
      y: 1,
      xanchor: "left" as "left",
      bgcolor: "rgba(23, 23, 23, 0.9)", // neutral-900 with opacity
      bordercolor: "#525252", // neutral-600
      borderwidth: 1,
      font: { color: "#d4d4d4" }, // neutral-300
      itemsizing: "constant" as "constant"
    },
    font: {
      family: "Arial, sans-serif",
      size: 12,
      color: "#d4d4d4" // neutral-300
    }
  }

  const config = {
    responsive: true,
    displayModeBar: true,
    modeBarButtonsToRemove: ["lasso2d", "select2d"] as any,
    displaylogo: false,
    toImageButtonOptions: {
      format: "png" as "png",
      filename: "tsp_graph_visualization",
      width: 1200,
      height: 800
    }
  }

  const plotDiv = document.getElementById("graph")!
  Plotly.newPlot(plotDiv, traces, layout, config)
}

// =======================
// 2. PHEROMONE GRAPH
// =======================

export function renderPheromones(
  graph: Graph,
  pheromones: PheromoneMap
) {
  const nodes = graph.nodes
  const edges = graph.edges

  // Calculate max pheromone value for normalization
  const maxPheromone = Math.max(
    ...pheromones.items.map(i => i.pheromone),
    0.001 // Avoid division by zero
  )

  const minPheromone = Math.min(
    ...pheromones.items.map(i => i.pheromone),
    0
  )

  const traces: any[] = []

  // Create a map for quick pheromone lookup
  const pheromoneMap = new Map<string, number>()
  pheromones.items.forEach(item => {
    const key = `${item.source_id}-${item.target_id}`
    pheromoneMap.set(key, item.pheromone)
  })

  // Draw all edges with very low opacity as base layer
  if (edges && edges.length > 0) {
    const edgeXs: number[] = []
    const edgeYs: number[] = []
    
    edges.forEach(edge => {
      const sourceNode = nodes.find(n => n.id === edge.source)
      const targetNode = nodes.find(n => n.id === edge.target)
      
      if (sourceNode && targetNode) {
        const [x1, y1] = getCoords(sourceNode)
        const [x2, y2] = getCoords(targetNode)
        
        edgeXs.push(x1, x2, null as any)
        edgeYs.push(y1, y2, null as any)
      }
    })
    
    traces.push({
      x: edgeXs,
      y: edgeYs,
      mode: "lines",
      type: "scatter",
      name: t("graphEdges"),
      line: {
        color: "rgba(115, 115, 115, 0.15)",
        width: 1
      },
      hoverinfo: "none",
      showlegend: true,
      zorder: 0 // Lowest z-order (background)
    })
  }

  // Sort pheromone items by intensity (ascending) so higher pheromone edges are drawn last (on top)
  const sortedPheromones = [...pheromones.items].sort((a, b) => a.pheromone - b.pheromone)

  // Draw pheromone edges with color intensity based on pheromone level
  sortedPheromones.forEach(e => {
    const s = nodes.find(n => n.id === e.source_id)
    const t = nodes.find(n => n.id === e.target_id)

    if (!s || !t) return

    const [x1, y1] = getCoords(s)
    const [x2, y2] = getCoords(t)

    // Normalize pheromone value (0 to 1)
    const intensity = e.pheromone / maxPheromone
    
    // Color gradient: low (blue) -> medium (amber) -> high (red)
    let color: string
    if (intensity < 0.33) {
      // Blue shades for low pheromone
      const r = 59
      const g = 130
      const b = 246
      color = `rgba(${r}, ${g}, ${b}, ${0.4 + intensity * 0.4})`
    } else if (intensity < 0.66) {
      // Amber shades for medium pheromone
      const t2 = (intensity - 0.33) / 0.33
      const r = Math.floor(245 + t2 * 10)
      const g = Math.floor(158 - t2 * 11)
      const b = Math.floor(11 - t2 * 11)
      color = `rgba(${r}, ${g}, ${b}, ${0.6 + intensity * 0.3})`
    } else {
      // Red shades for high pheromone
      const t2 = (intensity - 0.66) / 0.34
      const r = 239
      const g = Math.floor(68 - t2 * 68)
      const b = Math.floor(68 - t2 * 68)
      color = `rgba(${r}, ${g}, ${b}, ${0.8 + intensity * 0.2})`
    }

    // Line width based on pheromone intensity (2 to 7)
    const lineWidth = 2 + intensity * 5

    traces.push({
      x: [x1, x2],
      y: [y1, y2],
      mode: "lines",
      type: "scatter",
      name: `pheromone: ${e.pheromone.toFixed(4)}`,
      line: {
        color: color,
        width: lineWidth
      },
      opacity: 0.8 + intensity * 0.2,
      hoverinfo: "text",
      hovertext: `<b>Edge:</b> ${s.name || s.id} → ${t.name || t.id}<br>
                  <b>Pheromone level:</b> ${e.pheromone.toFixed(6)}<br>
                  <b>Intensity:</b> ${(intensity * 100).toFixed(1)}%<br>
                  <b>Weight:</b> ${getEdgeWeight(edges, e.source_id, e.target_id)?.toFixed(2) || 'N/A'}`,
      showlegend: false,
      zorder: Math.floor(intensity * 100) // Higher intensity = higher z-order
    })
  })

  // Draw nodes
  const xs: number[] = []
  const ys: number[] = []
  const nodeNames: string[] = []

  nodes.forEach(n => {
    const [x, y] = getCoords(n)
    xs.push(x)
    ys.push(y)
    nodeNames.push(n.name || n.id)
  })

  traces.push({
    x: xs,
    y: ys,
    mode: "markers+text",
    type: "scatter",
    name: t("cities"),
    text: nodeNames,
    textposition: "top center",
    textfont: {
      size: 10,
      color: "#d4d4d4",
      family: "Arial, sans-serif"
    },
    marker: {
      size: 14, // Slightly larger to stand out over edges
      color: "#3b82f6",
      symbol: "circle",
      line: {
        color: "#171717",
        width: 2
      }
    },
    hoverinfo: "text",
    hovertext: nodeNames.map((name, idx) =>
      `<b style="color:#f59e0b">${name}</b><br>ID: ${nodes[idx].id}<br>Coordinates: (${xs[idx].toFixed(2)}, ${ys[idx].toFixed(2)})`
    ),
    showlegend: true,
    zorder: 1000 // Nodes always on top
  })

  // Add a colorbar legend for pheromone intensity
  const colorbarTrace = {
    x: [null],
    y: [null],
    mode: "markers",
    type: "scatter",
    name: t("pheromoneIntensityLegend"),
    marker: {
      colorscale: [
        [0, "#3b82f6"],    // blue (low)
        [0.5, "#f59e0b"],  // amber (medium)
        [1, "#ef4444"]     // red (high)
      ],
      showscale: true,
      colorbar: {
        title: {
          text: t("pheromoneIntensity"),
          font: { color: "#d4d4d4", size: 12 },
          side: "right"
        },
        thickness: 15,
        len: 0.7,
        bgcolor: "rgba(23, 23, 23, 0.8)",
        tickfont: { color: "#a3a3a3" },
        tickformat: ".0%"
      },
      cmin: 0,
      cmax: 1,
      color: [0, 1]
    },
    hoverinfo: "none",
    showlegend: false,
    zorder: 0
  }
  
  traces.push(colorbarTrace)

  // Calculate layout
  const xRange = Math.max(...xs) - Math.min(...xs)
  const yRange = Math.max(...ys) - Math.min(...ys)
  const padding = Math.max(xRange, yRange) * 0.1

  const layout = {
    title: {
      text: t("pheromoneDistributionTitle"),
      font: {
        size: 16,
        family: "Arial, sans-serif",
        weight: "bold" as "bold",
        color: "#fef3c7"
      },
      x: 0.5,
      xanchor: "center" as  "center"
    },
    xaxis: {
      title: {
        text: t("xCoordinate"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" },
      range: [Math.min(...xs) - padding, Math.max(...xs) + padding]
    },
    yaxis: {
      title: {
        text: t("yCoordinate"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" },
      range: [Math.min(...ys) - padding, Math.max(...ys) + padding],
      scaleanchor: "x" as "x",
      scaleratio: 1
    },
    hovermode: "closest" as "closest",
    plot_bgcolor: "#262626",
    paper_bgcolor: "#171717",
    margin: {
      t: 50,
      r: 80, // Space for colorbar
      b: 50,
      l: 50
    },
    legend: {
      x: 1.05,
      y: 1,
      xanchor: "left" as "left",
      bgcolor: "rgba(23, 23, 23, 0.9)",
      bordercolor: "#525252",
      borderwidth: 1,
      font: { color: "#d4d4d4" }
    },
    font: {
      family: "Arial, sans-serif",
      size: 12,
      color: "#d4d4d4"
    }
  }

  const config = {
    responsive: true,
    displayModeBar: true,
    modeBarButtonsToRemove: ["lasso2d", "select2d"] as any,
    displaylogo: false,
    toImageButtonOptions: {
      format: "png" as "png",
      filename: "pheromone_distribution",
      width: 1200,
      height: 800
    }
  }

  const plotDiv = document.getElementById("pheromone-graph")!
  Plotly.newPlot(plotDiv, traces, layout, config)
}

// Helper function to get edge weight
function getEdgeWeight(edges: GraphEdge[], sourceId: string, targetId: string): number | null {
  const edge = edges.find(e => 
    (e.source === sourceId && e.target === targetId) ||
    (e.source === targetId && e.target === sourceId)
  )
  return edge ? edge.weight : null
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

  // Create index map for nodes
  const index = new Map<string, number>()
  const nodeNames: string[] = []
  nodes.forEach((n, i) => {
    index.set(n.id, i)
    nodeNames.push(n.name || n.id)
  })

  // Initialize matrix with zeros
  const matrix: number[][] = Array.from({ length: size }, () =>
    Array(size).fill(0)
  )

  // Fill matrix with pheromone values (directed)
  pheromones.items.forEach(e => {
    const i = index.get(e.source_id)!
    const j = index.get(e.target_id)!
    matrix[i][j] = e.pheromone
  })

  // Find max pheromone value for color scaling
  const maxPheromone = Math.max(
    ...pheromones.items.map(e => e.pheromone),
    0.001
  )

  // Create hover text for each cell
  const hoverText: string[][] = Array.from({ length: size }, () => Array(size).fill(""))
  for (let i = 0; i < size; i++) {
    for (let j = 0; j < size; j++) {
      const value = matrix[i][j]
      if (value > 0) {
        hoverText[i][j] = `<b>From:</b> ${nodeNames[i]}<br>
                           <b>To:</b> ${nodeNames[j]}<br>
                           <b>Pheromone:</b> ${value.toFixed(6)}<br>
                           <b>Normalized:</b> ${((value / maxPheromone) * 100).toFixed(1)}%`
      } else {
        hoverText[i][j] = `<b>From:</b> ${nodeNames[i]}<br>
                           <b>To:</b> ${nodeNames[j]}<br>
                           <b>Pheromone:</b> No direct edge or zero`
      }
    }
  }

  // Determine font size based on number of nodes
  const fontSize = size > 20 ? 8 : size > 10 ? 10 : 12
  const labelRotation = size > 15 ? -45 : 0

  const traces = [{
    z: matrix,
    type: "heatmap" as "heatmap",
    colorscale: [
      [0, "#3b82f6"],      // blue for low values
      [0.33, "#3b82f6"],
      [0.34, "#f59e0b"],   // amber for medium values
      [0.66, "#f59e0b"],
      [0.67, "#ef4444"],   // red for high values
      [1, "#ef4444"]
    ] as any,
    showscale: true,
    colorbar: {
      title: {
        text: t("pheromoneLevel"),
        font: { color: "#d4d4d4", size: 12 },
        side: "right" as "right",
      },
      thickness: 15,
      len: 0.7,
      bgcolor: "rgba(23, 23, 23, 0.8)",
      tickfont: { color: "#a3a3a3" },
      tickformat: ".2e"
    },
    x: nodeNames,
    y: nodeNames,
    xgap: 1,
    ygap: 1,
    hoverinfo: "text" as "text",
    hovertext: hoverText as any,
    hoverongaps: false,
    zmin: 0,
    zmax: maxPheromone,
    zsmooth: false as any // Keep discrete cells for better readability
  }]

  const layout = {
    title: {
      text: t("pheromoneHeatmapTitle"),
      font: {
        size: 16,
        family: "Arial, sans-serif",
        weight: "bold" as "bold",
        color: "#fef3c7"
      },
      x: 0.5,
      xanchor: "center" as "center"
    },
    xaxis: {
      title: {
        text: t("targetNode"),
        font: { color: "#a3a3a3" }
      },
      tickangle: labelRotation,
      tickfont: {
        size: fontSize,
        color: "#d4d4d4"
      },
      tickcolor: "#525252",
      showline: true,
      linecolor: "#525252",
      mirror: true,
      gridcolor: "#404040"
    },
    yaxis: {
      title: {
        text: t("sourceNode"),
        font: { color: "#a3a3a3" }
      },
      tickfont: {
        size: fontSize,
        color: "#d4d4d4"
      },
      tickcolor: "#525252",
      showline: true,
      linecolor: "#525252",
      mirror: true,
      gridcolor: "#404040",
      autorange: "reversed" as "reversed" // So that first node is at top
    },
    hovermode: "closest" as "closest",
    plot_bgcolor: "#262626",
    paper_bgcolor: "#171717",
    margin: {
      t: 60,
      r: 100, // Space for colorbar
      b: size > 15 ? 120 : 80,
      l: size > 15 ? 120 : 80
    },
    autosize: true, // Enable automatic sizing
    font: {
      family: "Arial, sans-serif",
      size: 12,
      color: "#d4d4d4"
    }
  }

  // Add annotation for max value cell if not too many nodes
  if (size <= 20) {
    let maxValue = 0
    let maxI = -1, maxJ = -1
    for (let i = 0; i < size; i++) {
      for (let j = 0; j < size; j++) {
        if (matrix[i][j] > maxValue) {
          maxValue = matrix[i][j]
          maxI = i
          maxJ = j
        }
      }
    }
  }

  const config = {
    responsive: true,
    displayModeBar: true,
    modeBarButtonsToRemove: ["lasso2d", "select2d", "zoomIn2d", "zoomOut2d", "autoScale2d"] as any,
    displaylogo: false,
    toImageButtonOptions: {
      format: "png" as "png",
      filename: "pheromone_heatmap",
      width: 1200,
      height: 800
    }
  }

  const plotDiv = document.getElementById("heatmap")!
  Plotly.newPlot(plotDiv, traces, layout, config)

}
// =======================
// 4. ALPHA / BETA
// =======================

export function renderCoeffs(coeffs: AvgCoeffs[]) {
  const x = coeffs.map(c => c.run)

  // Find min and max values for better y-axis scaling
  const allValues = [...coeffs.map(c => c.alpha), ...coeffs.map(c => c.beta)]
  const minValue = Math.min(...allValues)
  const maxValue = Math.max(...allValues)
  const padding = (maxValue - minValue) * 0.1

  const traces : any[] = [
    {
      x: x,
      y: coeffs.map(c => c.alpha),
      name: t("alphaLabel"),
      type: "scatter" as "scatter",
      mode: "lines+markers" as "lines+markers",
      line: {
        color: "#3b82f6", // blue-500
        width: 2.5,
        shape: "linear" as "linear"
      },
      marker: {
        size: 6,
        color: "#3b82f6",
        symbol: "circle" as "circle",
        line: {
          color: "#171717",
          width: 1
        }
      },
      hoverinfo: "text" as "text",
      hovertext: coeffs.map(c =>
        `<b>${t("alphaLabel")}</b><br>
          Run: ${c.run}<br>
          Value: ${c.alpha.toFixed(4)}<br>`
      ),
      showlegend: true
    },
    {
      x: x,
      y: coeffs.map(c => c.beta),
      name: t("betaLabel"),
      type: "scatter" as "scatter",
      mode: "lines+markers" as "lines+markers",
      line: {
        color: "#f59e0b", // amber-500
        width: 2.5,
        shape: "linear" as "linear"
      },
      marker: {
        size: 6,
        color: "#f59e0b",
        symbol: "diamond" as "diamond",
        line: {
          color: "#171717",
          width: 1
        }
      },
      hoverinfo: "text" as "text",
      hovertext: coeffs.map(c =>
        `<b>${t("betaLabel")}</b><br>
          Run: ${c.run}<br>
          Value: ${c.beta.toFixed(4)}<br>`
      ),
      showlegend: true
    }
  ]

  // Calculate layout
  const layout = {
    title: {
      text: t("parameterEvolutionTitle"),
      font: {
        size: 16,
        family: "Arial, sans-serif",
        weight: "bold" as "bold",
        color: "#fef3c7"
      },
      x: 0.5,
      xanchor: "center" as "center"
    },
    xaxis: {
      title: {
        text: t("generationRun"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" },
      tickformat: "d", // Integer format
      dtick: Math.ceil((Math.max(...x) - Math.min(...x)) / 10) // Auto tick spacing
    },
    yaxis: {
      title: {
        text: t("parameterValue"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: true,
      zerolinecolor: "#525252",
      zerolinewidth: 1,
      tickfont: { color: "#a3a3a3" },
      range: [Math.max(0, minValue - padding), maxValue + padding]
    },
    hovermode: "closest" as "closest",
    plot_bgcolor: "#262626",
    paper_bgcolor: "#171717",
    margin: {
      t: 50,
      r: 50,
      b: 50,
      l: 60
    },
    legend: {
      x: 0.02,
      y: 0.98,
      xanchor: "left" as "left",
      yanchor: "top" as "top",
      bgcolor: "rgba(23, 23, 23, 0.9)",
      bordercolor: "#525252",
      borderwidth: 1,
      font: { color: "#d4d4d4" },
      itemsizing: "constant" as "constant"
    },
    font: {
      family: "Arial, sans-serif",
      size: 12,
      color: "#d4d4d4"
    },
    annotations: [
      {
        x: 0.5,
        y: 1.02,
        xref: "paper" as "paper",
        yref: "paper" as "paper",
        text: t("alphaBetaHint"),
        showarrow: false,
        font: {
          size: 10,
          color: "#a3a3a3"
        }
      }
    ]
  }

  const config = {
    responsive: true,
    displayModeBar: true,
    modeBarButtonsToRemove: ["lasso2d", "select2d"] as any,
    displaylogo: false,
    toImageButtonOptions: {
      format: "png" as "png",
      filename: "parameter_evolution",
      width: 1200,
      height: 800
    }
  }

  const plotDiv = document.getElementById("coeffs")!
  Plotly.newPlot(plotDiv, traces, layout, config)
}

// =======================
// 5. BEST PATH SCORE
// =======================

export function renderBestScore(paths: Path[]) {
  // Sort by run number to ensure correct order
  const sortedPaths = [...paths].sort((a, b) => a.run - b.run)
  
  const runs = sortedPaths.map(p => p.run)
  const scores = sortedPaths.map(p => p.score)
  
  // Find best score (minimum for TSP)
  const bestScore = Math.min(...scores)
  const bestRun = runs[scores.indexOf(bestScore)]
  const worstScore = Math.max(...scores)
  const finalScore = scores[scores.length - 1]
  const initialScore = scores[0]
  
  // Calculate improvement percentage
  const improvement = ((initialScore - bestScore) / initialScore * 100).toFixed(1)
  
  // Create traces
  const traces : any[] = [
    {
      x: runs,
      y: scores,
      name: t("bestPathScore"),
      type: "scatter" as "scatter",
      mode: "lines+markers" as "markers" | "lines+markers",
      line: {
        color: "#10b981", // emerald-500
        width: 2.5,
        shape: "spline" as "spline"
      },
      marker: {
        size: 6,
        color: "#10b981",
        symbol: "circle" as "star" | "circle",
        line: {
          color: "#171717",
          width: 1
        }
      },
      hoverinfo: "text" as "text",
      hovertext: scores.map((score, idx) =>
        `<b>Generation ${runs[idx]}</b><br>
          Score: ${score.toFixed(2)}<br>
          ${idx > 0 ? `Improvement: ${((scores[idx-1] - score) / scores[idx-1] * 100).toFixed(2)}%` : 'Initial generation'}<br>
          ${score === bestScore ? 'Best score' : ''}`
      ),
      showlegend: true,
      fill: "tozeroy" as "tozeroy",
      fillcolor: "rgba(16, 185, 129, 0.1)" // emerald with low opacity
    }
  ]
  
  // Add marker for best score point
  traces.push({
    x: [bestRun],
    y: [bestScore],
    name: t("globalBest"),
    type: "scatter",
    mode: "markers" as "markers" | "lines+markers",
    marker: {
      size: 12,
      color: "#ef4444",
      symbol: "star" as "star" | "circle",
      line: {
        color: "#171717",
        width: 2
      }
    },
    hoverinfo: "text",
    hovertext: [`<b>${t("globalBest")}</b><br>
                  Generation: ${bestRun}<br>
                  Score: ${bestScore.toFixed(2)}<br>
                  Improvement: ${improvement}% from start`],
    showlegend: true
  })
  
  // Add annotation for best score
  const annotations = [
    {
      x: bestRun,
      y: bestScore,
      text: `Best: ${bestScore.toFixed(2)}`,
      showarrow: true,
      arrowhead: 2,
      arrowsize: 1,
      arrowwidth: 1,
      arrowcolor: "#ef4444",
      ax: 0,
      ay: -30,
      bgcolor: "rgba(23, 23, 23, 0.8)",
      bordercolor: "#ef4444",
      borderwidth: 1,
      font: { size: 11, color: "#fef3c7" }
    }
  ]
  
  // Add annotation for final score if it's not the best
  if (finalScore > bestScore) {
    annotations.push({
      x: runs[runs.length - 1],
      y: finalScore,
      text: `Final: ${finalScore.toFixed(2)}`,
      showarrow: true,
      arrowhead: 2,
      arrowsize: 1,
      arrowwidth: 1,
      arrowcolor: "#10b981",
      ax: 0,
      ay: 20,
      bgcolor: "rgba(23, 23, 23, 0.8)",
      bordercolor: "#10b981",
      borderwidth: 1,
      font: { size: 11, color: "#d4d4d4" }
    })
  }
  
  // // Add improvement annotation
  // annotations.push({
  //   x: 0.5,
  //   y: 1.05,
  //   xref: "paper" as "paper",
  //   yref: "paper" as "paper",
  //   text: `📈 Total improvement: ${improvement}% (${initialScore.toFixed(2)} → ${bestScore.toFixed(2)})`,
  //   showarrow: false,
  //   font: {
  //     size: 11,
  //     color: "#fef3c7"
  //   },
  //   bgcolor: "rgba(16, 185, 129, 0.2)",
  //   bordercolor: "#10b981",
  //   borderwidth: 1,
  //   borderpad: 4
  // })
  
  // Calculate y-axis range with padding
  const yRange = worstScore - bestScore
  const yPadding = yRange * 0.1
  
  const layout = {
    title: {
      text: t("bestPathEvolutionTitle"),
      font: {
        size: 16,
        family: "Arial, sans-serif",
        weight: "bold" as "bold",
        color: "#fef3c7"
      },
      x: 0.5,
      xanchor: "center" as "center"
    },
    xaxis: {
      title: {
        text: t("generation"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: false,
      tickfont: { color: "#a3a3a3" },
      tickformat: "d",
      dtick: Math.ceil((Math.max(...runs) - Math.min(...runs)) / 10)
    },
    yaxis: {
      title: {
        text: t("pathScoreLowerBetter"),
        font: { color: "#a3a3a3" }
      },
      showgrid: true,
      gridcolor: "#404040",
      gridwidth: 0.5,
      showline: true,
      linecolor: "#525252",
      mirror: true,
      zeroline: true,
      zerolinecolor: "#525252",
      zerolinewidth: 1,
      tickfont: { color: "#a3a3a3" },
      range: [Math.max(0, bestScore - yPadding), worstScore + yPadding]
    },
    hovermode: "closest" as "closest",
    plot_bgcolor: "#262626",
    paper_bgcolor: "#171717",
    margin: {
      t: 80, // Extra space for improvement annotation
      r: 50,
      b: 50,
      l: 60
    },
    legend: {
      x: 0.02,
      y: 0.98,
      xanchor: "left" as "left",
      yanchor: "top" as "top",
      bgcolor: "rgba(23, 23, 23, 0.9)",
      bordercolor: "#525252",
      borderwidth: 1,
      font: { color: "#d4d4d4" }
    },
    font: {
      family: "Arial, sans-serif",
      size: 12,
      color: "#d4d4d4"
    },
    annotations: annotations,
    shapes: [
      {
        type: "rect" as "rect",
        xref: "paper" as "paper",
        yref: "paper" as "paper",
        x0: 0,
        y0: 0,
        x1: 1,
        y1: 1,
        fillcolor: "rgba(16, 185, 129, 0.02)",
        layer: "below" as "below",
        line: { width: 0 }
      }
    ]
  }
  
  const config = {
    responsive: true,
    displayModeBar: true,
    modeBarButtonsToRemove: ["lasso2d", "select2d"] as any,
    displaylogo: false,
    toImageButtonOptions: {
      format: "png" as "png",
      filename: "best_score_evolution",
      width: 1200,
      height: 800
    }
  }
  
  const plotDiv = document.getElementById("score")!
  Plotly.newPlot(plotDiv, traces, layout, config)
  
  // Style modebar to appear only on hover
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