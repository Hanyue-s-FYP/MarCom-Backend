package agent

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

func CreateAgent(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create agent",
			LogMessage: fmt.Sprintf("failed to decode agent: %v", err),
		}
	}

	// append the id of the business into the agent
	if businessID, err := strconv.Atoi(r.Header.Get("UserID")); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create agent",
			LogMessage: fmt.Sprintf("failed to obtain user id when create agent: %v", err),
		}
	} else {
		agent.BusinessID = businessID
	}
	if err := models.AgentModel.Create(agent); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create agent",
			LogMessage: fmt.Sprintf("failed to create agent: %v", err),
		}
	}
	return &modules.ExecResponse{Message: "Successfully created agent"}, nil
}

func GetAgent(w http.ResponseWriter, r *http.Request) (*models.Agent, error) {
	// id of the agent accessible via route variable {id}
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

	agent, err := models.AgentModel.GetByID(idInt)
	if err != nil {
		var retErr utils.HttpError
		if errors.Is(err, models.ErrAgentNotFound) {
			retErr = utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Agent not found in database",
				LogMessage: "agent not found",
			}
		} else {
			retErr = utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain agent",
				LogMessage: fmt.Sprintf("failed to get agent by id: %v", err),
			}
		}
		return nil, retErr
	}

	return agent, nil
}

func GetAllAgents(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Agent], error) {
	agents, err := models.AgentModel.GetAll()
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain agents",
			LogMessage: fmt.Sprintf("failed to obtain agents: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Agent]{Data: agents}, nil
}

func GetAllAgentsByBusiness(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Agent], error) {
	// id of the business accessible via route variable {id}
	id := r.PathValue("id")
	if id == "" {
		return nil, utils.HttpError{
			Code:       http.StatusNotFound,
			Message:    "Expected ID in path, found empty string",
			LogMessage: "unexpected empty string in request when matching wildcard {id}",
		}
	}

	businessID, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse business ID from request",
			LogMessage: fmt.Sprintf("failed to parse business ID from request: %v", err),
		}
	}

	// if still at 0 means it is not populated (suiran not very likely this will happen)
	if businessID == 0 {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain agent",
			LogMessage: "failed to obtain user id when get agent by business id: unpopulated business id",
		}

	}

	agents, err := models.AgentModel.GetAllByBusinessID(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain agents",
			LogMessage: fmt.Sprintf("failed to obtain agents: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Agent]{Data: agents}, nil
}

func UpdateAgent(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update agent",
			LogMessage: fmt.Sprintf("failed to parse agent JSON: %v", err),
		}
	}

	if !canChangeAgent(agent.ID) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete agent, agent is being referenced by other environments",
		}
	}

	if err := models.AgentModel.Update(agent); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update agent",
			LogMessage: fmt.Sprintf("failed to update agent: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully updated agent"}, nil
}

func DeleteAgent(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	// id of the agent accessible via route variable {id}
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

	if !canChangeAgent(idInt) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete agent, agent is being referenced by other environments",
		}
	}

	if err = models.AgentModel.Delete(idInt); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to delete agent",
			LogMessage: fmt.Sprintf("failed to delete agent: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully deleted agent"}, nil
}

// cannot delete or update if agent exist in any of the environments
func canChangeAgent(id int) bool {
	env, err := models.EnvironmentModel.GetEnvironmentWithAgent(id)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to obtain environments with agent: %v", err))
		return false
	}

	if env == nil {
		return true
	} else {
		return false
	}
}
