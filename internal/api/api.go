package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/eoliveira-daitan/go-challenges/internal/repository"
)

type Server struct {
	Store repository.Repository
	http.Handler
}

func New(store repository.Repository) *Server {
	s := new(Server)
	s.Store = store

	router := http.NewServeMux()

	router.Handle("/tasks", http.HandlerFunc(s.HandleTasks))
	router.Handle("/tasks/", http.HandlerFunc(s.HandleTasks))

	s.Handler = router

	return s
}

func parseTask(body io.Reader) (task repository.Task, err error) {
	err = json.NewDecoder(body).Decode(&task)
	return task, err
}

func handleResponse(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(body))
	if err != nil {
		fmt.Printf("failed to write the failed payload: %v\n", err)
	}
}

func (s *Server) HandleTasks(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.Contains(path, "/tasks") {
		handleResponse(w, http.StatusNotImplemented, fmt.Sprintf("%q is not supported", path))
		return
	}

	switch r.Method {
	case "POST":
		task, err := parseTask(r.Body)
		if err != nil {
			handleResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		s.CreateTask(w, task)
	case "PUT":
		byIDPath := regexp.MustCompile("^/tasks/([0-9]+)$")
		m := byIDPath.MatchString(path)
		if !m {
			handleResponse(w, http.StatusNotImplemented, fmt.Sprintf("%q is not supported", path))
			return
		}

		str := byIDPath.FindStringSubmatch(path)[1]
		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			handleResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid id: %v", str))
			return
		}

		task, err := parseTask(r.Body)
		if err != nil {
			handleResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		s.UpdateTask(w, id, task)
	case "GET":
		byIDPath := regexp.MustCompile("^/tasks/([0-9]+)$")
		taskByID := byIDPath.MatchString(path)
		if taskByID {
			str := byIDPath.FindStringSubmatch(path)[1]
			id, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				handleResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid id: %v", str))
				return
			}

			s.GetTaskByID(w, id)
			return
		}

		c := r.URL.Query().Get("completed")
		if c != "" {
			var completed bool

			if c == "true" {
				completed = true
			}

			s.GetTasksByCompleted(w, completed)
			return
		}

		s.GetTasks(w)
	}

}

func (s *Server) CreateTask(w http.ResponseWriter, task repository.Task) {
	id, err := s.Store.CreateTask(task)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(w, http.StatusCreated, fmt.Sprint(id))
}

func (s *Server) UpdateTask(w http.ResponseWriter, id int64, task repository.Task) {
	id, err := s.Store.UpdateTask(id, task)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, fmt.Sprint(id))
}

func (s *Server) GetTaskByID(w http.ResponseWriter, id int64) {
	task, err := s.Store.GetTaskByID(id)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := json.Marshal(task)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, string(body))
}

func (s *Server) GetTasks(w http.ResponseWriter) {
	tasks, err := s.Store.ListTasks()
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := json.Marshal(tasks)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, string(body))
}

func (s *Server) GetTasksByCompleted(w http.ResponseWriter, completed bool) {
	tasks, err := s.Store.GetTasksByCompleted(completed)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := json.Marshal(tasks)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, string(body))
}
