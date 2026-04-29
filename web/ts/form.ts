
import { Graph } from "./models/graph.js"
import { HacoRunRequest } from "./models/haco/request.js"
import { TournamentSelection, BestSelection, BaseSelection } from "./models/haco/selection.js"
import { ArithmeticCrossover, BaseCrossover } from "./models/haco/crossover.js"
import { UniformMutation, GaussMutation, BaseMutation } from "./models/haco/mutation.js"
import { BaseOptimisation, NoOpLocalOptimisation, TwoOptLocalOptimisation } from "./models/haco/optimisation.js"
import { validateDto } from "./utils/validate.js"

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

  crossover_type: "arithmetic"
  
  selection_type: "best" | "tournament"
  tournament_k: number | undefined
  
  mutation_type: "uniform" | "gauss"
  mutation_min: number | undefined
  mutation_max: number | undefined
  mutation_mean: number | undefined
  mutation_std: number | undefined
  
  local_optimisation: "noop" | "2opt"
}

export async function buildHacoRequest(graph: Graph, input: HacoFormInput): Promise<HacoRunRequest> {

    let selection: BaseSelection
    if (input.selection_type === "tournament") {
        selection = new TournamentSelection(input.tournament_k!)
    } else {
        selection = new BestSelection()
    }

    let mutation: BaseMutation
    if (input.mutation_type === "uniform") {
        mutation = new UniformMutation(input.mutation_min!, input.mutation_max!)
    } else {
        mutation = new GaussMutation(input.mutation_mean!, input.mutation_std!)
    }

    let crossover: BaseCrossover
    if (input.crossover_type === "arithmetic") {
        crossover = new ArithmeticCrossover()
    } else {
        crossover = new ArithmeticCrossover()
    }

    let local_optimisation: BaseOptimisation
    if (input.local_optimisation === "noop") {
        local_optimisation = new NoOpLocalOptimisation()
    } else {
        local_optimisation = new TwoOptLocalOptimisation()
    }

    let req: HacoRunRequest = {
        graph,
        
        default_alpha: input.default_alpha,
        default_beta: input.default_beta,
        pheromone_multiplier: input.pheromone_multiplier,
        evaporation_rate: input.evaporation_rate,
        initial_pheromone: input.initial_pheromone,
        
        generation_count: input.generation_count,
        colony_size: input.colony_size,
        generation_period: input.generation_period,
        parent_count: input.parent_count,
        
        selection: selection,
        crossover: crossover,
        mutation: mutation,
        local_optimisation: local_optimisation
    }

    return validateDto(HacoRunRequest, req)
}