package models

type GeneticModel struct {
	Generations    int     `json:"generations"`
	PopulationSize int     `json:"population"`
	Mutation       float64 `json:"mutation"`
	N_Items        int     `json:"n_items"`
	Limit          int     `json:"limit"`
}

type Item struct {
	Value  int `json:"value"`
	Weight int `json:"weight"`
}

type Chromosome []bool

type Population []Chromosome

type Score struct {
	Value int        `json:"value"`
	Chrom Chromosome `json:"chromosome"`
}

// No repository is needed for this model.
