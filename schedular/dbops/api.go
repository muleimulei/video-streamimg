package dbops

import "log"

func AddVideoDeletionRecord(vid string) error {
	stmt, err := db.Prepare("INSERT INTO video_del_rec (video_id) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(vid)
	if err != nil {
		log.Println("AddVideoDeletion error ", err)
		return err
	}
	return nil
}
