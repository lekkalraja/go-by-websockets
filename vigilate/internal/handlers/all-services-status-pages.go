package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/lekkalraja/go-by-websockets/vigilate/internal/helpers"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	services, err := repo.DB.GetServicesByStatus(r.Context(), "healthy")

	if err != nil {
		helpers.ServerError(w, r, err)
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)

	err = helpers.RenderPage(w, r, "healthy", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	services, err := repo.DB.GetServicesByStatus(r.Context(), "warning")

	if err != nil {
		helpers.ServerError(w, r, err)
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)

	err = helpers.RenderPage(w, r, "warning", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	services, err := repo.DB.GetServicesByStatus(r.Context(), "problem")

	if err != nil {
		helpers.ServerError(w, r, err)
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)

	err = helpers.RenderPage(w, r, "problems", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	services, err := repo.DB.GetServicesByStatus(r.Context(), "pending")

	if err != nil {
		helpers.ServerError(w, r, err)
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)

	err = helpers.RenderPage(w, r, "pending", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
