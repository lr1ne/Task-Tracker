package main

import (
	"encoding/json"
	"fmt"
	"io"
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

//func createIfNotExists(filename string) error {
//	_, err := os.Stat(filename)
//	if os.IsNotExist(err) {
//		file, err := os.Create(filename)
//		if err != nil {
//			return fmt.Errorf("не удалось создать файл: %v", err)
//		}
//		defer file.Close()
//		fmt.Printf("Файл %s создан\n", filename)
//		return nil
//	}
//	fmt.Printf("Файл %s уже существует\n", filename)
//	return nil
//}

func main() {
	//if err := createIfNotExists("data.json"); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var tasks []DataBase
	err = json.Unmarshal(byteValue, &tasks)
	if err != nil {
		log.Fatalf("JSON parsing error: %v", err)
	}

	if len(os.Args) > 0 {
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
				for _, task := range tasks {
					fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
						task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
				}
			}
		case "add":
			fmt.Println("Test Test")
		case "update":
			// todo
		case "delete":
			// todo
		case "":
		default:
			fmt.Println("Нету такова")
		}
	}
}
