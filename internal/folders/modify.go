package folders

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	f := new(Folder)

	err := json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = f.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Update(h.db, int64(id), f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: get folder
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(f)
}

func Update(db *sql.DB, id int64, f *Folder) error {
	stmt := `UPDATE folders SET parent_id = $1, name = $2, modified_at = $3 WHERE id = $4`
	_, err := db.Exec(stmt, f.ParentID, f.Name, f.ModifiedAt, id)
	if err != nil {
		return err
	}
	return nil
}
