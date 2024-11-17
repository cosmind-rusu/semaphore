//go:build pro

package projects

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	log "github.com/sirupsen/logrus"
)

type runnerWithToken struct {
	db.Runner
	Token string `json:"token"`
}

func GetRunners(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	runners, err := helpers.Store(r).GetRunners(project.ID, false)

	if err != nil {
		panic(err)
	}

	helpers.WriteJSON(w, http.StatusOK, runners)
}

func AddRunner(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	var runner db.Runner
	if !helpers.Bind(w, r, &runner) {
		return
	}

	if runner.ProjectID == nil {
		log.Error(fmt.Sprintf("Project ID not provided for new project runner"))

		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID required",
		})
		return
	}

	if *runner.ProjectID != project.ID {
		log.Error(fmt.Sprintf("Project ID in body and URL must be the same: %v vs. %v", runner.ProjectID, project.ID))

		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	newRunner, err := helpers.Store(r).CreateRunner(runner)

	if err != nil {
		log.Warn("Runner is not created: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, runnerWithToken{
		Runner: newRunner,
		Token:  newRunner.Token,
	})
}

func RunnerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runnerID, err := helpers.GetIntParam("runner_id", w, r)
		project := context.Get(r, "project").(db.Project)

		if err != nil {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "runner_id required",
			})
			return
		}

		store := helpers.Store(r)

		runner, err := store.GetRunner(project.ID, runnerID)

		if err != nil {
			helpers.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": "Runner not found",
			})
			return
		}

		context.Set(r, "runner", &runner)
		next.ServeHTTP(w, r)
	})
}

func GetRunner(w http.ResponseWriter, r *http.Request) {
	runner := context.Get(r, "runner").(*db.Runner)

	helpers.WriteJSON(w, http.StatusOK, runner)
}

func UpdateRunner(w http.ResponseWriter, r *http.Request) {
	oldRunner := context.Get(r, "runner").(*db.Runner)

	var runner db.Runner
	if !helpers.Bind(w, r, &runner) {
		return
	}

	store := helpers.Store(r)

	runner.ID = oldRunner.ID
	runner.ProjectID = nil

	err := store.UpdateRunner(runner)

	if err != nil {
		helpers.WriteErrorStatus(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteRunner(w http.ResponseWriter, r *http.Request) {
	runner := context.Get(r, "runner").(*db.Runner)

	store := helpers.Store(r)

	err := store.DeleteRunner(*runner.ProjectID, runner.ID)

	if err != nil {
		helpers.WriteErrorStatus(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SetRunnerActive(w http.ResponseWriter, r *http.Request) {
	runner := context.Get(r, "runner").(*db.Runner)

	store := helpers.Store(r)

	var body struct {
		Active bool `json:"active"`
	}

	if !helpers.Bind(w, r, &body) {
		helpers.WriteErrorStatus(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	runner.Active = body.Active

	err := store.UpdateRunner(*runner)

	if err != nil {
		helpers.WriteErrorStatus(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
