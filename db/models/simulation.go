package models

import (
	"time"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
)

type SimulationStatus int

const (
	SimulationIdle = iota
	SimulationRunning
	SimulationCompleted
)

// cycle can start from 0 (which is the initialisation cycle to init the simulation)
type SimulationCycle struct {
	ID               int
	CycleNumber      int
	SimulationId     int
	SimulationEvents []SimulationEvent
}

// (BUY/SKIP/TALK): agent takes action, SIMULATION: high level simulation related events, like initializing agent, ACTION_RESP: response to BUY actions of an agent
type SimulationEventType int

const (
	SimulationEventBuy = iota
	SimulationEventSkip
	SimulationEventMessage
	SimulationEventSimulation
	SimulationEventActionResp
)

func SimulationEventTypeMapper(evType string) SimulationEventType {
	mapper := map[string]SimulationEventType{
		"BUY":         SimulationEventBuy,
		"SKIP":        SimulationEventSkip,
		"MESSAGE":     SimulationEventMessage,
		"SIMULATION":  SimulationEventSimulation,
		"ACTION_RESP": SimulationEventActionResp,
	}
	return mapper[evType]
}

type SimulationEvent struct {
	Agent            *Agent // nullable
	ID               int
	EventType        SimulationEventType
	EventDescription string
}

type Simulation struct {
	EnvironmentID     int
	BusinessID        int
	SimulationCycles  []SimulationCycle
	ID                int
	Name              string
	MaxCycleCount     int
	IsPriceOptEnabled bool
	Status            int // may change to string later, but enums in go are just ints
}

type simulationModel struct{}

var SimulationModel *simulationModel

// only get the simulation out, simulation cycle should be handled in the business logic part
func (*simulationModel) GetByID(id int) (*Simulation, error) {
	var simulation Simulation
	query := `SELECT * FROM Simulations WHERE id = ?`
	row := db.GetDB().QueryRow(query, id)
	err := row.Scan(&simulation.ID, &simulation.Name, &simulation.MaxCycleCount, &simulation.IsPriceOptEnabled, &simulation.Status, &simulation.EnvironmentID, &simulation.BusinessID)
	if err != nil {
		return nil, err
	}
	return &simulation, nil
}

func (*simulationModel) GetAll() ([]Simulation, error) {
	var simulations []Simulation
	query := `SELECT * FROM Simulations`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var simulation Simulation
		if err := rows.Scan(&simulation.ID, &simulation.Name, &simulation.MaxCycleCount, &simulation.IsPriceOptEnabled, &simulation.Status, &simulation.EnvironmentID, &simulation.BusinessID); err != nil {
			return nil, err
		}
		simulations = append(simulations, simulation)
	}

	return simulations, nil
}

func (*simulationModel) GetAllByBusinessID(businessID int) ([]Simulation, error) {
	var simulations []Simulation
	query := `SELECT * FROM Simulations WHERE business_id = ?`
	rows, err := db.GetDB().Query(query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var simulation Simulation
		if err := rows.Scan(&simulation.ID, &simulation.Name, &simulation.MaxCycleCount, &simulation.IsPriceOptEnabled, &simulation.Status, &simulation.EnvironmentID, &simulation.BusinessID); err != nil {
			return nil, err
		}
		simulations = append(simulations, simulation)
	}

	return simulations, nil
}

func (*simulationModel) Create(s Simulation) error {
	query := `
        INSERT INTO Simulations (name, max_cycle_count, is_price_opt_enabled, status, environment_id, business_id)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.GetDB().Exec(query, s.Name, s.MaxCycleCount, s.IsPriceOptEnabled, s.Status, s.EnvironmentID, s.BusinessID)
	return err
}

func (*simulationModel) Update(s Simulation) error {
	query := `
        UPDATE Simulations
        SET name = ?, max_cycle_count = ?, is_price_opt_enabled = ?, status = ?, environment_id = ?
        WHERE id = ?;
    `
	_, err := db.GetDB().Exec(query, s.Name, s.MaxCycleCount, s.IsPriceOptEnabled, s.Status, s.EnvironmentID, s.ID)
	return err
}

func (*simulationModel) NewSimulationCycle(simId int, cycle SimulationCycle) (int, error) {
	query := `
		INSERT INTO SimulationCycles (simulation_id, cycle_number, time)
		VALUES (?, ?, ?)
	`

	res, err := db.GetDB().Exec(query, simId, cycle.CycleNumber, time.Now())

	if err != nil {
		return 0, err
	}

	if id, err := res.LastInsertId(); err != nil {
		return 0, err
	} else {
		return int(id), nil
	}
}

func (*simulationModel) GetSimulationCyclesBySimID(simId int) ([]SimulationCycle, error) {
	query := `
		SELECT id, simulation_id, cycle_number
		FROM SimulationCycles
		WHERE simulation_id = ?
		ORDER BY time
	`

	rows, err := db.GetDB().Query(query, simId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []SimulationCycle
	for rows.Next() {
		var cycle SimulationCycle
		err := rows.Scan(&cycle.ID, &cycle.SimulationId, &cycle.CycleNumber)
		if err != nil {
			return nil, err
		}
		cycleEvents, err := getSimulationEventByCycleID(cycle.ID)
		if err != nil {
			return nil, err
		}
		cycle.SimulationEvents = cycleEvents
		cycles = append(cycles, cycle)
	}

	return cycles, nil
}

func (*simulationModel) GetSimulationCycleIdBySimCycle(simId, cycleNum int) (int, error) {
	query := `
		SELECT id
		FROM SimulationCycles
		WHERE simulation_id = ? AND cycle_number = ?
	`
	var id int
	row := db.GetDB().QueryRow(query, simId, cycleNum)

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (*simulationModel) NewSimulationEvent(cycleId int, event SimulationEvent) (int, error) {
	// business logic there should handle that the cycle exists
	query := `
		INSERT INTO SimulationEvents (event_type, event_description, agent_id, cycle_id, time)
		VALUES (?, ?, ?, ?, ?)
	`

	// returns nil or int depends if a is nil or an Agent
	agentId := func(a *Agent) any {
		if a != nil {
			return a.ID
		} else {
			return 0 // since auto increment, 0 is not used, so can be used to indicate this field has no value
		}
	}

	res, err := db.GetDB().Exec(query, event.EventType, event.EventDescription, agentId(event.Agent), cycleId, time.Now())

	if err != nil {
		return 0, err
	}

	if id, err := res.LastInsertId(); err != nil {
		return 0, err
	} else {
		return int(id), nil
	}
}

// TODO fetch all the simulation cycle of that
func getSimulationEventByCycleID(cycleId int) ([]SimulationEvent, error) {
	query := `
        SELECT id, event_type, event_description, agent_id
        FROM SimulationEvents
        WHERE cycle_id = ?
        ORDER BY time
    `

	rows, err := db.GetDB().Query(query, cycleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []SimulationEvent
	for rows.Next() {
		var event SimulationEvent
		var agentId int
		err := rows.Scan(&event.ID, &event.EventType, &event.EventDescription, &agentId)
		if err != nil {
			return nil, err
		}

		if agentId != 0 {
			agent, err := AgentModel.GetByID(agentId)
			if err != nil {
				return nil, err
			}
			event.Agent = agent
		}

		events = append(events, event)
	}

	return events, nil
}
