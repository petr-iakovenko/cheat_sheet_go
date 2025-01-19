
Работа с сервером на Go (Golang) включает несколько шагов: от написания простого сервера до добавления маршрутов, обработки запросов и ответов, а также настройки параметров и middleware. Вот пошаговый гайд:

---

### 1. **Создание базового сервера**

В Go стандартная библиотека `net/http` предоставляет всё необходимое для работы с HTTP-серверами.

Пример минимального сервера:

```go
package main

import (
 "fmt"
 "log"
 "net/http"
)

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, World!") // Отправляем простой текст в ответ
 })

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil)) // Запускаем сервер на порту 8080
}
```

**Объяснение:**

- `http.HandleFunc` связывает путь (`/`) с функцией-обработчиком.
- `http.ListenAndServe` запускает сервер на указанном порту (здесь `8080`).
- Функция-обработчик принимает `http.ResponseWriter` (для ответа клиенту) и `*http.Request` (для получения данных запроса).

---

### 2. **Добавление маршрутов**

Каждый маршрут — это отдельный путь, связанный с функцией-обработчиком.

Пример:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "Hello, user!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "This is a simple Go server.")
}

func main() {
 http.HandleFunc("/hello", helloHandler) // Регистрация маршрута "/hello"
 http.HandleFunc("/about", aboutHandler) // Регистрация маршрута "/about"

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

### 3. **Обработка GET и POST запросов**

Вы можете проверять тип запроса (`GET`, `POST`, и др.) через метод `r.Method`.

Пример:

```go
func formHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method == http.MethodPost {
  r.ParseForm() // Парсим тело запроса
  name := r.FormValue("name")
  fmt.Fprintf(w, "Hello, %s!", name)
 } else {
  fmt.Fprintf(w, "Please send a POST request.")
 }
}

func main() {
 http.HandleFunc("/form", formHandler)

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Тестирование:**

1. Отправьте `POST`-запрос с параметром `name`.
2. Можно использовать инструменты вроде `curl`:

   ```bash
   curl -X POST -d "name=John" http://localhost:8080/form
   ```

---

### 4. **Организация кода: маршруты в отдельный пакет**

Для больших проектов маршруты и обработчики лучше разделять на файлы и пакеты.

Пример структуры проекта:

```
/project
  /handlers
    handlers.go
  main.go
```

**`handlers/handlers.go`:**

```go
package handlers

import (
 "fmt"
 "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "Welcome to the Home Page!")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "About us page")
}
```

**`main.go`:**

```go
package main

import (
 "log"
 "net/http"
 "project/handlers"
)

func main() {
 http.HandleFunc("/", handlers.HomeHandler)
 http.HandleFunc("/about", handlers.AboutHandler)

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

### 5. **Работа с JSON**

Go имеет встроенные функции для работы с JSON через пакет `encoding/json`.

**Пример отправки JSON ответа:**

```go
import (
 "encoding/json"
 "net/http"
)

func jsonHandler(w http.ResponseWriter, r *http.Request) {
 data := map[string]string{"message": "Hello, JSON!"}
 w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок
 json.NewEncoder(w).Encode(data) // Кодируем map в JSON и отправляем
}

func main() {
 http.HandleFunc("/json", jsonHandler)

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Пример обработки JSON-запроса:**

```go
func postHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method == http.MethodPost {
  var input map[string]string
  json.NewDecoder(r.Body).Decode(&input) // Декодируем тело запроса
  response := map[string]string{"received": input["data"]}
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)
 } else {
  http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
 }
}
```

---

### 6. **Добавление Middleware**

Middleware — это функции, которые выполняются до вызова обработчика.

**Пример:**

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  log.Printf("Request: %s %s", r.Method, r.URL.Path)
  next(w, r) // Вызываем следующий обработчик
 }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "Hello, World!")
}

func main() {
 http.HandleFunc("/", loggingMiddleware(helloHandler))

 log.Println("Server is running on http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

### 7. **Тестирование**

Тестировать сервер можно с помощью:

- **Postman** или **Insomnia** — графические интерфейсы для работы с запросами.
- **curl** — инструмент командной строки.
- **Тесты в Go** (используя `net/http/httptest`).

**Пример теста:**

```go
import (
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestHelloHandler(t *testing.T) {
 req := httptest.NewRequest("GET", "/", nil)
 rr := httptest.NewRecorder()

 helloHandler(rr, req)

 if rr.Body.String() != "Hello, World!" {
  t.Errorf("Expected 'Hello, World!', got %s", rr.Body.String())
 }
}
```

---

### Ресурсы для изучения

1. [Документация Go: `net/http`](https://pkg.go.dev/net/http)
2. [Туториалы для начинающих](https://gobyexample.com/)
3. Популярные фреймворки: **Gin**, **Echo** (для более сложных серверов).
