package cli

import (
	"fmt"
	"os"
	"task-tracker/internal/tasks"
)

func RunCLI(service *tasks.Service) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-tracker <add|list> [name]")
		return
	}
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Provide task name")
			return
		}
		err := service.AddTask(os.Args[2])
		if err != nil {
			return
		}
	case "list":

		listTasks, err := service.ListTasks()
		if err != nil {
			return
		}
		if len(listTasks) == 0 {
			fmt.Println("No tasks found")
		} else {
			fmt.Println("Task List:")
			for _, task := range listTasks {
				fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated at: %s\nUpdated at: %s\n\n",
					task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
			}
		}
	case "help":
		fmt.Println("Task Tracker is a tool to write down and keep track of the things you need to do.\n\nUsage: \n\n        task-cli <command> [arguments]\n\nThe commands are:\n\n        add [\"name\"]                    add a task\n        list [todo/in-progress/done]    list all/todo/in-progress/done tasks\n        update [\"name\" ID]              update task name\n        delete [ID]                     delete a task\n        mark-in-progress [ID]           mark the status of the task as in progress\n        mark-done [ID]                  mark the status of the task as done\n        help                            show this message")
	default:
		fmt.Println("Task Tracker is a tool to write down and keep track of the things you need to do.\n\nUsage: \n\n        task-cli <command> [arguments]\n\nThe commands are:\n\n        add [\"name\"]                    add a task\n        list [todo/in-progress/done]    list all/todo/in-progress/done tasks\n        update [\"name\" ID]              update task name\n        delete [ID]                     delete a task\n        mark-in-progress [ID]           mark the status of the task as in progress\n        mark-done [ID]                  mark the status of the task as done\n        help                            show this message")
	}
}
