package entities

type Transaction struct {
	ID    string `json:"id"`
	Major int    `json:"major"`
	Minor int    `json:"minor"`
}
