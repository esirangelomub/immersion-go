package folders

import (
	"database/sql"
	"github.com/esirangelomub/immersion-go/internal/files"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFiles(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
}

func Delete(db *sql.DB, id int64) error {
	stmt := `UPDATE folders SET modified_at = $1, deleted = true WHERE id = $2`
	_, err := db.Exec(stmt, time.Now(), id)

	return err
}

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	return deleteSubFolder(db, folderID)
}

func deleteSubFolder(db *sql.DB, folderID int64) error {
	subFolders, err := getSubFolder(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subFolders))
	for _, rf := range removedFolders {
		err := Delete(db, rf.ID)
		if err != nil {
			break
		}

		err = deleteFolderContent(db, rf.ID)
		if err != nil {
			Update(db, rf.ID, &rf)
			break
		}

		removedFolders = append(removedFolders, rf)
	}

	if len(subFolders) != len(removedFolders) {
		for _, rf := range removedFolders {
			Update(db, rf.ID, &rf)
		}
	}

	return nil
}

func deleteFiles(db *sql.DB, folderID int64) error {
	f, err := files.List(db, folderID)
	if err != nil {
		return err
	}

	removedFiles := make([]files.File, 0, len(f))
	for _, file := range f {
		file.Deleted = true
		err := files.Update(db, file.ID, &file)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(f) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			files.Update(db, file.ID, &file)
		}
		return err
	}

	return nil
}
