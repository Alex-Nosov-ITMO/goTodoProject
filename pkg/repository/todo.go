package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)


type TodoRepository struct {
	Db *sqlx.DB
}

// Конструктор связки репозитория с БД
func NewTodoRepository(db *sqlx.DB) *TodoRepository {
	return &TodoRepository{Db: db}
}

// Методы репозитория
func (r *TodoRepository) GetTasks() ([]structures.Task,  error) {
	var tasks []structures.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit`
	
	rows, err := r.Db.Query(query, sql.Named("limit", 45))
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач из базы данных: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var task structures.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при считывании задач из базы данных: %v", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TodoRepository) CreateTask(task *structures.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
	response, err := r.Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании задачи в базе данных: %v", err)
	}

	id, err := response.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении id последней добавленной задачи: %v", err)
	}

	return id,nil
}

func (r *TodoRepository) DelTask(id int64) (error) {
	chekTask, err := r.GetTask(int64(id))
	if err != nil {
		return fmt.Errorf("ошибка при получении задачи из базы данных: %v", err)
	}
	if chekTask == (structures.Task{}){
		return fmt.Errorf("задача с id: %d не существует", id)
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err = r.Db.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи из базы данных: %v", err)
	}
	return nil
}

func (r *TodoRepository) GetTask(id int64) (structures.Task,  error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id"

	var task structures.Task
	row := r.Db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return structures.Task{}, fmt.Errorf("ошибка при получении задачи из базы данных: %v", err)
	}
	return task, nil
}


func (r *TodoRepository) UpdateTask(task *structures.Task) (error) {
	id, _ := strconv.Atoi(task.ID)	
	
	chekTask, err := r.GetTask(int64(id))
	if err != nil {
		return fmt.Errorf("ошибка при получении задачи из базы данных: %v", err)
	}
	if chekTask == (structures.Task{}){
		return fmt.Errorf("задача с id: %s не существует", task.ID)
	}
	
	query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"

	_, err = r.Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи в базе данных: %v", err)
	}
	return nil
}

