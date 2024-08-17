package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

type SimplifiedEnvironment struct {
	Name           string
	SimulatedCount int
}

func CreateEnvironment(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var env models.Environment
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		return nil, &utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create environment",
			LogMessage: fmt.Sprintf("failed to parse environment JSON: %v", err),
		}
	}

	if err := models.EnvironmentModel.Create(env); err != nil {
		return nil, &utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create environment",
			LogMessage: fmt.Sprintf("failed to create environment: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully created environment"}, nil
}

func GetEnvironment(w http.ResponseWriter, r *http.Request) (*models.Environment, error) {
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
			Message:    "Failed to parse environment ID from request",
			LogMessage: fmt.Sprintf("failed to parse environment ID from request: %v", err),
		}
	}

	environment, err := models.EnvironmentModel.GetByID(idInt)
	if err != nil {
		var retErr utils.HttpError
		if errors.Is(err, models.ErrEnvironmentNotFound) {
			retErr = utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Environment not found in database",
				LogMessage: "environment not found",
			}
		} else {
			retErr = utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain environment",
				LogMessage: fmt.Sprintf("failed to get environment by id: %v", err),
			}
		}
		return nil, retErr
	}

	return environment, nil
}

func GetAllEnvironments(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Environment], error) {
	environments, err := models.EnvironmentModel.GetAll()
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environments",
			LogMessage: fmt.Sprintf("failed to obtain environments: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Environment]{Data: environments}, nil
}

func GetAllEnvironmentsByBusiness(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Environment], error) {
	// just in case still want investor module, see role, if role is business then can directly take user id if role is business then id should be in path
	role, err := strconv.Atoi(r.Header.Get("role"))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environment",
			LogMessage: fmt.Sprintf("failed to obtain user role when get environment by business: %v", err),
		}
	}

	var businessID int
	if role == models.INVESTOR {
		// id of the environment shall be made accessible via route variable {id} if the request is being made by investor not business
		id := r.PathValue("id")
		if id == "" {
			return nil, utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Expected ID in path, found empty string",
				LogMessage: "unexpected empty string in request when matching wildcard {id}",
			}
		}

		businessID, err = strconv.Atoi(id)
		if err != nil {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to parse business ID from request",
				LogMessage: fmt.Sprintf("failed to parse business ID from request: %v", err),
			}
		}
	} else {
		if businessID, err = strconv.Atoi(r.Header.Get("UserID")); err != nil {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain environments",
				LogMessage: fmt.Sprintf("failed to obtain user id when get environments by business id: %v", err),
			}
		}
	}
	// if still at 0 means it is not populated (suiran not very likely this will happen)
	if businessID == 0 {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environments",
			LogMessage: "failed to obtain user id when get environments by business id: unpopulated business id",
		}

	}

	environments, err := models.EnvironmentModel.GetAllByBusinessID(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environments",
			LogMessage: fmt.Sprintf("failed to obtain environments: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Environment]{Data: environments}, nil
}

func GetSimplifiedEnvironmentsWithProduct(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[SimplifiedEnvironment], error) {
	// should have product id in path
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
			Message:    "Failed to parse product ID from request",
			LogMessage: fmt.Sprintf("failed to parse product ID from request: %v", err),
		}
	}

	envs, err := models.EnvironmentModel.GetEnvironmentWithProduct(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environments with product",
			LogMessage: fmt.Sprintf("failed to obtain environments with product id: %v", err),
		}
	}

	var simplifiedEnv []SimplifiedEnvironment
	for _, v := range envs {
		// TODO populate simulated count
		simplifiedEnv = append(simplifiedEnv, SimplifiedEnvironment{Name: v.Name, SimulatedCount: 0})
	}

	return &modules.SliceWrapper[SimplifiedEnvironment]{Data: simplifiedEnv}, nil
}

func GetSimplifiedEnvironmentsWithAgent(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[SimplifiedEnvironment], error) {
	// should have agent id in path
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
			Message:    "Failed to parse agent ID from request",
			LogMessage: fmt.Sprintf("failed to parse agent ID from request: %v", err),
		}
	}

	envs, err := models.EnvironmentModel.GetEnvironmentWithAgent(idInt)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environments with agent",
			LogMessage: fmt.Sprintf("failed to obtain environments with agent id: %v", err),
		}
	}

	var simplifiedEnv []SimplifiedEnvironment
	for _, v := range envs {
		// TODO populate simulated count
		simplifiedEnv = append(simplifiedEnv, SimplifiedEnvironment{Name: v.Name, SimulatedCount: 0})
	}

	return &modules.SliceWrapper[SimplifiedEnvironment]{Data: simplifiedEnv}, nil
}

func UpdateEnvironment(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var env models.Environment
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		return nil, &utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update environment",
			LogMessage: fmt.Sprintf("failed to parse environment JSON: %v", err),
		}
	}

	// check if there are any simulation using this environment, if any, cannot delete
	if !canChangeEnv(env.ID) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete environment, environment is being referenced by some simulation",
		}
	}

	if err := models.EnvironmentModel.Update(env); err != nil {
		return nil, &utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create environment",
			LogMessage: fmt.Sprintf("failed to update environment: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully updated environment"}, nil
}

func DeleteEnvironment(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
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
			Message:    "Failed to parse environment ID from request",
			LogMessage: fmt.Sprintf("failed to parse environment ID from request: %v", err),
		}
	}

	// check if there are any simulation using this environment, if any, cannot delete
	if !canChangeEnv(idInt) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete environment, environment is being referenced by some simulation",
		}
	}

	if err = models.EnvironmentModel.Delete(idInt); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to delete environment",
			LogMessage: fmt.Sprintf("failed to delete environment: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully deleted environment"}, nil
}

// if any simulation is referencing the environment, the environment cannot be changed
func canChangeEnv(envId int) bool {
	simulations, err := models.SimulationModel.GetAllByEnvID(envId)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to obtain simulations by environment ID: %v", err))
		return false
	}
	if simulations == nil {
		return true
	} else {
		return false
	}
}
