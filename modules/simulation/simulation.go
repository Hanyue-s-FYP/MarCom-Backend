package simulation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
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
