# CRUD

## 👨‍💻 Написание CRUD: Простейшие end-points

```go
package main

import (
 "fmt"
 "net/http"
 "strconv"
)

var counter int

func GetHandler(w http.ResponseWriter,r *http.Request) {
 if r.Method == http.MethodGet {
  fmt.Fprintln(w, "Counter is: ", strconv.Itoa(counter))
 } else {
  fmt.Fprintln(w, "Use only GET")
 }
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method == http.MethodPost {
  counter++
  fmt.Fprintln(w, "Counter was increased by 1 ")
 } else {
  fmt.Fprintln(w, "Use only Post")
 }
}

func main() {
 http.HandleFunc("/get", GetHandler)
 http.HandleFunc("/post", PostHandler)
 http.ListenAndServe("localhost:8080", nil)
}

```

Этот код создает простой HTTP-сервер с двумя эндпоинтами для получения и увеличения счетчика, используя встроенные возможности пакета net/http в Go.

Подробное объяснение работы кода:

1. **Обработчик GET-запросов (GetHandler)**:
   - Функция проверяет, является ли входящий HTTP-запрос методом GET
   - Если да, то выводит текущее значение глобальной переменной `counter`
   - Использует `fmt.Fprintln()` для записи ответа в `http.ResponseWriter`
   - Преобразует число счетчика в строку с помощью `strconv.Itoa()`
   - Если метод не GET, выводит сообщение об ошибке

2. **Обработчик POST-запросов (PostHandler)**:
   - Проверяет, является ли входящий HTTP-запрос методом POST
   - Если да, то увеличивает глобальную переменную `counter` на 1
   - Выводит сообщение о увеличении счетчика
   - Если метод не POST, выводит сообщение об ошибке

3. **Функция main()**:
   - Регистрирует обработчики для двух путей:
     - `/get` - связан с GetHandler (возвращает текущее значение счетчика)
     - `/post` - связан с PostHandler (увеличивает счетчик)
   - Запускает HTTP-сервер на локальном хосте (localhost) на порту 8080
   - Использует `http.HandleFunc()` для маршрутизации
   - `http.ListenAndServe()` начинает прослушивание входящих HTTP-соединений

---

## Модернизация с помощъю Echo

### Пакет Echo

Пакет Echo - это высокопроизводительный и минималистичный веб-фреймворк для Go, который значительно упрощает создание веб-приложений и API.

Основные возможности Echo:

1. **Маршрутизация**:
   - Быстрая и гибкая маршрутизация
   - Поддержка стандартных HTTP-методов (GET, POST, PUT, DELETE и др.)
   - Группировка маршрутов
   - Параметризация путей

   ```go
   e := echo.New()
   e.GET("/users/:id", getUserHandler)
   e.POST("/users", createUserHandler)
   ```

2. **Middleware**:  
(промежуточные обработчики)

   - Встроенные middleware для:
     - Логирования
     - Восстановления после паник
     - Cors
     - Безопасности
     - Сжатия
   - Возможность создания кастомных middleware

   ```go
   e.Use(middleware.Logger())
   e.Use(middleware.Recover())
   ```

3. **Обработка запросов**:
   - Простое связывание данных
   - Валидация входящих данных
   - Удобная работа с контекстом запроса

   ```go
   func(c echo.Context) error {
       user := new(User)
       if err := c.Bind(user); err != nil {
           return err
       }
       return c.JSON(http.StatusOK, user)
   }
   ```

4. **Производительность**:
   - Один из самых быстрых веб-фреймворков в Go
   - Минимальные накладные расходы
   - Эффективная маршрутизация

5. **Расширяемость**:
   - Поддержка шаблонов
   - Работа с WebSocket
   - Легкая интеграция с другими библиотеками

## 👨‍💻 Написание CRUD: добавление работы с JSON

```go
package main

import (
 "net/http"
 "github.com/labstack/echo/v4"
)

type Message struct {
 Text string `json:"text"`
}

type Response struct {
 Status string `json:"status"`
 Message string `json:"message"`
}

var messages []Message 

func GetHandler(c echo.Context) error {
 return c.JSON(http.StatusOK, &messages) // как работает под капотом?
}


func PostHandler(c echo.Context) error {
 var message Message
 if err := c.Bind(&message); err != nil { // как работает под капотом
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Could not add the message",
  })
 }
 messages = append(messages, message)
 return c.JSON(http.StatusOK, Response{
   Status: "Success",
   Message: "The message was successfully added",
 })
}

func main() {
 e := echo.New()
 e.GET("/messages", GetHandler)
 e.POST("/messages", PostHandler)

 e.Start(":8080")
}
```

Подробное объяснение работы кода:

**Обработчик: GetHandler**

```go
func GetHandler(c echo.Context) error {
 return c.JSON(http.StatusOK, &messages) // как работает под капотом?
}
```

***Под капотом `c.JSON` работает следующим образом:***

1. Метод `c.JSON()` принимает HTTP статус-код (200 OK) и данные (&messages) для сериализации
2. Внутри Echo вызывает `json.Marshal()` для преобразования слайса messages в JSON-строку
3. Устанавливает заголовок `"Content-Type: application/json"` в HTTP-ответе
4. Записывает статус-код в ответ (StatusOK = 200)
5. Записывает сериализованный JSON в тело ответа
6. Автоматически управляет потоками данных и закрытием соединения

**Обработчик: PostHandler**

```go
func PostHandler(c echo.Context) error {
 var message Message
 if err := c.Bind(&message); err != nil { // как работает под капотом
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Could not add the message",
  })
 }
 messages = append(messages, message)
 return c.JSON(http.StatusOK, Response{
   Status: "Success",
   Message: "The message was successfully added",
 })
}
```

***Под капотом `c.Bind` работает следующим образом:***

1. Echo анализирует заголовок `"Content-Type"` входящего запроса (обычно `"application/json"`)
2. На основе `Content-Type` выбирает соответствующий десериализатор
3. Для `"application/json"` вызывает `json.Unmarshal()` для преобразования JSON-тела запроса в структуру `Message`
4. Учитывает теги структуры (например, `json:"text"`) для правильного маппинга полей
5. Если формат некорректный или поля не соответствуют структуре, возвращает ошибку
6. Заполняет переданный указатель (&message) данными из запроса
7. Выполняет дополнительную валидацию, если она настроена

**Логика:**

1. Когда приходит GET-запрос на `/messages`, Echo вызывает `GetHandler`, который возвращает весь слайс сообщений в формате JSON
2. Когда приходит POST-запрос на /messages с JSON-телом, Echo вызывает `PostHandler`:
   - JSON десериализуется в структуру `Message`
   - Новое сообщение добавляется в слайс messages
   - Возвращается успешный ответ с подтверждением

---

## 👨‍💻 Написание CRUD: Patch и Delete

```go
package main

import (
 "net/http"
 "strconv"

 "github.com/labstack/echo/v4"
)

type Message struct {
 ID int `json:"id"`
 Text string `json:"text"`
}

type Response struct {
 Status string `json:"status"`
 Message string `json:"message"`
}

var messages = make(map[int]Message) 
var nextID = 1

func GetHandler(c echo.Context) error {
 var messageSlice []Message

 for _, v := range messages {
  messageSlice = append(messageSlice, v)
 }
 return c.JSON(http.StatusOK, &messageSlice)
}


func PostHandler(c echo.Context) error {
 var message Message
 if err := c.Bind(&message); err != nil {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Could not add the message",
  })
 }

 message.ID = nextID
 nextID++

 messages[message.ID] = message
 return c.JSON(http.StatusOK, Response{
  Status: "Success",
  Message: "The message was successfully added",
 })
}

func PatchHandler(c echo.Context) error {
 idParam := c.Param("id")
 id, err := strconv.Atoi(idParam)
 if err != nil {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Bad ID",
  })
 }
 
 var updMessage Message
 if err := c.Bind(&updMessage); err != nil {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Could not add the message",
  })
 }

 if _, exists := messages[id]; !exists {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "The message was not found",
  })
 }

 updMessage.ID = id
 messages[id] = updMessage

 return c.JSON(http.StatusOK, Response{
  Status: "Success",
  Message: "The message was updated",
 })
}

func DeleteHandler(c echo.Context) error {
 idParam := c.Param("id")
 id, err := strconv.Atoi(idParam)
 if err != nil {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "Bad ID",
  })
 }

 if _, exists := messages[id]; !exists {
  return c.JSON(http.StatusBadRequest, Response{
   Status: "Error",
   Message: "The message was not found",
  })
 }

 delete(messages, id)

 return c.JSON(http.StatusOK, Response{
  Status: "Success",
  Message: "The message was deleted",
 })
}

func main() {
 e := echo.New()
 e.GET("/messages", GetHandler)
 e.POST("/messages", PostHandler)
 e.PATCH("/messages/:id", PatchHandler)
 e.DELETE("/messages/:id", DeleteHandler)

 e.Start(":8080")
}
```

Строки `idParam := c.Param("id")` и `id, err := strconv.Atoi(idParam)` работают вместе для извлечения и конвертации параметра ID из URL-пути. Давайте разберем их подробно:

#### `idParam := c.Param("id")`

**Что происходит под капотом:**

1. Метод `c.Param("id")` в Echo фреймворке извлекает значение параметра "id" из URL-пути.
2. Например, если маршрут зарегистрирован как `/messages/:id`, а запрос пришел на `/messages/5`, то `c.Param("id")` вернет строку "5".
3. Echo реализует это следующим образом:
   - При регистрации маршрута `/messages/:id` Echo создает внутреннее дерево маршрутизации
   - При получении запроса Echo сопоставляет URL с деревом маршрутизации
   - Находит совпадающий маршрут и извлекает части пути, которые соответствуют параметрам
   - Сохраняет эти параметры в объекте контекста `c`
   - Метод `Param()` просто возвращает ранее сохраненное значение из контекста

#### `id, err := strconv.Atoi(idParam)`

**Что происходит под капотом:**

1. Функция `strconv.Atoi()` из стандартной библиотеки Go конвертирует строку в целое число (Atoi = ASCII to Integer).
2. Она принимает строку (`idParam`) и возвращает два значения:
   - Полученное целое число (`id`)
   - Ошибку (`err`), если преобразование не удалось
3. Внутри `strconv.Atoi()`:
   - Проверяется каждый символ строки (если строка содержит что-то кроме цифр и возможного знака минуса в начале, возникнет ошибка)
   - Обрабатываются знаки "+" и "-"
   - Последовательно вычисляется результат, умножая текущее значение на 10 и добавляя числовое значение текущего символа
   - Проверяются границы возможных значений int

#### Логика использования в контексте обработчиков

1. В обоих обработчиках (PatchHandler и DeleteHandler) сначала извлекается параметр "id" из URL.
2. Затем строковое значение преобразуется в целое число для использования в логике.
3. В случае ошибки (например, если в URL передано "abc" вместо числа) функция возвращает ошибку 400 Bad Request.
4. В случае успеха полученный ID используется для:
   - В PatchHandler: обновления сообщения по указанному индексу в слайсе messages
   - В DeleteHandler: начала процесса удаления сообщения

---

## 👨‍💻 Написание CRUD: postgres

### 🔍 Что такое GORM?

GORM (Go Object Relational Mapper) — это ORM-библиотека, используемая в этом коде. Она создает абстрактный слой между кодом Go и базой данных PostgreSQL, позволяя работать с записями базы данных, как если бы они были обычными структурами Go.

```go
import (
 "log"
 "net/http"
 "strconv"
 "gorm.io/driver/postgres"
 "github.com/labstack/echo/v4"
 "gorm.io/gorm"
)

type Message struct {
 Id int `json:"id"`
 Text string `json:"text"`
}

type Response struct {
 Status string `json:"status"`
 Message string `json:"message"`
}

var db *gorm.DB

func initDB() {
   dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5003 sslmode=disable"
   var err error 
   db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
   if err != nil {
      log.Fatalf("Could not connected to database: %v", err)
   }

 db.AutoMigrate(&Message{})
}

func GetHandler(c echo.Context) error {
   var messages []Message
 
   if err := db.Find(&messages).Error; err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Could not find the messages",
      })
   }

 return c.JSON(http.StatusOK, &messages)
}

func PostHandler(c echo.Context) error {
   var message Message
   if err := c.Bind(&message); err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Could not find the message",
      })
   }

   if err := db.Create(&message).Error; err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Could not find the message",
      })
   }  

   return c.JSON(http.StatusOK, Response{
      Status: "Success",
      Message: "The message was added",
   })
}


func PatchHandler(c echo.Context) error {
   idParam := c.Param("id")
   id, err := strconv.Atoi(idParam)
   if err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Bad ID",
      })
   }
   // id check existing
   var existingMessage Message
   if err := db.First(&existingMessage, id).Error; err != nil {
      return c.JSON(http.StatusNotFound, Response{
         Status: "Error",
         Message: "Message with this ID not found",
      })
   }

   var updMessage Message
   if err := c.Bind(&updMessage); err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Invalid input",
      })
   }

   if err := db.Model(&Message{}).Where("id = ?", id).Update("text", updMessage.Text).Error; err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Could not update the message",
      })
   }

   return c.JSON(http.StatusOK, Response{
      Status: "Success",
      Message: "The message was updated",
   })
}

func DeleteHandler(c echo.Context) error {
   idParam := c.Param("id")
   id, err := strconv.Atoi(idParam)
   if err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Bad ID",
      })
   } 

   // id check existing
   var existingMessage Message
   if err := db.First(&existingMessage, id).Error; err != nil {
      return c.JSON(http.StatusNotFound, Response{
         Status: "Error",
         Message: "Message with this ID not found",
   })
 }

   if err := db.Delete(&Message{}, id).Error; err != nil {
      return c.JSON(http.StatusBadRequest, Response{
         Status: "Error",
         Message: "Could not delete the message",
      })
   }

   return c.JSON(http.StatusOK, Response{
      Status: "Success",
      Message: "The message was deleted",
   })
}

func main() {
   initDB()

   e := echo.New()
   e.GET("/messages", GetHandler)
   e.POST("/messages", PostHandler)
   e.PATCH("/messages/:id", PatchHandler)
   e.DELETE("/messages/:id", DeleteHandler)

   e.Start(":8080")
}

```

### 🏗️ Настройка подключения к базе данных

```go
func initDB() {
    dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5003 sslmode=disable"
    var err error 
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Could not connected to database: %v", err)
    }

    db.AutoMigrate(&Message{})
}
```

Эта функция:

1. ⚙️ Создает строку DSN (Data Source Name) с параметрами подключения
2. 🔌 Устанавливает соединение с PostgreSQL с помощью `gorm.Open()`
3. 🚀 Запускает `AutoMigrate`, который создает или обновляет таблицу базы данных, чтобы она соответствовала структуре `Message`

### 💡 CRUD операции

#### 📖 Чтение (GET)

```go
func GetHandler(c echo.Context) error {
    var messages []Message
    
    if err := db.Find(&messages).Error; err != nil {
        // Обработка ошибок
    }

    return c.JSON(http.StatusOK, &messages)
}
```

Под капотом:

1. 🔎 `db.Find(&messages)` преобразуется в `SELECT * FROM messages`
2. 🔄 GORM выполняет запрос, извлекает все записи и заполняет срез `messages`
3. ✅ Возвращает результаты в формате JSON

#### ✏️ Создание (POST)

```go
func PostHandler(c echo.Context) error {
    var message Message
    if err := c.Bind(&message); err != nil {
        // Обработка ошибок
    }

    if err := db.Create(&message).Error; err != nil {
        // Обработка ошибок
    }

    return c.JSON(http.StatusOK, Response{
        Status: "Success",
        Message: "The message was added",
    })
}
```

Под капотом:

1. 📝 `c.Bind(&message)` разбирает JSON из тела запроса в структуру `message`
2. 🆕 `db.Create(&message)` генерирует и выполняет `INSERT INTO messages (text) VALUES (?)` с текстом сообщения
3. 🆔 GORM получает новый созданный ID и обновляет поле Id структуры

#### 🔄 Обновление (PATCH)

```go
func PatchHandler(c echo.Context) error {
    // Код валидации ID и привязки объекта...

    if err := db.Model(&Message{}).Where("id = ?", id).Update("text", updMessage.Text).Error; err != nil {
        // Обработка ошибок
    }
    
    return c.JSON(http.StatusOK, Response{
        Status: "Success",
        Message: "The message was updated",
    })
}
```

Под капотом:

1. 🔍 Сначала проверяет существование сообщения с помощью `db.First(&existingMessage, id)`
2. 🛠️ `db.Model(&Message{}).Where("id = ?", id).Update("text", updMessage.Text)` преобразуется в:

   ```sql
   UPDATE messages SET text = ? WHERE id = ?
   ```

3. ✅ GORM обрабатывает привязку параметров и предотвращение SQL-инъекций

#### 🗑️ Удаление (DELETE)

```go
func DeleteHandler(c echo.Context) error {
    // Код валидации ID...
    
    if err := db.Delete(&Message{}, id).Error; err != nil {
        // Обработка ошибок
    }

    return c.JSON(http.StatusOK, Response{
        Status: "Success",
        Message: "The message was deleted",
    })
}
```

Под капотом:

1. 🔍 Сначала проверяет существование сообщения
2. 🗑️ `db.Delete(&Message{}, id)` генерирует и выполняет:

   ```sql
   DELETE FROM messages WHERE id = ?
   ```

## 🚀 Расширенные возможности GORM, которые не показаны, но доступны

### 🔄 Транзакции

```go
tx := db.Begin()
// Операции внутри транзакции
if err := tx.Create(&someObject).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

### 🔍 Расширенные запросы

```go
// Условия WHERE
db.Where("text LIKE ?", "%search%").Find(&messages)

// Сортировка
db.Order("created_at desc").Find(&messages)

// Лимит и смещение
db.Limit(10).Offset(5).Find(&messages)

// Объединение связанных таблиц (с ассоциациями)
db.Joins("JOIN users ON users.id = messages.user_id").Find(&messages)
```

### 🎯 Работа с определенными полями

```go
// Выбрать определенные столбцы
db.Select("id", "text").Find(&messages)

// Исключить определенные столбцы
db.Omit("created_at").Find(&messages)
```

### 🧠 Использование чистого SQL при необходимости

```go
// Выполнить чистый SQL
db.Raw("SELECT * FROM messages WHERE text LIKE ?", "%search%").Scan(&messages)

// Именованные параметры
db.Where("text = @text", map[string]interface{}{"text": "Hello"}).Find(&messages)
```

## 🔑 Ключевые концепции GORM, которые стоит запомнить

1. ⚡ **Цепочка методов**: GORM позволяет объединять методы в цепочку (например, `db.Where(...).Order(...).Find(...)`)
2. 🛡️ **Безопасность**: GORM автоматически обеспечивает защиту от SQL-инъекций через параметризованные запросы
3. 🚦 **Обработка ошибок**: Всегда проверяйте `.Error` после каждой операции
4. 🧩 **Теги структур**: Управляйте именами столбцов, отношениями и валидациями с помощью тегов структур
5. 👓 **Ленивая загрузка**: По умолчанию запросы выполняются, когда это необходимо
6. 🎛️ **Колбэки**: GORM имеет хуки для настройки поведения (BeforeSave, AfterCreate и т.д.)
