package sdk

import (
	"fmt"
	"sync"
)

type Task struct {
	ID      int
	Job     string
	CheckID string
	Check   Check
}

func Worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	db := GetAppwriteDatabaseService()
	for task := range tasks {
		// Simulate task processing
		// fmt.Printf("WIP %d: Processing Task %d - %s\n", id, task.ID, task.Job)
		// time.Sleep(1 * time.Second) // Simulate work time

		switch task.Job {
		case "createCheck":
			fmt.Printf("WIP %d: Creating Check %d\n", id, task.ID)
			_, err := db.CreateCheckDocument(task.CheckID, task.Check)
			if err != nil {
				fmt.Println(err)
			}
		case "updateCheck":
			fmt.Printf("WIP %d: Updating Check %d\n", id, task.ID)
			_, err := db.UpdateCheckDocument(task.CheckID, task.Check)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("WIP %d: Unknown task: %s\n", id, task.Job)
		}

		wg.Done() // Mark the task as done
	}
}
