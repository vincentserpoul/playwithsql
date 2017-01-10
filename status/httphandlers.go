package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// EntityoneCreateHandler creates an entityone and returns it
func EntityoneCreateHandler(db *sqlx.DB, link *SQLLinkContainer) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var e Entityone
		err := e.Create(db, link)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		eJSON, err := json.Marshal(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(eJSON))

		return
	})
}

// EntityoneSelectHandler select an entityone and returns it
func EntityoneSelectHandler(db *sqlx.DB, link *SQLLinkContainer) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		es, err := SelectEntityoneByStatus(db, link, StatusCreated)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(es) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		esJSON, err := json.Marshal(es)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(esJSON))

		return
	})
}

// EntityoneSelectByPKHandler returns a selected entity
func EntityoneSelectByPKHandler(db *sqlx.DB, link *SQLLinkContainer) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		entityoneID, err := strconv.ParseInt(ps.ByName("entityoneID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		e, err := SelectEntityoneOneByPK(db, link, entityoneID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if e.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		eJSON, err := json.Marshal(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(eJSON))

		return
	})
}

// EntityoneDeleteByPKHandler updates an entityone to a deleted status
func EntityoneDeleteByPKHandler(db *sqlx.DB, link *SQLLinkContainer) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		entityoneID, err := strconv.ParseInt(ps.ByName("entityoneID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		e, err := SelectEntityoneOneByPK(db, link, entityoneID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if e.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = e.UpdateStatus(db, link, ActionCancel, StatusCancelled)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		return
	})
}
