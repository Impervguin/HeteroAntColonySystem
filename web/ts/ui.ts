function toggle(el: HTMLElement, show: boolean) {
  el.classList.toggle("hidden", !show)
}

// =======================
// GRAPH MODE
// =======================

const graphMode = document.getElementById("graph-mode") as HTMLSelectElement
const fileBlock = document.getElementById("graph-file-block")!
const existingBlock = document.getElementById("graph-existing-block")!

graphMode.addEventListener("change", () => {
  toggle(fileBlock, graphMode.value === "file")
  toggle(existingBlock, graphMode.value === "existing")
})

// =======================
// SELECTION
// =======================

const selection = document.getElementById("selection-type") as HTMLSelectElement
const tournament = document.getElementById("selection-tournament")!

selection.addEventListener("change", () => {
  toggle(tournament, selection.value === "tournament")
})

// =======================
// MUTATION
// =======================

const mutation = document.getElementById("mutation-type") as HTMLSelectElement
const uniform = document.getElementById("mutation-uniform")!
const gauss = document.getElementById("mutation-gauss")!

mutation.addEventListener("change", () => {
  toggle(uniform, mutation.value === "uniform")
  toggle(gauss, mutation.value === "gauss")
})

// INIT (важно)
selection.dispatchEvent(new Event("change"))
mutation.dispatchEvent(new Event("change"))
graphMode.dispatchEvent(new Event("change"))