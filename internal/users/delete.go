package users

import (
	"database/sql"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete User Name
	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Delete(db *sql.DB, id int64) error {
	stmt := "UPDATE users SET modified_at=$1, deleted=$2 WHERE id=$3"
	_, err := db.Exec(stmt, time.Now(), true, id)

	return err
}
