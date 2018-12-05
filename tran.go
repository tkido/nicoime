package main

// Tran is raw transformation data
type Tran struct {
	Word     string
	Read     string
	Redirect bool
}

// Trans is map of transformation
type Trans map[string]Tran
