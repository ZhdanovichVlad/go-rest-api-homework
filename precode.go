package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

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
// 1-ый обработчик выводит все задачи
func getTasks(w http.ResponseWriter, req *http.Request) {
	keys := make([]string, 0, len(tasks))
	for k, _ := range tasks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var respTotal []Task
	for i := 0; i < len(keys); i++ {
		respTotal = append(respTotal, tasks[keys[i]])

	}

	resp, err := json.Marshal(respTotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error w.Write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// 2-ой обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	if task.ID == "" {
		http.Error(w, "Не указано ID задачи", http.StatusBadRequest)
		fmt.Println("Не указано ID задачи")
		return

	}

	_, exist := tasks[task.ID]
	if exist {
		http.Error(w, "Задача c таким ID уже существует", http.StatusBadRequest)
		fmt.Println("Задача c таким ID уже существует")
		return
	}

	if len(task.Applications) == 0 {
		task.Applications = append(task.Applications, req.Header.Get("User-Agent"))
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// 3-ий обработчик выводит задачу по заданному ID
func getTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Запрашиваемая задача не найдена", http.StatusBadRequest)
		fmt.Println("Запрашиваемая задача не найдена")
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error w.Write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// 4-тый обработчик удаляет задачу
func deleteTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Удаляемая задача не найдена", http.StatusBadRequest)
		fmt.Println("Удаляемая задача не найдена")
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()
	// зарегистрированные обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
