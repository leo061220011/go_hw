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
	Success        bool
}

// структура для сообщений обратной связи
type Feedback struct {
	ID      int
	Name    string
	Email   string
	Phone   string
	Topic   string
	Message string
}

// структура для бронирований
type Booking struct {
	ID       int
	Name     string
	Email    string
	Phone    string
	Topic    string
	DateTime string
	Message  string
	Success  bool
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

	// Создать таблицы, если ещё нет
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone TEXT,
			email TEXT,
			GraficBudni TEXT,
			GraficSaturday TEXT,
			GraficSunday TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS feedback (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			phone TEXT,
			topic TEXT,
			message TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS bookings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    phone TEXT,
    topic TEXT,
    datetime TEXT,
    message TEXT
);`,
	}
	for _, query := range createTables {
		if _, err := db.Exec(query); err != nil {
			log.Fatal("Ошибка при создании таблицы:", err)
		}
	}

	// Если таблица contacts пустая – вставим данные
	_, err = db.Exec(`INSERT INTO contacts (phone, email, GraficBudni, GraficSaturday, GraficSunday)
		SELECT '+7 (903) 796-04-56', 'megapocherk@mail.ru', '9:00 - 20:00', '10:00 - 16:00', 'выходной'
		WHERE NOT EXISTS (SELECT 1 FROM contacts);`)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
	}

	// Шаблоны
	pages["index.html"] = mustPage("index.html")
	pages["checkout.html"] = mustPage("checkout.html")
	pages["math.html"] = mustPage("math.html")
	pages["english.html"] = mustPage("english.html")
	pages["about.html"] = mustPage("about.html")
	pages["feedback_admin.html"] = mustPage("feedback_admin.html")
	pages["bookings_admin.html"] = mustPage("bookings_admin.html")

	mux := http.NewServeMux()

	// Главная
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundHandler(w, r)
			return
		}
		render(w, "index.html", nil)
	})

	// checkout
	mux.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "checkout.html", nil)
			return
		}

		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			email := r.FormValue("email")
			phone := r.FormValue("phone")
			topic := r.FormValue("topic")
			datetime := r.FormValue("event-datetime")
			message := r.FormValue("message")

			if name == "" || email == "" || message == "" {
				http.Error(w, "Заполните обязательные поля", http.StatusBadRequest)
				return
			}

			_, err := db.Exec(
				"INSERT INTO bookings (name, email, phone, topic, datetime, message) VALUES (?, ?, ?, ?, ?, ?)",
				name, email, phone, topic, datetime, message,
			)
			if err != nil {
				http.Error(w, "Ошибка при сохранении", http.StatusInternalServerError)
				return
			}

			// Возвращаем ту же форму, но с Success
			booking := Booking{Success: true}
			render(w, "checkout.html", booking)
			return
		}

		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
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

	// // Обработчик формы обратной связи
	// mux.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Redirect(w, r, "/about", http.StatusSeeOther)
	// 		return
	// 	}

	// 	if err := r.ParseForm(); err != nil {
	// 		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
	// 		return
	// 	}

	// 	name := r.FormValue("name")
	// 	email := r.FormValue("email")
	// 	phone := r.FormValue("phone")
	// 	topic := r.FormValue("topic")
	// 	message := r.FormValue("message")

	// 	_, err := db.Exec(`INSERT INTO feedback (name, email, phone, topic, message) VALUES (?, ?, ?, ?, ?)`,
	// 		name, email, phone, topic, message)
	// 	if err != nil {
	// 		http.Error(w, "Ошибка записи в базу", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	log.Printf("Новое сообщение от %s <%s>: %s", name, email, message)

	// 	http.Redirect(w, r, "/about", http.StatusSeeOther)
	// })
	mux.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		topic := r.FormValue("topic")
		message := r.FormValue("message")

		if name == "" || email == "" || message == "" {
			http.Error(w, "Заполните обязательные поля", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			"INSERT INTO feedback (name, email, phone, topic, message) VALUES (?, ?, ?, ?, ?)",
			name, email, phone, topic, message,
		)
		if err != nil {
			http.Error(w, "Ошибка при сохранении", http.StatusInternalServerError)
			return
		}

		// Загружаем контакты снова, чтобы отрисовать about.html
		var info ContactInfo
		err = db.QueryRow("SELECT phone, email, GraficBudni, GraficSaturday, GraficSunday FROM contacts LIMIT 1").
			Scan(&info.Phone, &info.Email, &info.GraficBudni, &info.GraficSaturday, &info.GraficSunday)
		if err != nil {
			http.Error(w, "Ошибка загрузки контактов", http.StatusInternalServerError)
			return
		}

		info.Success = true // <-- указываем, что форма успешно отправлена

		render(w, "about.html", info)
	})

	// Админка: список сообщений
	mux.HandleFunc("/admin/feedback", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, phone, topic, message FROM feedback ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Ошибка загрузки сообщений", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []Feedback
		for rows.Next() {
			var f Feedback
			if err := rows.Scan(&f.ID, &f.Name, &f.Email, &f.Phone, &f.Topic, &f.Message); err != nil {
				http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
				return
			}
			messages = append(messages, f)
		}

		render(w, "feedback_admin.html", messages)
	})
	// Админка: список заявок
	mux.HandleFunc("/admin/bookings", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, phone, topic, datetime, message FROM bookings ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Ошибка загрузки заявок", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookings []Booking
		for rows.Next() {
			var b Booking
			if err := rows.Scan(&b.ID, &b.Name, &b.Email, &b.Phone, &b.Topic, &b.DateTime, &b.Message); err != nil {
				http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
				return
			}
			bookings = append(bookings, b)
		}

		render(w, "bookings_admin.html", bookings)
	})

	// Удаление сообщения
	mux.HandleFunc("/admin/feedback/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id") // <-- читаем из формы, а не из query
		if id == "" {
			http.Error(w, "Не указан id", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM feedback WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/feedback", http.StatusSeeOther)
	})
	// Удаление заявки
	mux.HandleFunc("/admin/bookings/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Не указан id", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM bookings WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
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
