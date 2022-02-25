package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

// TO-DOs
// Cover edge cases, and mainly, error handling

func TestListTasks(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := MySQLRepository{db: db}
	columns := []string{"id", "name", "completed"}

	cases := []struct {
		name string
		mock func(tasks []Task)
		want []Task
	}{
		{
			name: "List tasks",
			mock: func(tasks []Task) {
				rows := sqlmock.NewRows(columns)

				for _, task := range tasks {
					rows.AddRow(task.ID, task.Name, task.Completed)
				}

				mock.ExpectQuery("SELECT (.+) FROM task").
					WillReturnRows(rows)
			},
			want: []Task{
				{ID: 42, Name: "Test", Completed: false},
				{ID: 101, Name: "Another task", Completed: true},
			},
		},
		// TODO: add more test scenarios
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock(c.want)
			got, err := repo.ListTasks()
			if err != nil {
				t.Errorf("got err: %q", err.Error())
				return
			}

			diff := cmp.Diff(got, c.want)
			if diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestGetTasksByCompleted(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := MySQLRepository{db: db}
	columns := []string{"id", "name", "completed"}

	cases := []struct {
		name      string
		completed bool
		mock      func(completed bool, tasks []Task)
		want      []Task
	}{
		{
			name:      "List completed tasks",
			completed: true,
			mock: func(completed bool, tasks []Task) {
				rows := sqlmock.NewRows(columns)

				for _, task := range tasks {
					rows.AddRow(task.ID, task.Name, task.Completed)
				}

				mock.ExpectQuery("SELECT (.+) FROM task").
					WithArgs(completed).
					WillReturnRows(rows)
			},
			want: []Task{
				{ID: 42, Name: "Test", Completed: true},
				{ID: 101, Name: "Another task", Completed: true},
			},
		},
		{
			name:      "List incompleted tasks",
			completed: false,
			mock: func(completed bool, tasks []Task) {
				rows := sqlmock.NewRows(columns)

				for _, task := range tasks {
					rows.AddRow(task.ID, task.Name, task.Completed)
				}

				mock.ExpectQuery("SELECT (.+) FROM task").
					WithArgs(completed).
					WillReturnRows(rows)
			},
			want: []Task{
				{ID: 42, Name: "Test", Completed: false},
				{ID: 101, Name: "Another task", Completed: false},
			},
		},
		// TODO: add more test scenarios
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock(c.completed, c.want)
			got, err := repo.GetTasksByCompleted(c.completed)
			if err != nil {
				t.Errorf("got err: %q", err.Error())
				return
			}

			diff := cmp.Diff(got, c.want)
			if diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := MySQLRepository{db: db}
	columns := []string{"id", "name", "completed"}

	cases := []struct {
		name string
		ID   int64
		mock func(ID int64, want Task)
		want Task
	}{
		{
			name: "Get task 42",
			ID:   42,
			mock: func(ID int64, want Task) {
				mock.ExpectQuery("SELECT (.+) FROM task").
					WithArgs(ID).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(want.ID, want.Name, want.Completed))
			},
			want: Task{ID: 42, Name: "Test", Completed: false},
		},
		// TODO: add more test scenarios
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock(c.ID, c.want)
			got, err := repo.GetTaskByID(c.ID)
			if err != nil {
				t.Errorf("got err: %q", err.Error())
				return
			}

			diff := cmp.Diff(got, c.want)
			if diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := MySQLRepository{db: db}

	cases := []struct {
		name string
		task Task
		mock func(task Task, want int64)
		want int64
	}{
		{
			name: "Create simple task successfully",
			task: Task{Name: "Test", Completed: false},
			mock: func(task Task, want int64) {
				mock.ExpectExec("INSERT INTO task").
					WithArgs(task.Name, task.Completed).
					WillReturnResult(sqlmock.NewResult(want, 1))
			},
			want: 42,
		},
		// TODO: add more test scenarios
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock(c.task, c.want)
			got, err := repo.CreateTask(c.task)
			if err != nil {
				t.Errorf("got err: %q", err.Error())
				return
			}

			if got != c.want {
				t.Errorf("got %d, want %d", got, c.want)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := MySQLRepository{db: db}

	cases := []struct {
		name string
		task Task
		mock func(task Task)
		want int64
	}{
		{
			name: "Create simple task successfully",
			task: Task{ID: 42, Name: "New Name", Completed: true},
			mock: func(task Task) {
				mock.ExpectExec("UPDATE task SET name = (.+), completed = (.+) WHERE id = (.+)").
					WithArgs(task.Name, task.Completed, task.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want: 42,
		},
		// TODO: add more test scenarios
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock(c.task)
			got, err := repo.UpdateTask(c.task)
			if err != nil {
				t.Errorf("got err: %q", err.Error())
				return
			}

			if got != c.want {
				t.Errorf("got %d, want %d", got, c.want)
			}
		})
	}
}
