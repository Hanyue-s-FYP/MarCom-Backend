package models

import "github.com/Hanyue-s-FYP/Marcom-Backend/db"

type SimulationStatus int

const (
	SimulationIdle = iota
	SimulationRunning
	SimulationCompleted
)

// cycle can start from 0 (which is the initialisation cycle to init the simulation)
type SimulationCycle struct {
	ID               int
	SimulationEvents []SimulationEvent
}

// (BUY/SKIP/TALK): agent takes action, SIMULATION: high level simulation related events, like initializing agent, ACTION_RESP: response to BUY actions of an agent
type SimulationEventType int

const (
	SimulationEventBuy = iota
	SimulationEventSkip
	SimulationEventTalk
	SimulationEventSimulation
	SimulationEventActionResp
)

type SimulationEvent struct {
	Agent            *Agent // nullable
	ID               int
	Prompt           string
	EventType        int
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

// TODO fetch all the simulation cycle of that
func getSimulationCyclesAndEvents(simID int) ([]SimulationCycle, error) {
	return nil, nil
}
