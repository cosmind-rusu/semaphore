package projects

import (
	"github.com/gorilla/context"
	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"net/http"
)

func GetRunners(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	runners, err := helpers.Store(r).GetRunners(project.ID)

	if err != nil {
		panic(err)
	}

	var result = make([]db.Runner, 0)

	for _, runner := range runners {
		result = append(result, runner)
	}

	helpers.WriteJSON(w, http.StatusOK, result)
}

func AddRunner(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func RunnerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func GetRunner(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func UpdateRunner(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func DeleteRunner(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func SetRunnerActive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
