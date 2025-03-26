package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type DataBase struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ensureFileExists(filename string) error {
	if !fileExists(filename) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// Записываем пустой JSON массив, чтобы избежать ошибок при парсинге
		_, err = file.Write([]byte("[]"))
		if err != nil {
			return err
		}
		fmt.Println("File create: ", filename)
	}
	return nil
}

func main() {
	filename := "data.json"

	err := ensureFileExists(filename)
	if err != nil {
		log.Fatalf("File create error: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("File open error: %v", err)
	}
	defer file.Close()

	var tasks []DataBase
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatalf("JSON decode error: %v", err)
	}

	if len(os.Args) == 1 {
		fmt.Println("Task Tracker is a tool to write down and keep track of the things you need to do.\n\nUsage: \n\n        task-cli <command> [arguments]\n\nThe commands are:\n\n        add [\"name\"]                    add a task\n        list [todo/in-progress/done]    list all/todo/in-progress/done tasks\n        update [\"name\" ID]              update task name\n        delete [ID]                     delete a task\n        mark-in-progress [ID]           mark the status of the task as in progress\n        mark-done [ID]                  mark the status of the task as done\n        help                            show this message")
	} else {
		switch os.Args[1] {
		case "list":
			if len(os.Args) > 2 {
				found := false

				if os.Args[2] == "todo" {
					fmt.Println("All todo tasks: ")
					for _, task := range tasks {
						if task.Status == "todo" {
							found = true
							fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
								task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
						}
					}
				} else if os.Args[2] == "in-progress" {
					fmt.Println("All in progress tasks: ")
					for _, task := range tasks {
						if task.Status == "in-progress" {
							found = true
							fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
								task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
						}
					}
				} else if os.Args[2] == "done" {
					fmt.Println("All done tasks: ")
					for _, task := range tasks {
						if task.Status == "done" {
							found = true
							fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
								task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
						}
					}
				}
				if !found {
					fmt.Printf("No tasks with \"%s\".\n", os.Args[2])
				}
			} else {
				if len(tasks) == 0 {
					fmt.Println("No tasks found.")
					return
				} else {
					for _, task := range tasks {
						fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
							task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
					}
				}
			}
		case "add":
			if len(os.Args) > 2 {
				var maxID int
				for _, task := range tasks {
					if task.ID > maxID {
						maxID = task.ID
					}
				}
				ID := maxID + 1
				newTask := DataBase{
					ID:          ID,
					Description: os.Args[2],
					Status:      "todo",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				tasks = append(tasks, newTask)

				file, err := os.Create("data.json")
				if err != nil {
					fmt.Println("Error when opening a file for recording:", err)
					return
				}
				defer file.Close()

				err = json.NewEncoder(file).Encode(tasks)
				if err != nil {
					fmt.Println("JSON encoding error:", err)
					return
				}

				fmt.Printf("Task added successfully (ID: %d)\n", ID)
			} else {
				fmt.Println("Use <add \"Your new task\"> ")
			}
		case "update":
			if len(os.Args) > 2 {
				var updateID int
				_, err := fmt.Sscanf(os.Args[2], "%d", &updateID)
				if err != nil {
					fmt.Println("ID parsing error: ", err)
					return
				}

				var updatedTasks []DataBase
				for _, task := range tasks {
					if task.ID == updateID {
						task.Description = os.Args[2]
						updatedTasks = append(updatedTasks, task)
						break
					}
				}
			} else {
				fmt.Println("Use <update ID>")
			}
		case "delete":
			if len(os.Args) > 2 {
				var deletedId int
				_, err := fmt.Sscanf(os.Args[2], "%d", &deletedId)
				if err != nil {
					fmt.Println("ID parsing error: ", err)
					return
				}

				found := false
				var updatedTasks []DataBase
				for _, task := range tasks {
					if task.ID != deletedId {
						updatedTasks = append(updatedTasks, task)
					} else {
						found = true
					}
				}

				if !found {
					fmt.Printf("Error: task with id %d not found\n", deletedId)
					return
				}

				file, err := os.Create("data.json")
				if err != nil {
					fmt.Println("Error when opening a file for recording: ", err)
					return
				}
				defer file.Close()

				err = json.NewEncoder(file).Encode(updatedTasks)
				if err != nil {
					fmt.Println("JSON encoding error:", err)
					return
				}
				fmt.Println("Delete successfully!")
			} else {
				fmt.Println("Use <delete ID>")
			}
		case "mark-in-progress":
			if len(os.Args) > 2 {
				var inProgressId int
				_, err := fmt.Sscanf(os.Args[2], "%d", &inProgressId)
				if err != nil {
					fmt.Println("ID parsing error: ", err)
				}

				var updatedTasks []DataBase
				for _, task := range tasks {
					if task.ID == inProgressId {
						task.Status = "in-progress"
						updatedTasks = append(updatedTasks, task)
						break
					}
				}

				file, err := os.Create("data.json")
				if err != nil {
					fmt.Println("Error when opening a file for recording: ", err)
					return
				}
				defer file.Close()

				err = json.NewEncoder(file).Encode(updatedTasks)
				if err != nil {
					fmt.Println("JSON encoding error: ", err)
					return
				}

			} else {
				fmt.Println("Use <mark-in-progress ID>")
			}
		case "mark-done":
			if len(os.Args) > 2 {
				var markDoneID int
				_, err := fmt.Sscanf(os.Args[2], "%d", &markDoneID)
				if err != nil {
					fmt.Println("ID parsing error: ", err)
				}

				var updatedTasks []DataBase
				for _, task := range tasks {
					if task.ID == markDoneID {
						task.Status = "done"
						updatedTasks = append(updatedTasks, task)
						break
					}
				}

				file, err := os.Create("data.json")
				if err != nil {
					fmt.Println("Error when opening a file for recording: ", err)
					return
				}
				defer file.Close()

				err = json.NewEncoder(file).Encode(updatedTasks)
				if err != nil {
					fmt.Println("JSON encoding error: ", err)
					return
				}

			} else {
				fmt.Println("Use <mark-done ID>")
			}
		case "help":
			fmt.Println("Task Tracker is a tool to write down and keep track of the things you need to do.\n\nUsage: \n\n        task-cli <command> [arguments]\n\nThe commands are:\n\n        add [\"name\"]                    add a task\n        list [todo/in-progress/done]    list all/todo/in-progress/done tasks\n        update [\"name\" ID]              update task name\n        delete [ID]                     delete a task\n        mark-in-progress [ID]           mark the status of the task as in progress\n        mark-done [ID]                  mark the status of the task as done\n        help                            show this message")
		default:
			fmt.Println("Task Tracker is a tool to write down and keep track of the things you need to do.\n\nUsage: \n\n        task-cli <command> [arguments]\n\nThe commands are:\n\n        add [\"name\"]                    add a task\n        list [todo/in-progress/done]    list all/todo/in-progress/done tasks\n        update [\"name\" ID]              update task name\n        delete [ID]                     delete a task\n        mark-in-progress [ID]           mark the status of the task as in progress\n        mark-done [ID]                  mark the status of the task as done\n        help                            show this message")
		}
	}
}
