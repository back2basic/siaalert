package sdk

import (
	"fmt"
	"sync"

	"github.com/back2basic/siadata/siaalert/strict"
)

func SdkWorker(id int, tasks <-chan strict.TaskCheckDoc, wg *sync.WaitGroup) {
	db := GetAppwriteDatabaseService()
	defer func() {
		wg.Done()
	}()
	for task := range tasks {
		switch task.Job {
		case "createCheck":
			// fmt.Printf("WIP %d: Creating Check %d\n", id, task.ID)
			_, err := db.CreateCheckDocument(task.CheckID, task.Check)
			if err != nil {
				fmt.Println(err)
			}
		case "updateCheck":
			// fmt.Printf("WIP %d: Updating Check %d\n", id, task.ID)
			_, err := db.UpdateCheckDocument(task.CheckID, task.Check)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("WIP %d: Unknown task: %s\n", id, task.Job)
		}
	}
}
