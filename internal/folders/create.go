package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
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

	id, err := Insert(h.db, f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	f.ID = id

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(f)
}

func Insert(db *sql.DB, f *Folder) (int64, error) {
	stmt := `INSERT INTO folders (parent_id, name, created_at, modified_at) VALUES ($1, $2, $3, $4) RETURNING id`
	result, err := db.Exec(stmt, f.ParentID, f.Name, f.CreatedAt, f.ModifiedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
