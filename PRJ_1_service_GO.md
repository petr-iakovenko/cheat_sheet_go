напиши структуру сервиса на go:

сервис будет иметь

- через api подключаться к google календарю пользователя
- обновлять токен доступа к календарю пользователя
- сохранять данные о прошедших эвентах и будущих эвентах в мапу (в будущем это будет БД postgres)
- по запросу пользователя выводить информацию по введенным настройкам от пользователя о прошедших и будущий эвентах из календаря

<I’ll answer as the world-famous Go programmer with the Google Open Source Peer Reward for contributions to API service design>

TL;DR: Структура сервиса на Go для интеграции с Google Calendar, обновления токенов, хранения и отображения событий, включая взаимодействие с пользователем, может быть организована следующим образом.

Общая структура сервиса:
````
project/
├── main.go                # Точка входа
├── config/                # Конфигурации
│   └── config.go
├── handlers/              # API-обработчики
│   ├── events.go
│   └── auth.go
├── googleapi/             # Логика взаимодействия с Google API
│   ├── calendar.go
│   └── auth.go
├── models/                # Структуры данных
│   └── event.go
├── storage/               # Локальная база (или интеграция с БД)
│   └── storage.go
└── utils/                 # Утилиты
    └── logger.go
````

Пример реализации:

1. main.go

````go
package main

import (
 "log"
 "net/http"
 "project/handlers"
)

func main() {
 // Регистрируем маршруты
 http.HandleFunc("/auth", handlers.AuthHandler)
 http.HandleFunc("/events", handlers.EventsHandler)

 // Запускаем сервер
 log.Println("Server is running on http://localhost:8080")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  log.Fatal("Error starting server:", err)
 }
}
````

2. config/config.go

````go
package config

import (
 "os"
)

type Config struct {
 ClientID     string
 ClientSecret string
 RedirectURL  string
}

func LoadConfig() Config {
 return Config{
  ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
  ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
  RedirectURL:  "<http://localhost:8080/auth/callback>",
 }
}
````
3. handlers/auth.go

````go
package handlers

import (
 "fmt"
 "net/http"
 "project/googleapi"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
 // Генерация URL для авторизации
 authURL := googleapi.GenerateAuthURL()
 http.Redirect(w, r, authURL, http.StatusFound)
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
 code := r.URL.Query().Get("code")
 if code == "" {
  http.Error(w, "Authorization code is missing", http.StatusBadRequest)
  return
 }

 // Получаем токен
 token, err := googleapi.ExchangeCodeForToken(code)
 if err != nil {
  http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
  return
 }

 fmt.Fprintf(w, "Token saved successfully: %v", token)
}
````

4. handlers/events.go

````go
package handlers

import (
 "encoding/json"
 "net/http"
 "project/storage"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodGet {
  http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
  return
 }

 // Получаем данные событий
 events := storage.GetAllEvents()

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(events)
}
````

5. googleapi/calendar.go

````go
package googleapi

import (
 "context"
 "fmt"

 "google.golang.org/api/calendar/v3"
 "google.golang.org/api/option"
)

func FetchEvents(token string) ([]*calendar.Event, error) {
 ctx := context.Background()
 service, err := calendar.NewService(ctx, option.WithTokenSource(GetTokenSource(token)))
 if err != nil {
  return nil, fmt.Errorf("failed to create calendar service: %w", err)
 }

 events, err := service.Events.List("primary").Do()
 if err != nil {
  return nil, fmt.Errorf("failed to fetch events: %w", err)
 }

 return events.Items, nil
}
````

6. storage/storage.go
```` go
package storage

import (
 "sync"
)

type Event struct {
 ID       string
 Title    string
 Start    string
 End      string
 IsFuture bool
}

var (
 eventMap = make(map[string]Event)
 mu       sync.RWMutex
)

func SaveEvent(event Event) {
 mu.Lock()
 defer mu.Unlock()
 eventMap[event.ID] = event
}

func GetAllEvents() []Event {
 mu.RLock()
 defer mu.RUnlock()

 events := make([]Event, 0, len(eventMap))
 for _, event := range eventMap {
  events = append(events, event)
 }
 return events
}
````
Основные шаги:

 1. Авторизация Google API:
 - Используйте OAuth2 для получения и обновления токенов.
 - Генерация authURL и обработка обратного вызова (/auth и /auth/callback).
 2. Работа с Google Calendar API:
 - Получение событий с помощью библиотеки google.golang.org/api/calendar/v3.
 - Убедитесь, что токен регулярно обновляется.
 3. Хранение данных:
 - Используйте map для временного хранения событий.
 - В будущем замените на PostgreSQL.
 4. API для взаимодействия с пользователем:
 - Обработчик /events возвращает события по запросу.

Если потребуется больше деталей по каждому компоненту, уточните! 😊
