# Как работает сервер в Go и что такое `net/http`

## Основы HTTP-сервера в Go

### Что такое сервер?

Сервер — это программа, которая принимает запросы от клиентов (например, браузеров) и отвечает на них. В Go для создания сервера используется стандартный пакет `net/http`, который предоставляет всё необходимое для работы с HTTP.

Мы разберём всё простыми словами и на примерах.

---

### Как сервер принимает запросы?

1. **Сервер слушает порт** (например, `:8080`), ожидая запросы.
2. **Когда приходит запрос**, сервер передаёт его обработчику (функции).
3. **Обработчик отвечает клиенту**, возвращая нужные данные (например, HTML-страницу или JSON).

---

### Минимальный HTTP-сервер в Go

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Создаём обработчик для корневого пути "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!") // Отправляем текст ответа клиенту
    })

    // Запускаем сервер на порту 8080
    fmt.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

1. **`http.HandleFunc("/", handler)`**:  
   - Регистрирует обработчик для пути `/`.
   - Когда кто-то заходит на ваш сервер (например, `http://localhost:8080/`), вызывается указанная функция.

2. **`http.ListenAndServe(":8080", nil)`**:  
   - Запускает сервер и слушает запросы на порту `8080`.  
   - Второй аргумент (`nil`) означает, что используется стандартный маршрутизатор (регистрирует обработчики через `http.HandleFunc`).

3. **Функция-обработчик:**
   - Функция принимает два аргумента:
     - `w http.ResponseWriter`: Для отправки ответа клиенту.
     - `r *http.Request`: Содержит информацию о запросе (например, путь, метод, заголовки).

---

### Как это работает?

1. Запустите сервер.
2. Перейдите в браузере на `http://localhost:8080/`.
3. Увидите текст: `Hello, World!`.

---

## Как обрабатывать разные пути

Сервер может отвечать по разным путям. Например:

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Обработчик для пути "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to the homepage!")
    })

    // Обработчик для пути "/about"
    http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "This is the about page.")
    })

    // Запускаем сервер на порту 8080
    fmt.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

1. Перейдите на:
   - `http://localhost:8080/` — откроется главная страница.
   - `http://localhost:8080/about` — откроется страница "About".

---

## Как обрабатывать разные методы (GET, POST)

HTTP-запросы могут быть разных типов: **GET** (получить данные), **POST** (отправить данные), и другие. Вы можете проверить тип запроса и обработать его:

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            // Обработка GET-запроса
            fmt.Fprintln(w, "This is a GET request.")
        } else if r.Method == http.MethodPost {
            // Обработка POST-запроса
            r.ParseForm() // Парсим данные из тела запроса
            name := r.FormValue("name")
            fmt.Fprintf(w, "Hello, %s!", name)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

1. Перейдите на `http://localhost:8080/form` (GET-запрос).  
   Ответ: `This is a GET request.`

2. Отправьте POST-запрос с параметром `name` (например, через Postman):  
   Ответ: `Hello, [имя]!`

---

## Отправка JSON-ответов

Часто серверы отправляют данные в формате JSON. Это можно сделать с помощью стандартного пакета `encoding/json`.

```go
package main

import (
    "encoding/json"
    "net/http"
)

func main() {
    http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
        user := map[string]string{
            "name":  "John Doe",
            "email": "john@example.com",
        }

        w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок
        json.NewEncoder(w).Encode(user)                   // Отправляем JSON
    })

    fmt.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

1. Перейдите на `http://localhost:8080/user`.  
   Ответ:

   ```json
   {
       "name": "John Doe",
       "email": "john@example.com"
   }
   ```

---

## Как читать параметры запроса (query parameters)

Параметры запроса передаются в URL после `?`. Например: `http://localhost:8080/search?query=golang`.

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query().Get("query") // Получаем значение параметра "query"
        if query == "" {
            fmt.Fprintln(w, "No query provided")
        } else {
            fmt.Fprintf(w, "You searched for: %s", query)
        }
    })

    fmt.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

1. Перейдите на `http://localhost:8080/search?query=golang`.  
   Ответ: `You searched for: golang`.

---

## Подведение итогов

### Что вы узнали

1. **`http.HandleFunc`**:
   - Регистрация обработчиков для разных путей.
2. **Обработчики запросов**:
   - Функции, принимающие `http.ResponseWriter` (для ответа) и `*http.Request` (для запроса).
3. **Типы запросов (GET, POST)**:
   - Как проверять и обрабатывать разные методы.
4. **JSON-ответы**:
   - Отправка данных в формате JSON.
5. **Параметры запроса**:
   - Как читать данные из URL.

---

### Минимальный шаблон для вашего первого сервера

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to my server!")
    })

    fmt.Println("Server is running on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

___
___
___
___
___
___
___
___
___
___

# Цель — построить **структуру Go-проекта для простого REST API**.

---

## **Что мы строим?**

Мы делаем **API для управления пользователями**. Наше приложение должно:
1. Уметь добавлять пользователей.
2. Возвращать список всех пользователей.
3. Удалять пользователей.

**API (REST)** — это способ, с помощью которого другие программы (например, мобильные приложения) могут "общаться" с нашим приложением через HTTP-запросы.

---

## **Как будет выглядеть запросы к нашему API?**

1. **Добавить пользователя**  
   Запрос:
   ```
   POST /users
   {
       "name": "John Doe"
   }
   ```
   Ответ:
   ```
   {
       "message": "user created"
   }
   ```

2. **Получить список пользователей**  
   Запрос:
   ```
   GET /users
   ```
   Ответ:
   ```
   [
       {"id": 1, "name": "John Doe"},
       {"id": 2, "name": "Jane Smith"}
   ]
   ```

3. **Удалить пользователя**  
   Запрос:
   ```
   DELETE /users/1
   ```
   Ответ:
   ```
   {
       "message": "user deleted"
   }
   ```

---

## **Как мы будем проектировать приложение?**

1. **Структура кода**:  
   Мы разобьём код на модули, чтобы легко добавлять новые функции и поддерживать проект.

2. **Модули:**
   - **Точка входа**: Главная программа запускает сервер.
   - **Сервер**: Управляет маршрутизацией (куда отправлять запросы).
   - **Обработчики запросов**: Обрабатывают HTTP-запросы.
   - **Бизнес-логика**: Управляет данными и действиями.
   - **Хранилище**: Сохраняет данные (пока только в памяти).

3. **Что это нам даёт?**
   - **Понятная структура**: Каждый модуль выполняет только свою задачу.
   - **Масштабируемость**: Если нужно добавить новую функцию, это легко сделать.
   - **Поддерживаемость**: Код легче читать и менять.

---

### **Шаг 1: Точка входа (`main.go`)**

Точка входа запускает сервер. Это наша стартовая точка.

**Файл `cmd/app/main.go`:**
```go
package main

import (
    "log"
    "project/internal/server"
)

func main() {
    // Создаём сервер
    srv := server.NewServer()

    // Логируем и запускаем сервер
    log.Println("Starting server on :8080")
    log.Fatal(srv.ListenAndServe(":8080")) // Запуск сервера на порту 8080
}
```

**Почему так?**
- `main.go` нужен для запуска приложения. Вся сложная логика (сервер, маршруты) скрыта в других модулях.

---

### **Шаг 2: Сервер и маршрутизация**

Сервер управляет HTTP-запросами и отправляет их в нужный обработчик.

#### **Файл `internal/server/server.go`:**
```go
package server

import (
    "net/http"
    "project/internal/server/router"
)

// Server — структура для хранения HTTP-сервера
type Server struct {
    httpServer *http.Server
}

// NewServer создаёт новый сервер
func NewServer() *Server {
    mux := router.NewRouter() // Регистрируем маршруты
    return &Server{
        httpServer: &http.Server{
            Addr:    ":8080",
            Handler: mux,
        },
    }
}

// ListenAndServe запускает сервер
func (s *Server) ListenAndServe(addr string) error {
    s.httpServer.Addr = addr
    return s.httpServer.ListenAndServe()
}
```

**Файл `internal/server/router.go`:**
```go
package router

import (
    "net/http"
    "project/internal/user"
)

// NewRouter создаёт маршруты
func NewRouter() *http.ServeMux {
    mux := http.NewServeMux()

    // Регистрируем обработчики
    mux.HandleFunc("/users", user.HandleUsers) // GET, POST
    mux.HandleFunc("/users/", user.HandleUser) // DELETE

    return mux
}
```

**Почему так?**
- Сервер принимает запросы и перенаправляет их в обработчики, регистрируемые в `router`.
- Мы отделяем маршруты от основной логики, чтобы их легко менять.

---

### **Шаг 3: Обработчики запросов**

Обработчики отвечают на HTTP-запросы: **создают пользователя**, **возвращают список** или **удаляют пользователя**.

#### **Файл `internal/user/handler.go`:**
```go
package user

import (
    "encoding/json"
    "net/http"
    "project/pkg/utils"
    "strconv"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        users := userService.GetAll() // Получаем всех пользователей
        utils.JSONResponse(w, users, http.StatusOK)
    case http.MethodPost:
        var u User
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
            utils.JSONResponse(w, map[string]string{"error": "invalid input"}, http.StatusBadRequest)
            return
        }
        userService.Create(u) // Создаём пользователя
        utils.JSONResponse(w, map[string]string{"message": "user created"}, http.StatusCreated)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func HandleUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Path[len("/users/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.JSONResponse(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
        return
    }

    if err := userService.Delete(id); err != nil {
        utils.JSONResponse(w, map[string]string{"error": "user not found"}, http.StatusNotFound)
        return
    }

    utils.JSONResponse(w, map[string]string{"message": "user deleted"}, http.StatusOK)
}
```

**Почему так?**
- Обработчики принимают запросы и вызывают бизнес-логику.
- Мы разбиваем логику на несколько частей, чтобы каждая выполняла одну задачу.

---

### **Шаг 4: Бизнес-логика**

Бизнес-логика определяет, **что делать с данными**.

#### **Файл `internal/user/service.go`:**
```go
package user

import "errors"

type UserService struct {
    repo *UserRepository
}

var userService = &UserService{repo: NewUserRepository()}

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func (s *UserService) Create(u User) {
    s.repo.Create(u)
}

func (s *UserService) GetAll() []User {
    return s.repo.GetAll()
}

func (s *UserService) Delete(id int) error {
    return s.repo.Delete(id)
}
```

**Почему так?**
- Бизнес-логика управляет данными. Она "знает", что делать, но не знает, как сохранять данные — это задача репозитория.

---

### **Шаг 5: Хранилище**

Хранилище отвечает за **сохранение и удаление данных**.

#### **Файл `internal/user/repository.go`:**
```go
package user

import "errors"

type UserRepository struct {
    users  []User
    nextID int
}

func NewUserRepository() *UserRepository {
    return &UserRepository{nextID: 1}
}

func (r *UserRepository) Create(u User) {
    u.ID = r.nextID
    r.nextID++
    r.users = append(r.users, u)
}

func (r *UserRepository) GetAll() []User {
    return r.users
}

func (r *UserRepository) Delete(id int) error {
    for i, user := range r.users {
        if user.ID == id {
            r.users = append(r.users[:i], r.users[i+1:]...)
            return nil
        }
    }
    return errors.New("user not found")
}
```

**Почему так?**
- Репозиторий отвечает за доступ к данным, а не за то, как их использовать. Это делает код модульным и легко расширяемым.

---

### **Шаг 6: Утилиты**

#### **Файл `pkg/utils/response.go`:**
```go
package utils

import (
    "encoding/json"
    "net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

**Почему так?**
- Мы выносим общие функции в отдельный модуль, чтобы не дублировать код.

---

### Итоговый чек-лист

1. **Точка входа (`main.go`):**
   - Отвечает только за запуск приложения.
2. **Сервер:**
   - Управляет маршрутизацией.
3. **Обработчики:**
   - Обрабатывают запросы (HTTP) и вызывают бизнес-логику.
4. **Бизнес-логика:**
   - Управляет данными (что делать).
5. **Репозиторий:**
   - Отвечает за хранение данных (где хранить).
6. **Утилиты:**
   - Общие функции для всех модулей.

Этот подход помогает создавать **модульный, читаемый и масштабируемый код**. Если что-то осталось непонятным — дай знать!

