//go:build pro

package sql

import (
	"github.com/Masterminds/squirrel"
	"github.com/semaphoreui/semaphore/db"
)

func (d *SqlDb) GetRunner(projectID int, runnerID int) (runner db.Runner, err error) {
	err = d.getObject(projectID, db.RunnerProps, runnerID, &runner)
	return
}

func (d *SqlDb) GetRunners(projectID int, activeOnly bool) (runners []db.Runner, err error) {
	err = d.getObjects(projectID, db.RunnerProps, db.RetrieveQueryParams{}, func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		if activeOnly {
			builder = builder.Where("active=?", activeOnly)
		}

		return builder
	}, &runners)
	return
}

func (d *SqlDb) DeleteRunner(projectID int, runnerID int) (err error) {
	err = d.deleteObject(projectID, db.RunnerProps, runnerID)
	return
}
