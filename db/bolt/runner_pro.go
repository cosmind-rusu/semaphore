//go:build pro

package bolt

import (
	"github.com/semaphoreui/semaphore/db"
	"go.etcd.io/bbolt"
)

func (d *BoltDb) GetRunner(projectID int, runnerID int) (runner db.Runner, err error) {
	err = d.getObject(projectID, db.RunnerProps, intObjectID(runnerID), &runner)
	return
}

func (d *BoltDb) GetRunners(projectID int, activeOnly bool) (runners []db.Runner, err error) {
	runners = make([]db.Runner, 0)
	err = d.getObjects(projectID, db.RunnerProps, db.RetrieveQueryParams{}, func(i interface{}) bool {
		runner := i.(db.Runner)
		if activeOnly {
			return runner.Active
		}
		return true
	}, &runners)
	return
}

func (d *BoltDb) DeleteRunner(projectID int, runnerID int) (err error) {
	return d.db.Update(func(tx *bbolt.Tx) error {
		return d.deleteObject(projectID, db.RunnerProps, intObjectID(runnerID), tx)
	})
}
