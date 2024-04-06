package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже описаны обработчики для каждого эндпоинта
// 1-ый обработчки выводи все таски
func getTasks(w http.ResponseWriter, req *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//            ("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// 2-ой Обработчик для отправки задачи на сервер
func postTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	//            ("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// 3-ий обработчки выводит таск по заданному ID
func getTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не нейдена", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//            ("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// 4-тый обработчик удаляет таск
func deletTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не нейдена", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	//            ("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()
	// зарегестрированные обработчики
	r.Get("/tasks", getTasks)
	r.Post("/task", postTask)
	r.Get("/task/{id}", getTask)
	r.Delete("/task/{id}", deletTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
