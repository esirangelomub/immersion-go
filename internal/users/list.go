package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	// Get User
	list, err := SelectAll(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SelectAll(db *sql.DB) ([]User, error) {
	stmt := "SELECT * FROM users WHERE deleted=false"
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.CreatedAt, &u.ModifiedAt, &u.Deleted, &u.LastLogin)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}
