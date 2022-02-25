package repository

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

type Repository interface {
	ListTasks() ([]Task, error)
	GetTaskByID(ID int64) (Task, error)
	GetTasksByCompleted(completed bool) ([]Task, error)
	CreateTask(task Task) (int64, error)
	UpdateTask(task Task) (int64, error)
}
