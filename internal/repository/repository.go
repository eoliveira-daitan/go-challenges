package repository

type Task struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type Repository interface {
	ListTasks() ([]Task, error)
	GetTaskByID(id int64) (Task, error)
	GetTasksByCompleted(completed bool) ([]Task, error)
	CreateTask(task Task) (int64, error)
	UpdateTask(id int64, task Task) (int64, error)
}
