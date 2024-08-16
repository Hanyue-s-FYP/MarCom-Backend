package simulation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
	core_pb "github.com/Hanyue-s-FYP/Marcom-Backend/proto"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func CreateSimulation(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var simulation models.Simulation

	if err := json.NewDecoder(r.Body).Decode(&simulation); err != nil {

		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create simulation",
			LogMessage: fmt.Sprintf("failed to decode simulation: %v", err),
		}
	}

	// append the id of the business into the simulation
	if businessID, err := strconv.Atoi(r.Header.Get("UserID")); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create simulation",
			LogMessage: fmt.Sprintf("failed to obtain user id when create simulation: %v", err),
		}
	} else {
		simulation.BusinessID = businessID
	}

	if err := models.SimulationModel.Create(simulation); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create simulation",
			LogMessage: fmt.Sprintf("failed to write simulation to db: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully created simulation"}, nil
}

func GetSimulation(w http.ResponseWriter, r *http.Request) (*models.Simulation, error) {
	// id of the simulation accessible via route variable {id}
	id := r.PathValue("id")
	if id == "" {
		return nil, utils.HttpError{
			Code:       http.StatusNotFound,
			Message:    "Expected ID in path, found empty string",
			LogMessage: "unexpected empty string in request when matching wildcard {id}",
		}
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse simulation ID from request",
			LogMessage: fmt.Sprintf("failed to parse simulation ID from request: %v", err),
		}
	}

	simulation, err := models.SimulationModel.GetByID(idInt)
	if err != nil {
		var retErr utils.HttpError
		if errors.Is(err, models.ErrProductNotFound) {
			retErr = utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Simulation not found in database",
				LogMessage: "simulation not found",
			}
		} else {
			retErr = utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain simulation",
				LogMessage: fmt.Sprintf("failed to get simulation by id: %v", err),
			}
		}
		return nil, retErr
	}

	return simulation, nil
}

func GetAllSimulations(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Simulation], error) {
	simulations, err := models.SimulationModel.GetAll()
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain simulations",
			LogMessage: fmt.Sprintf("failed to obtain simulations by business id: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Simulation]{Data: simulations}, nil
}

func GetSimulationsByBusinessID(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Simulation], error) {
	// id of the business accessible via route variable {id}
	id := r.PathValue("id")
	if id == "" {
		return nil, utils.HttpError{
			Code:       http.StatusNotFound,
			Message:    "Expected ID in path, found empty string",
			LogMessage: "unexpected empty string in request when matching wildcard {id}",
		}
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse business ID from request",
			LogMessage: fmt.Sprintf("failed to parse business ID from request: %v", err),
		}
	}

	simulations, err := models.SimulationModel.GetAllByBusinessID(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain simulations",
			LogMessage: fmt.Sprintf("failed to obtain simulations by business id: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Simulation]{Data: simulations}, nil
}

func UpdateSimulation(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var simulation models.Simulation

	if err := json.NewDecoder(r.Body).Decode(&simulation); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update simulation",
			LogMessage: fmt.Sprintf("failed to decode simulation: %v", err),
		}
	}

	if err := models.SimulationModel.Update(simulation); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update simulation",
			LogMessage: fmt.Sprintf("failed to write simulation to db: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully update simulation"}, nil
}

type SimulationStartRequest struct {
	ID int
}

func StartSimulation(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	id := r.PathValue("id")
	if id == "" {
		return nil, utils.HttpError{
			Code:       http.StatusNotFound,
			Message:    "Expected ID in path, found empty string",
			LogMessage: "unexpected empty string in request when matching wildcard {id}",
		}
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse simulation ID from request",
			LogMessage: fmt.Sprintf("failed to parse simulation ID from request: %v", err),
		}
	}

	sim, err := models.SimulationModel.GetByID(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to start simulation",
			LogMessage: fmt.Sprintf("failed to obtain simulation: %v", err),
		}
	}

	// obtain environment, products in the environment, agents in the environment
	environment, err := models.EnvironmentModel.GetByID(sim.EnvironmentID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to start simulation",
			LogMessage: fmt.Sprintf("failed to obtain environment: %v", err),
		}
	}

	var transformed_agents []*core_pb.Agent
	for _, a := range environment.Agents {
		var transformed_attrs []*core_pb.AgentAttribute
		for _, attr := range a.Attributes {
			transformed_attrs = append(transformed_attrs, &core_pb.AgentAttribute{Key: attr.Key, Value: attr.Value})
		}
		transformed_agents = append(transformed_agents, &core_pb.Agent{
			Id:    int32(a.ID),
			Name:  a.Name,
			Desc:  a.GeneralDescription,
			Attrs: transformed_attrs,
		})
	}

	var transformed_products []*core_pb.Product
	for _, p := range environment.Products {
		transformed_products = append(transformed_products, &core_pb.Product{
			Id:    int32(p.ID),
			Name:  p.Name,
			Desc:  p.Description,
			Price: float32(p.Price),
			Cost:  float32(p.Cost),
		})
	}

	var (
		startSimulationMessage string
		startSimulationErr     error
	)
	utils.UseCoreGRPCClient(func(client core_pb.MarcomServiceClient) {
		slog.Info("Sending simulation data to simulation server")
		startSimulationResp, err := client.StartSimulation(context.Background(), &core_pb.SimulationRequest{
			Id:          int32(sim.ID),
			EnvDesc:     environment.Description,
			Agents:      transformed_agents,
			Products:    transformed_products,
			TotalCycles: int32(sim.MaxCycleCount),
		})
		if err != nil {
			startSimulationErr = err
			return
		}
		slog.Info(fmt.Sprintf("Start Simulation Response: %s", startSimulationResp.Message))
		startSimulationMessage = startSimulationResp.Message
	})

	if startSimulationErr != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to start simulation",
			LogMessage: fmt.Sprintf("failed to send simulation data to grpc server: %v", err),
		}
	}

	go streamSimulationUpdate(sim.ID)

	return &modules.ExecResponse{Message: startSimulationMessage}, nil
}

// should not be called by client (stream is handled by start simulation)
func streamSimulationUpdate(id int) {
	utils.UseCoreGRPCClient(func(client core_pb.MarcomServiceClient) {
		slog.Info("Requesting StreamUpdates from grpc server")
		stream, err := client.StreamSimulationUpdates(context.Background(), &core_pb.StreamRequest{SimulationId: int32(id)})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to call StreamSimulationUpdates on simulation server: %v", err))
		}
		for {
			simulationEvent, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error(fmt.Sprintf("%v.StreamSimulationUpdates(_) = _, %v", client, err))
				break
			}
			slog.Info(fmt.Sprintf("Got event from simulation server: %v", simulationEvent))
			dbSimulationEvent := models.SimulationEvent{
				EventType:        models.SimulationEventTypeMapper(simulationEvent.Action),
				EventDescription: simulationEvent.Content,
			}
			if simulationEvent.AgentId != 0 {
				agent, err := models.AgentModel.GetByID(int(simulationEvent.AgentId))
				if err != nil {
					slog.Error(fmt.Sprintf("failed to obtain agent of the event: %v", err))
					continue
				}
				dbSimulationEvent.Agent = agent
			}
			if cycleId, err := models.SimulationModel.GetSimulationCycleIdBySimCycle(int(simulationEvent.SimulationId), int(simulationEvent.Cycle)); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err := models.SimulationModel.NewSimulationCycle(int(simulationEvent.SimulationId), models.SimulationCycle{
						CycleNumber:  int(simulationEvent.Cycle),
						SimulationId: int(simulationEvent.SimulationId),
					})
					if err != nil {
						slog.Error(fmt.Sprintf("failed to create cycle of simulation: %v", err))
						continue
					}
					err = models.SimulationModel.NewSimulationEvent(cycleId, dbSimulationEvent)
					if err != nil {
						slog.Error(fmt.Sprintf("failed to create event of simulation cycle: %v", err))
						continue
					}
				} else {
					slog.Error(fmt.Sprintf("failed to obtain cycle of simulation: %v", err))
					continue
				}
			} else {
				err := models.SimulationModel.NewSimulationEvent(cycleId, dbSimulationEvent)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to create event of simulation cycle: %v", err))
					continue
				}
			}
		}
	})
}
