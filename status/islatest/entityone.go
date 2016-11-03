package islatest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// Entityone represents an event
type Entityone struct {
	ID          int64     `db:"entityone_id"`
	TimeCreated time.Time `db:"time_created"`
	Status
}

// Status of the entity
type Status struct {
	entityID    int64     `db:"entityone_id"`
	actionID    ActionID  `db:"action_id"`
	statusID    StatusID  `db:"status_id"`
	TimeCreated time.Time `db:"time_created"`
}

// ActionID represents the action performed on the tradeoffer request
type ActionID int

const (
	// ActionCreate is triggered when the Entityone is created
	ActionCreate ActionID = 1
	// ActionCancel  is triggered when the Entityone is cancelled
	ActionCancel ActionID = 999
	// ActionClose is triggered when the Entityone is closed
	ActionClose ActionID = 500
)

func (s ActionID) String() string {
	return strconv.Itoa(int(s))
}

// StatusID represents the state of the tradeoffer, see constants
type StatusID int

const (
	// StatusCreated is when a Entityone is just created
	StatusCreated StatusID = 1
	// StatusCancelled when a Entityone is cancelled
	StatusCancelled StatusID = 999
	// StatusClosed is not changeable anymore, final status
	StatusClosed StatusID = 1000
)

func (s StatusID) String() string {
	return strconv.Itoa(int(s))
}

// SQLLink is used to define SQL interactions
type SQLLink interface {
	InitDB(exec sqlx.Execer, dbName string) (errExec error)
	DestroyDB(exec sqlx.Execer, dbName string) (errExec error)
	MigrateUp(exec sqlx.Execer) (errExec error)
	MigrateDown(exec sqlx.Execer) (errExec error)
	InsertOne(sqlx.Ext) (int64, error)
	SaveStatus(exec sqlx.Execer, entityID int64, actionID int, statusID int) error
}

// Create will create an entityone
func (e *Entityone) Create(db *sqlx.DB, link SQLLink) (err error) {
	tx := db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	e.ID, err = link.InsertOne(db)
	if err != nil {
		return fmt.Errorf("Entityone createEntityone: %v", err)
	}

	err = link.SaveStatus(tx, e.ID, int(ActionCreate), int(StatusCreated))
	if err != nil {
		return fmt.Errorf("Entityone createEntityone: %v", err)
	}

	e.Status = Status{
		entityID:    e.ID,
		actionID:    ActionCreate,
		statusID:    StatusCreated,
		TimeCreated: time.Now(),
	}

	return nil
}

// UpdateStatus will update the status of an Entityone into db
func (e *Entityone) UpdateStatus(exec sqlx.Ext, link SQLLink, actionID ActionID, statusID StatusID) error {
	err := link.SaveStatus(exec, e.ID, int(actionID), int(statusID))

	if err != nil {
		return fmt.Errorf("entityone UpdateStatus(): %v", err)
	}

	// Update status
	e.actionID = actionID
	e.statusID = statusID

	return nil
}
