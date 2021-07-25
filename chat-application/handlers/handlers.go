package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(), // remove in production
)

func Home(w http.ResponseWriter, r *http.Request) {
	err := render(w, "home", nil)
	if err != nil {
		log.Printf("Something went wrong while rendering : %v", err)
	}
}

func render(w http.ResponseWriter, path string, variables jet.VarMap) error {
	template, err := views.GetTemplate(path)

	if err != nil {
		return err
	}

	err = template.Execute(w, variables, nil)

	if err != nil {
		return err
	}

	return nil
}
