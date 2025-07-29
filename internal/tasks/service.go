package tasks

import (
	"fmt"
	"task-tracker/internal/models"
	"time"
)

type Service struct {
	storage Storage
}

type Storage interface {
	AddTask(task models.DataBase) error
	GetTasks() ([]models.DataBase, error)
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddTask(description string) error {
	if description != "" {
		currentTasks, err := s.storage.GetTasks()
		if err != nil {
			return err
		}
		maxID := 0
		for _, task := range currentTasks {
			if task.ID > maxID {
				maxID = task.ID
			}
		}
		newTask := models.DataBase{
			ID:          maxID + 1,
			Description: description,
			Status:      "todo",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = s.storage.AddTask(newTask)
		if err != nil {
			return err
		}
		fmt.Printf("Task added successfully (ID: %d)\n", maxID)
	} else {
		return fmt.Errorf("task description cannot be empty")
	}
	return nil
}

func (s *Service) ListTasks() ([]models.DataBase, error) {
	return s.storage.GetTasks()
}

// todo update command
// todo delete command
// todo mark-done command
// todo in-progress command
// todo help command
