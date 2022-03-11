package repository

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (Task) TableName() string {
	return "task"
}

type OrmRepository struct {
	db *gorm.DB
}

func NewOrmRepository(cfg DBConfig) (*OrmRepository, error) {
	t := "%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(t, cfg.User, cfg.Pass, cfg.Protocol, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB. Error:\n%s", err)
	}

	return &OrmRepository{db}, nil
}

func (or *OrmRepository) ListTasks() ([]Task, error) {
	var tasks []Task

	result := or.db.Find(&tasks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list tasks. Error:\n%s", result.Error)
	}

	return tasks, nil
}

func (or *OrmRepository) GetTaskByID(ID int64) (Task, error) {
	var task Task

	result := or.db.First(&task, ID)

	if result.Error != nil {
		return task, fmt.Errorf("failed to get a task by id: %d. Error:\n%v", ID, result.Error)
	}

	return task, nil
}

func (or *OrmRepository) GetTasksByCompleted(completed bool) ([]Task, error) {
	var tasks []Task

	result := or.db.Where("Completed = ?", completed).Find(&tasks)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get tasks by completed=%v. Error:\n%s", completed, result.Error)
	}

	return tasks, nil
}

func (or *OrmRepository) CreateTask(task Task) (int64, error) {
	result := or.db.Create(&task)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create task: %+v. Error:\n%v", task, result.Error)
	}

	id := task.ID
	return id, nil
}

func (or *OrmRepository) UpdateTask(id int64, task Task) (int64, error) {
	toUpdate := Task{ID: id}
	or.db.First(&toUpdate)

	toUpdate.Name = task.Name
	toUpdate.Completed = task.Completed

	result := or.db.Save(&toUpdate)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to update task: %+v. Error:\n%v", task, result.Error)
	}

	affected := result.RowsAffected
	if affected == 0 {
		return 0, fmt.Errorf("failed to update task: %+v. Error:\ntask(%d) not found", task, id)
	}

	return id, nil
}
