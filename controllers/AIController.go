package controllers

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"restAPI/models"
	"sort"
)

type GeneticHandler struct{}

type GeneticHandlerResponse struct {
	Top_10_fitness  []float64     `json:"top_10_fitness"`
	Average_fitness []float64     `json:"average_fitness"`
	Top_fitness     []float64     `json:"top_fitness"`
	Items           []models.Item `json:"items"`
}

func NewGeneticHandler() *GeneticHandler {
	return &GeneticHandler{}
}

func (g *GeneticHandler) RunGenetic(w http.ResponseWriter, r *http.Request) {

	genetic := models.GeneticModel{}
	resultSet := GeneticHandlerResponse{}

	// Get the parameters from the form request
	err := json.NewDecoder(r.Body).Decode(&genetic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	top10, average, top, items := Evolve(genetic.Generations, genetic.N_Items, genetic.PopulationSize, genetic.Mutation, genetic.Limit)

	resultSet.Top_10_fitness = top10
	resultSet.Average_fitness = average
	resultSet.Top_fitness = top
	resultSet.Items = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultSet)
}

// Evolve - execute the genetic algorithm for a given number Generations, N_Items, Population, and Mutation,
// Keeping track of the top 10% of the population, the average fitness, and the top fitness
func Evolve(generations int, n_items int, populationSize int, mutation float64, limit int) ([]float64, []float64, []float64, []models.Item) {
	items := createItems(n_items, limit)
	population := createStartingPopulation(populationSize, n_items)
	top10 := []float64{}
	average := []float64{}
	top := []float64{}
	for i := 0; i < generations; i++ {
		scores := scorePopulation(population, items, limit)
		sortScoresAndPopulation(scores, population)

		scoreValues := []float64{}
		for j := 0; j < len(scores); j++ {
			scoreValues = append(scoreValues, float64(scores[j].Value))
		}

		// top10 is the average of the top 10 of the population
		top10 = append(top10, Average(scoreValues[:len(scoreValues)/10]))
		// average is the average of the entire population
		average = append(average, Average(scoreValues))
		// top is the top fitness of the entire population
		top = append(top, scoreValues[0])

		reproduce(population)
		mutate(population, mutation)
	}
	return top10, average, top, items
}

// Create Items array from given N_Items and Limit
// Item Values and Weights are randomly generated up to 1/4th of the Limit value
func createItems(n_items int, limit int) []models.Item {
	items := []models.Item{}
	for i := 0; i < n_items; i++ {
		items = append(items, models.Item{
			Value:  rand.Intn(10),
			Weight: rand.Intn(10),
		})
	}
	return items
}

// Create an array which contains true or false, randomly generated.
func makeChromosome(chromosomeLength int) models.Chromosome {
	responseArray := models.Chromosome{}
	for i := 0; i < chromosomeLength; i++ {
		responseArray = append(responseArray, rand.Intn(2) == 1)
	}
	return responseArray
}

// CreateStartingPopulation creates an array of chromosomes
func createStartingPopulation(populationSize int, chromosomeLength int) models.Population {
	populationArray := models.Population{}
	for i := 0; i < populationSize; i++ {
		populationArray = append(populationArray, makeChromosome(chromosomeLength))
	}
	return populationArray
}

// getPhenotype returns the phenotype of a chromosome, which is (1) the sum of the values and (2) the sum of the weights of the items in the chromosome
func getPhenotype(chromosome models.Chromosome, items []models.Item) (int, int) {
	var value int
	var weight int
	for i := 0; i < len(chromosome); i++ {
		if chromosome[i] {
			value += items[i].Value
			weight += items[i].Weight
		}
	}
	return value, weight
}

// score calculates the fitness of a chromosome, which is the sum of the values of the items in the chromosome
// If the weight of the chromosome is greater than the limit, the fitness is 0
func score(chromosome models.Chromosome, items []models.Item, limit int) models.Score {
	value, weight := getPhenotype(chromosome, items)
	if weight > limit {
		value -= 100
	}
	return models.Score{
		Value: value,
		Chrom: chromosome,
	}
}

// scorePopulation calculates the fitness of each chromosome in the population
func scorePopulation(population models.Population, items []models.Item, limit int) []models.Score {
	scores := []models.Score{}
	for i := 0; i < len(population); i++ {
		scores = append(scores, score(population[i], items, limit))
	}
	return scores
}

// sortPopulation sorts the population by fitness
func sortScoresAndPopulation(scores []models.Score, population models.Population) {
	// Sort the scores by value, create a new sortedPopulation array to return
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[j].Value < scores[i].Value
	})

	// reset and reorder population
	population = population[:0]
	for i := 0; i < len(scores); i++ {
		population = append(population, scores[i].Chrom)
	}
}

// Cut population in half, then duplicate the top half
func reproduce(population models.Population) {
	// // middle element
	// middle := (len(population) / 2) - 1
	// for i := 0; i < len(population)/2; i++ {
	// 	population[middle+i] = population[i]
	// }
	// element 62% of the way through the array
	middle := int(math.Ceil(float64(len(population)) * 0.62))
	j := 0
	for i := middle; i < len(population); i++ {
		// set population[i] to a new chromosome with the same values as population[j]
		population[i] = models.Chromosome{}
		for k := 0; k < len(population[j]); k++ {
			population[i] = append(population[i], population[j][k])
		}
		j++
	}
}

// Mutate a population by randomly flipping true and false in place, with the given Mutation proabability
func mutate(population models.Population, mutation float64) {
	for i := 0; i < len(population); i++ {
		for j := 0; j < len(population[i]); j++ {
			if rand.Float64() < mutation {
				population[i][j] = !(population[i][j])
			}
		}
	}
}

func Average(xs []float64) float64 {
	total := 0.0
	for _, x := range xs {
		total += x
	}
	return total / float64(len(xs))
}
