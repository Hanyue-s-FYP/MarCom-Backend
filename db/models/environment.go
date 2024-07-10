package models

type Environment struct {
	Business
	Products    []Product
	Agents      []Agent
	ID          int
	Name        string
	Description string
}
