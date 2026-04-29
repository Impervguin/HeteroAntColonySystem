import {
  parseTSP,
  getTSP,
  runHacoDetails,
  graphStats,
  renderGraphStats
} from "./api.js"

import {
  renderGraph,
  renderAll
} from "./render.js"

import { HacoRunRequest } from "./models/haco/request.js"
import { GraphStatsRequest, GraphStatsResponse } from "./models/tsp.js"
import { HacoRunDetailsResponse } from "./models/haco/response.js"
import { Graph } from "./models/graph.js"
import { HacoFormInput, buildHacoRequest } from "./form.js"

import './shims/global-shim.js'
import * as Plotly from 'plotly.js-dist-min';

// =======================
// State Management
// =======================

let currentGraph: Graph | null = null
let currentResult: HacoRunDetailsResponse | null = null
let currentFilename: string | null = null

// =======================
// DOM Elements
// =======================

const elements = {
  fileInput: document.getElementById("file-input") as HTMLInputElement,
  graphSelect: document.getElementById("graph-select") as HTMLSelectElement,
  runBtn: document.getElementById("run-btn") as HTMLButtonElement,
  
  // Basic params
  defaultAlpha: document.getElementById("default-alpha") as HTMLInputElement,
  defaultBeta: document.getElementById("default-beta") as HTMLInputElement,
  pheromoneMultiplier: document.getElementById("pheromone-multiplier") as HTMLInputElement,
  evaporationRate: document.getElementById("evaporation-rate") as HTMLInputElement,
  initialPheromone: document.getElementById("initial-pheromone") as HTMLInputElement,
  generationCount: document.getElementById("generation-count") as HTMLInputElement,
  colonySize: document.getElementById("colony-size") as HTMLInputElement,
  generationPeriod: document.getElementById("generation-period") as HTMLInputElement,
  parentCount: document.getElementById("parent-count") as HTMLInputElement,
  
  // Selection
  selectionType: document.getElementById("selection-type") as HTMLSelectElement,
  tournamentK: document.getElementById("tournament-k") as HTMLInputElement,
  
  // Mutation
  mutationType: document.getElementById("mutation-type") as HTMLSelectElement,
  mutationMin: document.getElementById("mutation-min") as HTMLInputElement,
  mutationMax: document.getElementById("mutation-max") as HTMLInputElement,
  mutationMean: document.getElementById("mutation-mean") as HTMLInputElement,
  mutationStd: document.getElementById("mutation-std") as HTMLInputElement,
  
  // Crossover
  crossoverType: document.getElementById("crossover-type") as HTMLSelectElement,

  // Local optimisation
  localOptimisation: document.getElementById("local-optimisation") as HTMLSelectElement,
}

// =======================
// Helper Functions
// =======================

function showError(message: string) {
  console.error(message)
  alert(`Error: ${message}`)
}

function disableRunButton(disabled: boolean) {
  elements.runBtn.disabled = disabled
  elements.runBtn.textContent = disabled ? "Running..." : "Run HACO"
}

function getFormInput(): HacoFormInput | null {
  try {
    const selectionType = elements.selectionType.value
    const mutationType = elements.mutationType.value
    
    return {
      default_alpha: parseFloat(elements.defaultAlpha.value),
      default_beta: parseFloat(elements.defaultBeta.value),
      pheromone_multiplier: parseFloat(elements.pheromoneMultiplier.value),
      evaporation_rate: parseFloat(elements.evaporationRate.value),
      initial_pheromone: parseFloat(elements.initialPheromone.value),
      
      generation_count: parseInt(elements.generationCount.value, 10),
      colony_size: parseInt(elements.colonySize.value, 10),
      generation_period: parseInt(elements.generationPeriod.value, 10),
      parent_count: parseInt(elements.parentCount.value, 10),
      
      selection_type: selectionType as "best" | "tournament",
      tournament_k: selectionType === "tournament"
        ? parseInt(elements.tournamentK.value, 10)
        : undefined,
      
      crossover_type: elements.crossoverType.value as "arithmetic",
      
      mutation_type: mutationType as "uniform" | "gauss",
      mutation_min: mutationType === "uniform"
        ? parseFloat(elements.mutationMin.value)
        : undefined,
      mutation_max: mutationType === "uniform"
        ? parseFloat(elements.mutationMax.value)
        : undefined,
      mutation_mean: mutationType === "gauss"
        ? parseFloat(elements.mutationMean.value)
        : undefined,
      mutation_std: mutationType === "gauss"
        ? parseFloat(elements.mutationStd.value)
        : undefined,   
      local_optimisation: elements.localOptimisation.value as "noop" | "2opt"
    }
  } catch (error) {
    showError(`Invalid form data: ${error}`)
    return null
  }
}

// =======================
// File Handling
// =======================

async function handleFileUpload(file: File) {
  try {
    console.log(`Uploading file: ${file.name}`)
    const response = await parseTSP(file)
    
    if (!response.graph) {
      throw new Error("No graph data received")
    }
    
    currentGraph = response.graph
    currentFilename = file.name
    currentResult = null
    
    await updateGraphStats(currentGraph)
    
    // Render initial graph without path
    renderGraph(currentGraph)
    
    // Clear other plots
    Plotly.newPlot("pheromone-graph", [], {})
    Plotly.newPlot("heatmap", [], {})
    Plotly.newPlot("coeffs", [], {})
    Plotly.newPlot("score", [], {})
    
    console.log(`File loaded successfully: ${currentGraph.nodes.length} nodes, ${currentGraph.edges.length} edges`)
  } catch (error) {
    showError(`Failed to load TSP file: ${error}`)
    currentGraph = null
  }
}

async function handleFileChoose(file: string) {
  try {
    console.log(`Getting file: ${file}`)
    const response = await getTSP(file)
    
    if (!response.graph) {
      throw new Error("No graph data received")
    }
    
    currentGraph = response.graph
    currentFilename = file
    currentResult = null

    await updateGraphStats(currentGraph)
    
    // Render initial graph without path
    renderGraph(currentGraph)
    
    // Clear other plots
    Plotly.newPlot("pheromone-graph", [], {})
    Plotly.newPlot("heatmap", [], {})
    Plotly.newPlot("coeffs", [], {})
    Plotly.newPlot("score", [], {})

    console.log(`File loaded successfully: ${currentGraph.nodes.length} nodes, ${currentGraph.edges.length} edges`)
  } catch (error) {
    showError(`Failed to load TSP file: ${error}`)
    currentGraph = null
  }
}

// =======================
// HACO Execution
// =======================

async function runHaco() {
  if (!currentGraph) {
    showError("Please load a TSP file first")
    return
  }
  
  const formInput = getFormInput()
  if (!formInput) return
  
  disableRunButton(true)
  
  try {
    console.log("Building HACO request...")
    const request = await buildHacoRequest(currentGraph, formInput)


    console.log("Running HACO optimization...")
    const result = await runHacoDetails(request)
    
    currentResult = result
    
    // Render all visualizations
    renderAll(currentGraph, currentResult)
    
    console.log(`HACO completed! Best score: ${result.best_score}`)
  } catch (error) {
    showError(`HACO execution failed: ${error}`)
  } finally {
    disableRunButton(false)
  }
}

// =======================
// UI State Management
// =======================

async function updateGraphStats(graph: Graph | null) {
  try {
    const container = document.getElementById("graph-stats")!
    if (graph === null) {
      container.innerHTML = ""
      return
    }

    const request: GraphStatsRequest = {
      graph: graph
    }
    
    const stats = await graphStats(request)
    
    // Обновляем рекомендации в форме
    if (stats.recommended_pheromone_multiplier) {
      const pheromoneInput = document.getElementById("pheromone-multiplier") as HTMLInputElement
      const evaporationInput = document.getElementById("evaporation-rate") as HTMLInputElement
      
      // Показываем рекомендацию как placeholder или предзаполняем
      if (pheromoneInput && !pheromoneInput.value) {
        pheromoneInput.placeholder = `${stats.recommended_pheromone_multiplier.toFixed(2)}`
      }
      
      if (evaporationInput && !evaporationInput.value) {
        evaporationInput.placeholder = `${stats.recommended_evaporation_rate.toFixed(3)}`
      }
    }
    
    // Отправляем статистику на бэкенд для рендеринга компонента
    const newHTML = await renderGraphStats(stats)
    
    container.innerHTML = newHTML
  } catch (error) {
    console.error("Failed to load graph statistics:", error)
  }
}

function updateMutationFields() {
  const mutationType = elements.mutationType.value
  const uniformFields = [elements.mutationMin, elements.mutationMax]
  const gaussFields = [elements.mutationMean, elements.mutationStd]
  
  if (mutationType === "uniform") {
    uniformFields.forEach(field => field.disabled = false)
    gaussFields.forEach(field => field.disabled = true)
  } else {
    uniformFields.forEach(field => field.disabled = true)
    gaussFields.forEach(field => field.disabled = false)
  }
}

function updateSelectionFields() {
  const selectionType = elements.selectionType.value
  const tournamentKField = elements.tournamentK
  
  tournamentKField.disabled = selectionType !== "tournament"
}

// =======================
// Event Handlers
// =======================

function setupEventListeners() {
  // File input
  elements.fileInput.addEventListener("change", async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
      await handleFileUpload(file)
    }
  })

  elements.graphSelect.addEventListener("change", async (e) => {
    const file = elements.graphSelect.selectedOptions[0]
    if (file) {
      await handleFileChoose(file.value)
    }
  })
  
  // Run button
  elements.runBtn.addEventListener("click", runHaco)
  
  // UI state changes
  elements.mutationType.addEventListener("change", updateMutationFields)
  elements.selectionType.addEventListener("change", updateSelectionFields)
  
  // Initial UI state
  updateMutationFields()
  updateSelectionFields()
}


// =======================
// Initialization
// =======================

function init() {
  console.log("Initializing HACO Visualizer...")
  setupEventListeners()
  
  // Clear initial plots
  const plotIds = ["graph", "pheromone-graph", "heatmap", "coeffs", "score"]
  plotIds.forEach(id => {
    const div = document.getElementById(id)
    if (div) {
      Plotly.newPlot(id, [], {})
    }
  })
}

// Start the application
init()