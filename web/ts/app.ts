import {
  parseTSP,
  getTSP,
  runHacoDetails,
} from "./api.js"

import {
  renderGraph,
  renderAll
} from "./render.js"

import { HacoRunRequest } from "./models/haco/request.js"
import { HacoRunDetailsResponse } from "./models/haco/response.js"
import { Graph } from "./models/graph.js"
import { HacoFormInput, buildHacoRequest } from "./form.js"

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
  runBtn: document.getElementById("run-btn") as HTMLButtonElement,
  
  // Basic params
  defaultAlpha: document.getElementById("default_alpha") as HTMLInputElement,
  defaultBeta: document.getElementById("default_beta") as HTMLInputElement,
  pheromoneMultiplier: document.getElementById("pheromone_multiplier") as HTMLInputElement,
  evaporationRate: document.getElementById("evaporation_rate") as HTMLInputElement,
  initialPheromone: document.getElementById("initial_pheromone") as HTMLInputElement,
  generationCount: document.getElementById("generation_count") as HTMLInputElement,
  colonySize: document.getElementById("colony_size") as HTMLInputElement,
  generationPeriod: document.getElementById("generation_period") as HTMLInputElement,
  parentCount: document.getElementById("parent_count") as HTMLInputElement,
  
  // Selection
  selectionType: document.getElementById("selection_type") as HTMLSelectElement,
  tournamentK: document.getElementById("tournament_k") as HTMLInputElement,
  
  // Mutation
  mutationType: document.getElementById("mutation_type") as HTMLSelectElement,
  mutationMin: document.getElementById("mutation_min") as HTMLInputElement,
  mutationMax: document.getElementById("mutation_max") as HTMLInputElement,
  mutationMean: document.getElementById("mutation_mean") as HTMLInputElement,
  mutationStd: document.getElementById("mutation_std") as HTMLInputElement,
  
  // Local optimisation
  localOptimisation: document.getElementById("local_optimisation") as HTMLSelectElement,
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
      
      crossover_type: "arithmetic",
      
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
// Keyboard Shortcuts
// =======================

function setupKeyboardShortcuts() {
  document.addEventListener("keydown", (e) => {
    // Ctrl/Cmd + Enter to run
    if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
      e.preventDefault()
      if (elements.runBtn.disabled === false) {
        runHaco()
      }
    }
    
    // Ctrl/Cmd + O to open file
    if ((e.ctrlKey || e.metaKey) && e.key === "o") {
      e.preventDefault()
      elements.fileInput.click()
    }
  })
}

// =======================
// Initialization
// =======================

async function loadDefaultGraph() {
  // Optional: Load a default TSP file if available
  // Uncomment and adjust if you have a default file on server
  /*
  try {
    const response = await getTSP("default.tsp")
    if (response.graph) {
      currentGraph = response.graph
      renderGraph(currentGraph)
      console.log("Default graph loaded")
    }
  } catch (error) {
    console.log("No default graph available")
  }
  */
}

function init() {
  console.log("Initializing HACO Visualizer...")
  setupEventListeners()
  setupKeyboardShortcuts()
  loadDefaultGraph()
  
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