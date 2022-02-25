package main

import (
	"fmt"
	"os"

	"github.com/eoliveira-daitan/go-challenges/internal/repository"
)

func main() {
	var repo repository.Repository

	cfg := repository.MysqlConfig{
		User:   os.Getenv("DBUSER"),
		Pass:   os.Getenv("DBPASS"),
		Host:   "127.0.0.1",
		Port:   os.Getenv("DBPORT"),
		DBName: os.Getenv("DBNAME"),
	}

	repo, err := repository.NewMySQLRepository(cfg)
	handleErr(err)

	taskID, err := repo.CreateTask(repository.Task{Name: "Create a dummy task", Completed: false})
	handleErr(err)
	fmt.Printf("Task created with ID: %d\n", taskID)

	task, err := repo.GetTaskByID(taskID)
	handleErr(err)
	fmt.Printf("Task found: %+v\n", task)

	task.Completed = true
	taskID, err = repo.UpdateTask(task)
	handleErr(err)
	fmt.Printf("Task updated: %d\n", taskID)

	anotherTaskID, err := repo.CreateTask(repository.Task{Name: "Create a second task", Completed: false})
	handleErr(err)
	fmt.Printf("Task created with ID: %d\n", anotherTaskID)

	tasks, err := repo.ListTasks()
	handleErr(err)
	fmt.Printf("Tasks found: %+v\n", tasks)

	tasks, err = repo.GetTasksByCompleted(true)
	handleErr(err)
	fmt.Printf("Tasks completed: %+v\n", tasks)

	tasks, err = repo.GetTasksByCompleted(false)
	handleErr(err)
	fmt.Printf("Tasks incompleted: %+v\n", tasks)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
