package main

import (
	"log"
	"task-tracker/internal/cli"
	"task-tracker/internal/storage"
	"task-tracker/internal/tasks"
)

func main() {
	jsonStorage, err := storage.NewJSONStorage("data.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	cli.RunCLI(tasks.NewService(jsonStorage))
}
