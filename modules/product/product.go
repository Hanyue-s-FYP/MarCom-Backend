package product

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

func CreateProduct(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create product",
			LogMessage: fmt.Sprintf("failed to decode product: %v", err),
		}
	}

	// append the id of the business into the product
	if businessID, err := strconv.Atoi(r.Header.Get("UserID")); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create product",
			LogMessage: fmt.Sprintf("failed to obtain user id when create product: %v", err),
		}
	} else {
		product.BusinessID = businessID
	}
	if err := models.ProductModel.Create(product); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create product",
			LogMessage: fmt.Sprintf("failed to create product: %v", err),
		}
	}
	return &modules.ExecResponse{Message: "Successfully created product"}, nil
}

func GetProduct(w http.ResponseWriter, r *http.Request) (*models.Product, error) {
	// id of the product accessible via route variable {id}
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

	product, err := models.ProductModel.GetByID(idInt)
	if err != nil {
		var retErr utils.HttpError
		if errors.Is(err, models.ErrProductNotFound) {
			retErr = utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Product not found in database",
				LogMessage: "product not found",
			}
		} else {
			retErr = utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain product",
				LogMessage: fmt.Sprintf("failed to get product by id: %v", err),
			}
		}
		return nil, retErr
	}

	return product, nil
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Product], error) {
	products, err := models.ProductModel.GetAll()
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain products",
			LogMessage: fmt.Sprintf("failed to obtain products: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Product]{Data: products}, nil
}

func GetAllProductsByBusiness(w http.ResponseWriter, r *http.Request) (*modules.SliceWrapper[models.Product], error) {
	// id of the product accessible via route variable {id}
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

	products, err := models.ProductModel.GetAllByBusinessID(businessID)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain products",
			LogMessage: fmt.Sprintf("failed to obtain products: %v", err),
		}
	}

	return &modules.SliceWrapper[models.Product]{Data: products}, nil
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update product",
			LogMessage: fmt.Sprintf("failed to parse product JSON: %v", err),
		}
	}

	if !canChangeProduct(product.ID) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete product, product is being referenced in other environments",
		}
	}

	if err := models.ProductModel.Update(product); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update product",
			LogMessage: fmt.Sprintf("failed to update product: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully updated product"}, nil
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	// id of the product accessible via route variable {id}
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

	if !canChangeProduct(idInt) {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: "Failed to delete product, product is being referenced in other environments",
		}
	}

	if err = models.ProductModel.Delete(idInt); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to delete product",
			LogMessage: fmt.Sprintf("failed to delete product: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully deleted product"}, nil
}

// cannot update or delete if product is used by other environment
func canChangeProduct(id int) bool {
	env, err := models.EnvironmentModel.GetEnvironmentWithProduct(id)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to obtain environments with product: %v", err))
		return false
	}

	if env == nil {
		return true
	} else {
		return false
	}
}
