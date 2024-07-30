package models

// cycle can start from 0 (which is the initialisation cycle to init the simulation)
type SimulationCycle struct {
	ID               int
	Profit           float64
	SimulationEvents []SimulationEvent
}

type SimulationEvent struct {
	Agent            *Agent // nullable
	ID               int
	Prompt           string
	EventType        int
	EventDescription string
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
