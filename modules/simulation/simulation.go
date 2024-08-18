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
	"strings"
	"sync"

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

	if !canUpdateSimulation(simulation) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to update simulation, simulation can not be running or completed or ran before",
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

func DeleteSimulation(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	// id of the environment shall be made accessible via route variable {id}
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

	if err = models.SimulationModel.Delete(idInt); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to delete simulation",
			LogMessage: fmt.Sprintf("failed to delete simulation: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully deleted simulation"}, nil
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
			LogMessage: fmt.Sprintf("failed to send simulation data to grpc server: %v", startSimulationErr),
		}
	}

	sim.Status = models.SimulationRunning
	models.SimulationModel.Update(*sim)
	sim, err = models.SimulationModel.GetByID(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to start simulation",
			LogMessage: fmt.Sprintf("failed to obtain simulation: %v", err),
		}
	}
	slog.Info(fmt.Sprintf("in StartSimulation: {SimulationStatus: %d}", sim.Status))

	go streamSimulationUpdate(sim.ID)

	return &modules.ExecResponse{Message: startSimulationMessage}, nil
}

type ComplexSimulationEvent struct {
	models.SimulationEvent
	CycleId int
}

type SimulationEventUpdateChannels struct {
	simulationUpdateListener chan ComplexSimulationEvent
	simulationUpdateEnd      chan struct{}
}

var (
	simulationUpdateListeners    map[int]SimulationEventUpdateChannels = make(map[int]SimulationEventUpdateChannels, 0)
	simulationUpdateListenerLock sync.Mutex
)

// not a API Func, rather, is already a http handler (SSE require unique handling)
func ListenToSimulationUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream") // declare that it is indeed SSE

	// assume simulation id is in the path
	id := r.PathValue("id")
	if id == "" {
		utils.ResponseError(w, utils.HttpError{
			Code:       http.StatusNotFound,
			Message:    "Expected ID in path, found empty string",
			LogMessage: "unexpected empty string in request when matching wildcard {id}",
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.ResponseError(w, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse simulation ID from request",
			LogMessage: fmt.Sprintf("failed to parse simulation ID from request: %v", err),
		})
		return
	}

	// check if the simulation is running (to minimise waste of resource)
	sim, err := models.SimulationModel.GetByID(idInt)
	if err != nil {
		utils.ResponseError(w, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain simulation from request",
			LogMessage: fmt.Sprintf("failed to obtain simulation from request: %v", err),
		})
		return
	}
	slog.Info(fmt.Sprintf("in ListenToSimulation: {SimulationStatus: %d}", sim.Status))

	if sim.Status != models.SimulationRunning {
		utils.ResponseError(w, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to listen to simulation event updates",
			LogMessage: "failed to listen to simulation event updates: simulation is not running",
		})
		return
	}

	ctx := r.Context()
	updateCh := make(chan ComplexSimulationEvent) // channel should be closed only when user disconnect or simulation complete or simulation paused
	updateEndCh := make(chan struct{})
	simulationUpdateListenerLock.Lock()
	simulationUpdateListeners[idInt] = SimulationEventUpdateChannels{
		simulationUpdateListener: updateCh,
		simulationUpdateEnd:      updateEndCh,
	}
	simulationUpdateListenerLock.Unlock()

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.ResponseError(w, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to listen to simulation event updates, SSE is not supported",
			LogMessage: fmt.Sprintf("failed to listen to simulation event updates: %v", "cannot cast http.Writer to http.Flusher"),
		})
		return
	}
OuterLoop:
	for {
		select {
		case <-ctx.Done(): // exit the loop when disconnected
			slog.Info(fmt.Sprintf("Client disconnected: %v", ctx.Err()))
			break OuterLoop
		case <-updateEndCh:
			slog.Info("sending stop signal to client")
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("event: %s\n", "simulation-stopped"))
			sb.WriteString(fmt.Sprintf("data: %v\n\n", "STOP"))
			fmt.Fprint(w, sb.String())
			flusher.Flush()
			break OuterLoop
		case e := <-updateCh:
			jsonBytes, err := json.Marshal(e)
			if err != nil {
				// don't need special message for internal server error, at least don't need to be returned to the client
				utils.ResponseError(w, utils.HttpError{
					Code:       http.StatusInternalServerError,
					LogMessage: fmt.Sprintf("failed to marshal JSON: %v", err),
				})
				break OuterLoop
			}
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("event: %s\n", "simulation-event"))
			sb.WriteString(fmt.Sprintf("data: %v\n\n", string(jsonBytes)))
			fmt.Fprint(w, sb.String())
			flusher.Flush()
		}
	}

	slog.Info(fmt.Sprintf("closing update channel for simulation id: %d", idInt))
	// remove the listener from the map (code reach here only in 2 conditions: user disconnect or simulation completed)
	simulationUpdateListenerLock.Lock()
	close(updateCh)
	close(updateEndCh)
	delete(simulationUpdateListeners, idInt)
	simulationUpdateListenerLock.Unlock()
}

func GetSimulationCyclesBySimID(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.SimulationCycle], error) {
	// id of the environment shall be made accessible via route variable {id}
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

	simCycles, err := models.SimulationModel.GetSimulationCyclesBySimID(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain simulation cycles",
			LogMessage: fmt.Sprintf("failed to obtain simulation cycles: %v", err),
		}
	}

	return &modules.SliceWrapper[models.SimulationCycle]{Data: simCycles}, nil
}

func GetSimulationCycleByCycleID(w http.ResponseWriter, r *http.Request) (*models.SimulationCycle, error) {
	// id of the environment shall be made accessible via route variable {id}
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
			Message:    "Failed to parse simulation cycle ID from request",
			LogMessage: fmt.Sprintf("failed to parse simulation cycle ID from request: %v", err),
		}
	}

	simCycle, err := models.SimulationModel.GetSimulationCycleByCycleID(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain simulation cycle",
			LogMessage: fmt.Sprintf("failed to obtain simulation cycle: %v", err),
		}
	}

	return simCycle, nil
}

func PauseSimulation(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
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

	var (
		pauseSimulationMessage string
		pauseSimulationErr     error
	)
	utils.UseCoreGRPCClient(func(client core_pb.MarcomServiceClient) {
		slog.Info("Sending pause request to grpc server")
		pauseSimulationResp, err := client.PauseSimulation(context.Background(), &core_pb.PauseRequest{
			SimulationId: int32(idInt),
		})
		if err != nil {
			pauseSimulationErr = err
			return
		}
		pauseSimulationMessage = pauseSimulationResp.Message
	})

	if pauseSimulationErr != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to pause simulation",
			LogMessage: fmt.Sprintf("failed to send pause simulation request to grpc server: %v", pauseSimulationErr),
		}
	}

	// in case

	return &modules.ExecResponse{Message: pauseSimulationMessage}, nil
}

// should not be called by client (stream is handled by start simulation)
func streamSimulationUpdate(id int) {
	isComplete := false
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
			if simulationEvent.Action == "COMPLETE" {
				isComplete = true
				break
			}
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
					id, err := models.SimulationModel.NewSimulationCycle(int(simulationEvent.SimulationId), models.SimulationCycle{
						CycleNumber:  int(simulationEvent.Cycle),
						SimulationId: int(simulationEvent.SimulationId),
					})
					if err != nil {
						slog.Error(fmt.Sprintf("failed to create cycle of simulation: %v", err))
						continue
					}
					err = newSimulationEventWithUpdate(int(simulationEvent.SimulationId), id, dbSimulationEvent)
					if err != nil {
						slog.Error(fmt.Sprintf("failed to create event of simulation cycle: %v", err))
						continue
					}
				} else {
					slog.Error(fmt.Sprintf("failed to obtain cycle of simulation: %v", err))
					continue
				}
			} else {
				err := newSimulationEventWithUpdate(int(simulationEvent.SimulationId), cycleId, dbSimulationEvent)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to create event of simulation cycle: %v", err))
					continue
				}
			}
		}
	})

	sim, err := models.SimulationModel.GetByID(id)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to obtain simulation: %v", err))
		return
	}
	if isComplete {
		slog.Info(fmt.Sprintf("Simulation with ID %d completed", id))

		sim.Status = models.SimulationCompleted
		err = models.SimulationModel.Update(*sim)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to update simulation status: %v", err))
			return
		}
	} else {
		slog.Info(fmt.Sprintf("Simulation with ID %d stopped streaming", id))

		sim.Status = models.SimulationIdle
		err = models.SimulationModel.Update(*sim)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to update simulation status: %v", err))
			return
		}
	}
	// I guess might need to send a signal to notify no more updates coming
	simulationUpdateListenerLock.Lock() // it might occur that the map is being deleted when this happened, so lock first
	// check if there is any listener of this event, if there is, check if it is open and send updates to it
	if ch, ok := simulationUpdateListeners[id]; ok {
		ch.simulationUpdateEnd <- struct{}{}
	}
	simulationUpdateListenerLock.Unlock()
}

// creates new simulation event that also handles sending update to listener (if any)
func newSimulationEventWithUpdate(simId, cycleId int, ev models.SimulationEvent) error {
	id, err := models.SimulationModel.NewSimulationEvent(cycleId, ev)
	if err != nil {
		return err
	}

	ev.ID = id
	simulationUpdateListenerLock.Lock() // it might occur that the map is being deleted when this happened, so lock first
	// check if there is any listener of this event, if there is, check if it is open and send updates to it
	if ch, ok := simulationUpdateListeners[simId]; ok {
		ch.simulationUpdateListener <- ComplexSimulationEvent{
			SimulationEvent: ev,
			CycleId:         cycleId,
		}
	}
	simulationUpdateListenerLock.Unlock()
	return nil
}

// only can update simulation if simulation is not completed or is not running or has never been started (can see if there are any cycles already)
func canUpdateSimulation(simulation models.Simulation) bool {
	if simulation.Status == models.SimulationCompleted || simulation.Status == models.SimulationRunning {
		return false
	}

	cycles, err := models.SimulationModel.GetSimulationCyclesBySimID(simulation.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to obtain simulation cycle by simulation id: %v", err))
	}

	if cycles == nil {
		return true
	} else {
		return false
	}
}
