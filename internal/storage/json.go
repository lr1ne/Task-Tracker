package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"task-tracker/internal/models"
)

type JSONStorage struct {
	filename string
	tasks    []models.DataBase
}

func NewJSONStorage(filename string) (*JSONStorage, error) {
	err := ensureFileExists(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var tasks []models.DataBase
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return &JSONStorage{filename: filename, tasks: tasks}, nil
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

		_, err = file.Write([]byte("[]"))
		if err != nil {
			return err
		}
		fmt.Println("File create: ", filename)
	}
	return nil
}

func (s *JSONStorage) AddTask(task models.DataBase) error {
	s.tasks = append(s.tasks, task)
	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.tasks)
	if err != nil {
		return err
	}
	return nil
}

func (s *JSONStorage) GetTasks() ([]models.DataBase, error) {
	var tasks []models.DataBase
	if s.tasks != nil {
		return s.tasks, nil
	}
	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	s.tasks = tasks
	return tasks, nil
}
