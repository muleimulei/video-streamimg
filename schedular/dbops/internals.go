package dbops

import (
	"database/sql"
	"log"
)

func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmt, err := db.Prepare("SELECT video_id from video_del_rec LIMIT ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Query VideoDeletion error ", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	var ids []string
	for rows.Next() {
		var vid string
		if err := rows.Scan(&vid); err != nil {
			return ids, err
		}
		ids = append(ids, vid)
	}
	return ids, nil
}

func DelVideoDeletion(vid string) error {
	stmt, err := db.Prepare("DELETE FROM video_del_rec WHERE video_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vid)
	if err != nil {
		log.Println("Deleting video error : ", err)
		return err
	}

	return nil
}
