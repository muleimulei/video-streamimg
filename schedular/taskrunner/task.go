package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"video_server/schedular/dbops"
)

func deleteVideo(path string) error {
	err := os.Remove(VIDEO_DIR + path)
	if err != nil && os.IsExist(err) {
		log.Println("Deleting video error ", err)
		return err
	}
	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Println("video clear dispatcher error : ", err)
		return err
	}

	if res == nil {
		return errors.New("all tasks finished")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDeletion(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		return err == nil
	})
	return err
}
