package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

//go:embed templates/*
var templatesFS embed.FS

var templates *template.Template

// Simple in-memory todo storage
type Todo struct {
	ID   int
	Text string
	Done bool
}

var (
	todos   []Todo
	todosMu sync.RWMutex
	nextID  = 1
)

func main() {
	var err error
	templates, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}

	// Initialize with some example todos
	todos = []Todo{
		{ID: 1, Text: "Deploy this app on Dalang.io", Done: false},
		{ID: 2, Text: "Try HTMX interactions", Done: false},
		{ID: 3, Text: "Customize the brutalist theme", Done: true},
	}
	nextID = 4

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/todos", handleTodos)
	http.HandleFunc("/todos/add", handleAddTodo)
	http.HandleFunc("/todos/toggle/", handleToggleTodo)
	http.HandleFunc("/todos/delete/", handleDeleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	todosMu.RLock()
	defer todosMu.RUnlock()

	templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Todos": todos,
	})
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	todosMu.RLock()
	defer todosMu.RUnlock()

	templates.ExecuteTemplate(w, "todos.html", map[string]interface{}{
		"Todos": todos,
	})
}

func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	todos = append(todos, Todo{
		ID:   nextID,
		Text: text,
		Done: false,
	})
	nextID++
	todosMu.Unlock()

	handleTodos(w, r)
}

func handleToggleTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path /todos/toggle/{id}
	var id int
	_, err := fmt.Sscanf(r.URL.Path, "/todos/toggle/%d", &id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = !todos[i].Done
			break
		}
	}
	todosMu.Unlock()

	handleTodos(w, r)
}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path /todos/delete/{id}
	var id int
	_, err := fmt.Sscanf(r.URL.Path, "/todos/delete/%d", &id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	for i := range todos {
		if todos[i].ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	todosMu.Unlock()

	handleTodos(w, r)
}
