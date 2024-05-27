package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/nextDate"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/repository"
)

type TodoService struct{
	rps repository.RepositoryInterface
}

// Конструктор связки сервиса и репозитория
func NewTodoService(repos repository.RepositoryInterface) *TodoService {
	return &TodoService{rps: repos}
}

// Методы сервиса
func (s *TodoService) GetTasks(search string) (task []structures.Task,  err error) {
	var tasks []structures.Task 

	if search == "" {
		tasks, err = s.rps.GetTasks()
		if err != nil {
			return nil, err
		}
	}else {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			tasks, err = s.rps.GetTasksWithDate(date.Format("20060102"))
			if err != nil {
				return nil, err
			}
		} else {
			tasks, err = s.rps.GetTasksWithStr(search)
			if err != nil {
				return nil, err
			}
		}	
	}

	if len(tasks) == 0 {
		return []structures.Task{}, nil
	} 
	
	return tasks, nil 
}


func (s *TodoService) CreateTask(task *structures.Task) (int64, error) {

	if task == nil {
		return 0, errors.New("task is nil")
	}

	if task.Title == "" {
		return 0, errors.New("title is empty")
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Service: CreateTask: невозможно преобразовать дату: %s. Ошибка: %v", task.Date, err)
		return 0, errors.New("невозможно преобразовать дату")
		}


	if task.Date < time.Now().Format("20060102") {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			task.Date, err = nextDate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
		}
	}
	return s.rps.CreateTask(task)
}

func (s *TodoService) DelTask(id int64) (err error) {
	return s.rps.DelTask(id)
}

func (s *TodoService) GetTask(id int64) (task structures.Task,  err error) {
	return s.rps.GetTask(id)
}

func (s *TodoService) UpdateTask(task *structures.Task) (err error) {
	if task == nil {
		return errors.New("task is nil")
	}

	if task.ID == "" {
		return errors.New("id is empty")
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		log.Printf("Service: UpdateTask: невозможно преобразовать id: %s. Ошибка: %v", task.ID, err)
		return errors.New("невозможно преобразовать id")
	}
	
	if task.Title == "" {
		return errors.New("title is empty")
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Service: UpdateTask: невозможно преобразовать дату: %s. Ошибка: %v", task.Date, err)
		return errors.New("невозможно преобразовать дату")
	}


	if task.Date < time.Now().Format("20060102") {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			task.Date, err = nextDate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return err
			}
		}
	}

	return s.rps.UpdateTask(task)
}