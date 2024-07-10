package models

type SimulationCycle struct {
	ID           int
	Profit       float64
	AgentActions []CycleAgentAction
}

type CycleAgentAction struct {
	Agent
	ID                int
	Prompt            string
	ActionType        int
	ActionDescription string
}

type Simulation struct {
	Environment
	Business
	SimulationCycles  []SimulationCycle
	ID                int
	Name              string
	MaxCycleCount     int
	IsPriceOptEnabled bool
	Status            int // may change to string later, but enums in go are just ints
}
