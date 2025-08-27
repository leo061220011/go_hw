package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	PrettyDT string
	Error    string
}

// структура для пользователя (админ)
type User struct {
	ID       int
	Username string
	Password string
}

// простая проверка рабочее время
func isWorkingHours(t time.Time) bool {
	weekday := t.Weekday()
	hour := t.Hour()
	minute := t.Minute()

	switch weekday {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		return (hour > 9 && hour < 21) || (hour == 9 && minute >= 0) || (hour == 21 && minute == 0)
	case time.Saturday:
		return (hour > 10 && hour < 16) || (hour == 10 && minute >= 0) || (hour == 16 && minute == 0)
	case time.Sunday:
		return false
	default:
		return false
	}
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

	// Создание таблиц
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
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password TEXT
		);`,
	}
	for _, query := range createTables {
		if _, err := db.Exec(query); err != nil {
			log.Fatal("Ошибка при создании таблицы:", err)
		}
	}

	// Дефолтный контакт
	_, err = db.Exec(`INSERT INTO contacts (phone, email, GraficBudni, GraficSaturday, GraficSunday)
		SELECT '+7 (903) 796-04-56', 'megapocherk@mail.ru', '9:00 - 20:00', '10:00 - 16:00', 'выходной'
		WHERE NOT EXISTS (SELECT 1 FROM contacts);`)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
	}

	// Дефолтный админ (если нет)
	_, err = db.Exec(`INSERT INTO users (username, password)
		SELECT 'admin', 'admin' WHERE NOT EXISTS (SELECT 1 FROM users);`)
	if err != nil {
		log.Fatal("Ошибка при вставке админа:", err)
	}

	// Шаблоны
	pages["index.html"] = mustPage("index.html")
	pages["checkout.html"] = mustPage("checkout.html")
	pages["math.html"] = mustPage("math.html")
	pages["english.html"] = mustPage("english.html")
	pages["about.html"] = mustPage("about.html")
	pages["feedback_admin.html"] = mustPage("feedback_admin.html")
	pages["bookings_admin.html"] = mustPage("bookings_admin.html")
	pages["admin.html"] = mustPage("admin.html")
	pages["login.html"] = mustPage("login.html")

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
			booking := Booking{
				Name:     r.FormValue("name"),
				Email:    r.FormValue("email"),
				Phone:    r.FormValue("phone"),
				Topic:    r.FormValue("topic"),
				DateTime: r.FormValue("event-datetime"),
				Message:  r.FormValue("message"),
			}

			// Проверка обязательных полей
			if booking.Name == "" || booking.Email == "" || booking.Message == "" || booking.DateTime == "" {
				booking.Error = "Заполните обязательные поля"
				render(w, "checkout.html", booking)
				return
			}

			// Парсим дату-время
			t, err := time.Parse("2006-01-02T15:04", booking.DateTime)
			if err != nil {
				booking.Error = "Некорректный формат даты/времени"
				render(w, "checkout.html", booking)
				return
			}

			// Проверка графика работы
			if !isWorkingHours(t) {
				booking.Error = "Выбранное время не входит в график работы"
				render(w, "checkout.html", booking)
				return
			}

			// Проверка на занятость
			var exists int
			err = db.QueryRow("SELECT COUNT(*) FROM bookings WHERE datetime = ?", booking.DateTime).Scan(&exists)
			if err != nil {
				booking.Error = "Ошибка проверки времени"
				render(w, "checkout.html", booking)
				return
			}
			if exists > 0 {
				booking.Error = "Данное время уже забронировано. Выберите другое"
				render(w, "checkout.html", booking)
				return
			}

			// Сохраняем
			_, err = db.Exec(
				"INSERT INTO bookings (name, email, phone, topic, datetime, message) VALUES (?, ?, ?, ?, ?, ?)",
				booking.Name, booking.Email, booking.Phone, booking.Topic, booking.DateTime, booking.Message,
			)
			if err != nil {
				booking.Error = "Ошибка при сохранении"
				render(w, "checkout.html", booking)
				return
			}

			// Успешная заявка
			render(w, "checkout.html", Booking{Success: true})
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

		info.Success = true // указываем, что форма успешно отправлена

		render(w, "about.html", info)
	})

	// // Удаление сообщения
	mux.HandleFunc("/admin/feedback/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id") // читаем из формы, а не из query
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

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "login.html", nil)
			return
		}
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			var user User
			err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
				Scan(&user.ID, &user.Username, &user.Password)
			if err != nil || user.Password != password {
				render(w, "login.html", "Неверный логин или пароль")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_user",
				Value:    user.Username,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
			})

			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	})

	// ADMIN INDEX
	mux.Handle("/admin", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin.html", nil)
	})))

	// ADMIN FEEDBACK
	mux.Handle("/admin/feedback", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	// ADMIN BOOKINGS
	mux.Handle("/admin/bookings", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			t, err := time.Parse("2006-01-02T15:04", b.DateTime)
			if err == nil {
				b.PrettyDT = t.Format("02-01-2006 15:04")
			} else {
				b.PrettyDT = b.DateTime
			}
			bookings = append(bookings, b)
		}

		render(w, "bookings_admin.html", bookings)
	})))

	//STATIC
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

	// 404
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

//HELPERS

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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("static", "404.html")
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, path)
}

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_user")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
