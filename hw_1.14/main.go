package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var pages = map[string]*template.Template{}
var db *sql.DB

// структура для контактов
type ContactInfo struct {
	Phone          string
	Email          string
	GraficBudni    string
	GraficSaturday string
	GraficSunday   string
}

// загрузка шаблонов
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
	// Подключение к базе
	var err error
	db, err = sql.Open("sqlite", "./ucheb_center.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}
	defer db.Close()

	// Создать таблицы, если  ещё нет
	createTable := `
	CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phone TEXT,
		email TEXT,
		GraficBudni TEXT,
		GraficSaturday TEXT,
		GraficSunday TEXT
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}

	// Если таблица пустая – вставим данные
	_, err = db.Exec(`INSERT INTO contacts (phone, email, GraficBudni, GraficSaturday, GraficSunday)
		SELECT '+7 (903) 796-04-56', 'megapocherk@mail.ru', '9:00 - 20:00', '10:00 - 16:00', 'выходной'
		WHERE NOT EXISTS (SELECT 1 FROM contacts);`)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
	}

	// Шаблоны
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

	// Единый handler для предметов
	mux.HandleFunc("/subject", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Не указан параметр ?name", http.StatusBadRequest)
			return
		}

		pageName := name + ".html"
		if _, ok := pages[pageName]; !ok {
			notFoundHandler(w, r)
			return
		}

		render(w, pageName, nil)
	})

	// О нас (данные из БД)
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		var info ContactInfo
		err := db.QueryRow("SELECT phone, email, GraficBudni, GraficSaturday, GraficSunday FROM contacts LIMIT 1").
			Scan(&info.Phone, &info.Email, &info.GraficBudni, &info.GraficSaturday, &info.GraficSunday)
		if err != nil {
			http.Error(w, "Ошибка загрузки контактов", http.StatusInternalServerError)
			return
		}
		render(w, "about.html", info)
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
