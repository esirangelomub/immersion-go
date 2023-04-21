package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	u := new(User)

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the password
	err = u.SetPassword(u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the User
	err = u.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the User
	id, err := Insert(h.db, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Insert(db *sql.DB, u *User) (int64, error) {
	stmt := "INSERT INTO users(name, login, password, modified_at) VALUES($1, $2, $3, $4) RETURNING id"
	result, err := db.Exec(stmt, u.Name, u.Login, u.Password, u.ModifiedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
