package dashboard

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func GetDashboardProduct(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.DashboardProduct], error) {
	businessID, err := strconv.Atoi(r.Header.Get("UserID"))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain product",
			LogMessage: fmt.Sprintf("failed to obtain user id when obtain dashboard product: %v", err),
		}
	}

	prods, err := models.ProductModel.GetDashboardData(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain dashboard product data",
			LogMessage: fmt.Sprintf("failed to obtain dashboard product data: %v", err),
		}
	}
	return &modules.SliceWrapper[models.DashboardProduct]{Data: prods}, nil
}

func GetDashboardAgent(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.DashboardAgent], error) {
	businessID, err := strconv.Atoi(r.Header.Get("UserID"))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain agents",
			LogMessage: fmt.Sprintf("failed to obtain user id when obtain dashboard agent: %v", err),
		}
	}

	agents, err := models.AgentModel.GetDashboardData(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain agents",
			LogMessage: fmt.Sprintf("failed to obtain dashboard agent: %v", err),
		}
	}

	return &modules.SliceWrapper[models.DashboardAgent]{Data: agents}, nil
}

func GetDashboardEnvironment(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.DashboardEnvironment], error) {
	businessID, err := strconv.Atoi(r.Header.Get("UserID"))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environment",
			LogMessage: fmt.Sprintf("failed to obtain user id when obtain dashboard environment: %v", err),
		}
	}

	environments, err := models.EnvironmentModel.GetDashboardData(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain environment",
			LogMessage: fmt.Sprintf("failed to obtain dashboard environment: %v", err),
		}
	}

	return &modules.SliceWrapper[models.DashboardEnvironment]{Data: environments}, nil
}
