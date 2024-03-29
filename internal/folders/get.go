package folders

import (
	"database/sql"
	"encoding/json"
	"github.com/esirangelomub/immersion-go/internal/files"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (h *handler) Get(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := GetFolder(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := GetFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := FolderContent{
		Folder:  *f,
		Content: c,
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(fc)
}

func GetFolder(db *sql.DB, folderID int64) (*Folder, error) {
	stmt := `SELECT id, parent_id, name, created_at, modified_at FROM folders WHERE id = $1`
	row := db.QueryRow(stmt, folderID)

	var f Folder
	err := row.Scan(&f.ID, &f.ParentID, &f.Name, &f.CreatedAt, &f.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetFolderContent(db *sql.DB, folderID int64) ([]FolderResource, error) {
	subFolders, err := getSubFolder(db, folderID)
	if err != nil {
		return nil, err
	}

	fr := make([]FolderResource, 0, len(subFolders))
	for _, sf := range subFolders {
		r := FolderResource{
			ID:         sf.ID,
			Name:       sf.Name,
			Type:       "directory",
			CreatedAt:  sf.CreatedAt,
			ModifiedAt: sf.ModifiedAt,
		}

		fr = append(fr, r)
	}

	folderFiles, err := files.List(db, folderID)
	if err != nil {
		return nil, err
	}

	for _, f := range folderFiles {
		r := FolderResource{
			ID:         f.ID,
			Name:       f.Name,
			Type:       f.Type,
			CreatedAt:  f.CreatedAt,
			ModifiedAt: f.ModifiedAt,
		}

		fr = append(fr, r)
	}

	return fr, nil
}

func getSubFolder(db *sql.DB, folderID int64) ([]Folder, error) {
	stmt := `SELECT id, parent_id, name, created_at, modified_at FROM folders WHERE parent_id = $1 AND default = false`
	rows, err := db.Query(stmt, folderID)
	if err != nil {
		return nil, err
	}

	f := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreatedAt, &folder.ModifiedAt)
		if err != nil {
			return nil, err
		}

		f = append(f, folder)
	}

	return f, nil
}
