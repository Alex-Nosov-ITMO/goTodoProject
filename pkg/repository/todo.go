package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"


	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	myErr "github.com/Alex-Nosov-ITMO/go_project_final/internal/myErrors"
)

type TodoRepository struct {
	Db *sqlx.DB
}

// Конструктор связки репозитория с БД
func NewTodoRepository(db *sqlx.DB) *TodoRepository {
	return &TodoRepository{Db: db}
}

// Методы репозитория
func (r *TodoRepository) GetTasks() ([]structures.Task, error) {
	var tasks []structures.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit`

	rows, err := r.Db.Query(query, sql.Named("limit", 45))
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasks", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task structures.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, myErr.WithMassage("Repository: GetTasks", err)
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasks", err)
	}

	return tasks, nil
}

func (r *TodoRepository) GetTasksWithStr(str string) ([]structures.Task, error) {
	var tasks []structures.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :title OR comment LIKE :comment ORDER BY date ASC LIMIT :limit`

	rows, err := r.Db.Query(query, sql.Named("title", "%"+str+"%"), sql.Named("comment", "%"+str+"%"), sql.Named("limit", 45))
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasksWithStr", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task structures.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, myErr.WithMassage("Repository: GetTasksWithStr", err)
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasksWithStr", err)
	}

	return tasks, nil
}

func (r *TodoRepository) GetTasksWithDate(date string) ([]structures.Task, error) {
	var tasks []structures.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date ASC LIMIT :limit`

	rows, err := r.Db.Query(query, sql.Named("date", date), sql.Named("limit", 45))
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasksWithDate", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task structures.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, myErr.WithMassage("Repository: GetTasksWithDate", err)
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, myErr.WithMassage("Repository: GetTasksWithDate", err)
	}

	return tasks, nil
}

func (r *TodoRepository) CreateTask(task *structures.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
	response, err := r.Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, myErr.WithMassage("Repository: CreateTask", err)
	}

	id, err := response.LastInsertId()
	if err != nil {
		return 0, myErr.WithMassage("Repository: CreateTask", err)
	}

	return id, nil
}

func (r *TodoRepository) DelTask(id int64) error {
	chekTask, err := r.GetTask(int64(id))
	if err != nil {
		return myErr.WithMassage("Repository: DelTask", err)
	}
	if chekTask == (structures.Task{}) {
		return myErr.WithMassage("Repository: DelTask", fmt.Errorf("задача с id: %d не существует", id))
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err = r.Db.Exec(query, sql.Named("id", id))
	if err != nil {
		return myErr.WithMassage("Repository: DelTask", err)
	}
	return nil
}

func (r *TodoRepository) GetTask(id int64) (structures.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id"

	var task structures.Task
	row := r.Db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return structures.Task{}, myErr.WithMassage("Repository: GetTask", err)
	}
	return task, nil
}

func (r *TodoRepository) UpdateTask(task *structures.Task) error {
	id, _ := strconv.Atoi(task.ID)

	chekTask, err := r.GetTask(int64(id))
	if err != nil {
		return myErr.WithMassage("Repository: UpdateTask", err)
	}

	if chekTask == (structures.Task{}) {
		return myErr.WithMassage("Repository: UpdateTask", fmt.Errorf("задача с id: %s не существует", task.ID))
	}

	query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"

	_, err = r.Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return myErr.WithMassage("Repository: UpdateTask", err)
	}
	return nil
}
