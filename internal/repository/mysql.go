package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db *sql.DB
}

const (
	allTasks        = "SELECT * FROM task;"
	taskByID        = "SELECT * FROM task WHERE id = ?;"
	createTask      = "INSERT INTO task (name, completed) VALUES (?, ?);"
	updateTask      = "UPDATE task SET name = ?, completed = ? WHERE id = ?;"
	taskByCompleted = "SELECT * FROM task WHERE completed = ?;"
)

func NewMySQLRepository(cfg DBConfig) (*MySQLRepository, error) {
	c := mysql.Config{
		User:   cfg.User,
		Passwd: cfg.Pass,
		Net:    cfg.Protocol,
		Addr:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DBName: cfg.DBName,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open DB. Error:\n%s", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB. Error:\n%s", err)
	}

	return &MySQLRepository{db}, nil
}

func (mr *MySQLRepository) ListTasks() ([]Task, error) {
	var tasks []Task

	rows, err := mr.db.Query(allTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks. Error:\n%s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks. Error:\n%v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list tasks. Error:\n%v", err)
	}

	return tasks, nil
}

func (mr *MySQLRepository) GetTaskByID(ID int64) (Task, error) {
	var task Task

	row := mr.db.QueryRow(taskByID, ID)

	err := row.Scan(&task.ID, &task.Name, &task.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("task(%d) not found", ID)
		}
		return task, fmt.Errorf("failed to get a task by id: %d. Error:\n%v", ID, err)
	}
	return task, nil
}

func (mr *MySQLRepository) GetTasksByCompleted(completed bool) ([]Task, error) {
	var tasks []Task

	rows, err := mr.db.Query(taskByCompleted, completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by completed=%v. Error:\n%s", completed, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, fmt.Errorf("failed to get tasks by completed=%v. Error:\n%v", completed, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get tasks by completed=%v. Error:\n%v", completed, err)
	}

	return tasks, nil
}

func (mr *MySQLRepository) CreateTask(task Task) (int64, error) {
	result, err := mr.db.Exec(createTask, task.Name, task.Completed)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %+v. Error:\n%v", task, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %+v. Error:\n%v", task, err)
	}
	return id, nil
}

func (mr *MySQLRepository) UpdateTask(id int64, task Task) (int64, error) {
	result, err := mr.db.Exec(updateTask, task.Name, task.Completed, id)
	if err != nil {
		return 0, fmt.Errorf("failed to update task: %+v. Error:\n%v", task, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to update task: %+v. Error:\n%v", task, err)
	}

	if affected == 0 {
		return 0, fmt.Errorf("failed to update task: %+v. Error:\ntask(%d) not found", task, id)
	}

	return id, nil
}
