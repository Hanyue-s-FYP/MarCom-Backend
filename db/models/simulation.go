package models

import "database/sql"

type SimulationCycle struct {
	ID           int
	Profit       sql.NullFloat64
	AgentActions []CycleAgentAction
}

type CycleAgentAction struct {
	ID                int
	Prompt            sql.NullString
	ActionType        sql.NullString
	ActionDescription sql.NullString
	AgentID           int
}

type Simulation struct {
	Environment
	Business
	SimulationCycles  []SimulationCycle
	ID                int
	Name              sql.NullString
	MaxCycleCount     sql.NullInt64
	IsPriceOptEnabled sql.NullBool
	Status            sql.NullString
}
