package files

import (
	"database/sql"
)

func List(db *sql.DB, folderID int64) ([]File, error) {
	stmt := `SELECT * FROM files WHERE folder_id = $1 AND deleted = false`
	rows, err := db.Query(stmt, folderID)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
	for rows.Next() {
		var f File

		err := rows.Scan(&f.ID, &f.FolderID, &f.OwnerID, &f.Name, &f.Type, &f.Path, &f.CreatedAt, &f.ModifiedAt, &f.Deleted)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}
