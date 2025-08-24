package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var pages = map[string]*template.Template{}

func mustPage(page string) *template.Template {
	t, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", page),
	)
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}
	return t
}

func main() {
	// Подготовим шаблоны один раз при старте
	pages["index.html"] = mustPage("index.html")
	pages["math.html"] = mustPage("math.html")
	pages["english.html"] = mustPage("english.html")
	pages["about.html"] = mustPage("about.html")

	mux := http.NewServeMux()

	// Главная
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundHandler(w, r)
			return
		}
		render(w, "index.html", nil)
	})

	// Математика
	mux.HandleFunc("/math", func(w http.ResponseWriter, r *http.Request) {
		render(w, "math.html", nil)
	})
	// Английский язык
	mux.HandleFunc("/english", func(w http.ResponseWriter, r *http.Request) {
		render(w, "english.html", nil)
	})

	// O нас
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		render(w, "about.html", nil)
	})

	// Статика
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// favicon.ico
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("static", "favicon.ico")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			notFoundHandler(w, r)
			return
		}
		http.ServeFile(w, r, path)
	})

	// Неизвестные маршруты -> 404
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			notFoundHandler(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func render(w http.ResponseWriter, page string, data any) {
	t := pages[page]
	if t == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	// Рендерим каркас layout.html, в котором подставятся блоки из конкретной страницы
	if err := t.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Кастомный 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("static", "404.html")
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, path)
}
