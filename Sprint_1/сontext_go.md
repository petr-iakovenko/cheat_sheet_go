
# 🌟 Контекст в Go: полное понятное руководство

Context в Go - это специальный инструмент для управления выполнением горутин, особенно для отмены операций по сигналу или таймауту. Он упрощает передачу сигналов отмены между разными функциями и горутинами.

---

## 🤔 Что такое Context и зачем он нужен?

Представьте, что вы запустили длительную операцию (например, запрос к базе данных), но вдруг пользователь закрыл страницу. Теперь нужно остановить все связанные с этим запросом операции. Как это сделать?

Context - это встроенный в Go механизм для:

- 🚫 Отмены операций (по запросу или таймауту)
- 🔄 Передачи сигнала отмены через цепочку вызовов функций
- 📦 Передачи дополнительной информации между функциями (хотя это менее распространённое применение)

---

## 🔄 Сравнение: канал отмены vs Context

### 📦 Решение с обычным каналом отмены

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

// Функция, которая выполняет длительную операцию
// и может быть отменена через канал
func executeWithCancelChannel(cancelChan <-chan struct{}, operation func() string) (string, error) {
    // Канал для результата операции
    resultChan := make(chan string)
    
    // 🧵 Запускаем операцию в отдельной горутине
    go func() {
        result := operation()
        resultChan <- result
    }()
    
    // 🔀 Ждём либо результат, либо сигнал отмены
    select {
    case result := <-resultChan:
        // ✅ Операция успешно завершилась
        return result, nil
    case <-cancelChan:
        // 🚫 Получен сигнал отмены
        return "", errors.New("операция отменена")
    }
}

func main() {
    // 📋 Создаём длительную операцию
    slowOperation := func() string {
        fmt.Println("🔹 Операция начата")
        time.Sleep(2 * time.Second)
        fmt.Println("✅ Операция завершена")
        return "Результат операции"
    }
    
    // 🛑 Создаём канал для отмены
    cancelChan := make(chan struct{})
    
    // ⏱️ Запускаем таймер, который отменит операцию через 1 секунду
    go func() {
        fmt.Println("⏱️ Таймер запущен на 1 секунду")
        time.Sleep(1 * time.Second)
        fmt.Println("🚫 Таймер сработал, отменяем операцию")
        close(cancelChan)
    }()
    
    // 🚀 Запускаем операцию
    result, err := executeWithCancelChannel(cancelChan, slowOperation)
    if err != nil {
        fmt.Println("❌ Ошибка:", err)
    } else {
        fmt.Println("✅ Результат:", result)
    }
}
```

### 📦 Решение с использованием Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Функция теперь принимает context.Context вместо канала
func executeWithContext(ctx context.Context, operation func() string) (string, error) {
    resultChan := make(chan string)
    
    // 🧵 Операция запускается точно так же
    go func() {
        result := operation()
        resultChan <- result
    }()
    
    // 🔀 Select работает аналогично, но с ctx.Done() вместо канала
    select {
    case result := <-resultChan:
        return result, nil
    case <-ctx.Done():
        // 🚫 Получен сигнал отмены через ctx.Done()
        return "", ctx.Err() // Встроенная ошибка контекста
    }
}

func main() {
    slowOperation := func() string {
        fmt.Println("🔹 Операция начата")
        time.Sleep(2 * time.Second)
        fmt.Println("✅ Операция завершена")
        return "Результат операции"
    }
    
    // 🔑 Создаём контекст с таймаутом вместо ручной отмены
    fmt.Println("⏱️ Создаём контекст с таймаутом 1 секунда")
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    
    // 🧹 Важно всегда вызывать cancel() для освобождения ресурсов
    defer cancel()
    
    // 🚀 Запускаем операцию с контекстом
    result, err := executeWithContext(ctx, slowOperation)
    if err != nil {
        fmt.Println("❌ Ошибка:", err)
    } else {
        fmt.Println("✅ Результат:", result)
    }
}
```

---

## 🔍 Почему Context лучше простого канала отмены?

1. **Стандартизация** 📋
   - Все пакеты Go используют одинаковый интерфейс `context.Context`
   - Не нужно придумывать свои решения для каждого случая

2. **Разные типы отмены в одном интерфейсе** 🔄
   - Ручная отмена: `context.WithCancel()`
   - Таймаут: `context.WithTimeout()` или `context.WithDeadline()`
   - Функция `executeWithContext()` работает со всеми одинаково!

3. **Безопасная многократная отмена** 🛡️
   - Можно вызвать `cancel()` сколько угодно раз без паники
   - С обычным каналом вызов `close(cancelChan)` дважды приведёт к панике

4. **Иерархия контекстов** 🌳
   - Можно создавать дочерние контексты, которые автоматически отменяются при отмене родительских

---

## 📚 Базовые операции с Context

### 1️⃣ Создание пустого контекста

```go
// Создаём корневой (пустой) контекст
ctx := context.Background()
```

`Background()` - это пустой контекст, который никогда не отменяется. Используется как отправная точка для создания других контекстов.

### 2️⃣ Контекст с отменой

```go
// Создаём контекст с возможностью отмены
ctx, cancel := context.WithCancel(context.Background())

// Отменяем контекст (можно в любой момент)
cancel()

// Или с отложенной отменой при выходе из функции
defer cancel()
```

### 3️⃣ Контекст с таймаутом

```go
// Контекст, который автоматически отменится через 5 секунд
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // Всё равно нужно вызвать для освобождения ресурсов
```

### 4️⃣ Контекст с конкретным временем отмены

```go
// Контекст, который отменится в указанное время
deadline := time.Now().Add(1 * time.Hour)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### 5️⃣ Проверка, был ли контекст отменён

```go
// Если контекст отменён, ctx.Err() вернёт причину
if ctx.Err() != nil {
    fmt.Println("Контекст отменён:", ctx.Err())
}

// В горутине можно ждать сигнала отмены через ctx.Done()
select {
case <-ctx.Done():
    fmt.Println("Получен сигнал отмены:", ctx.Err())
    return
case <-time.After(100 * time.Millisecond):
    fmt.Println("Продолжаем работу...")
}
```

---

## 🌟 Практические примеры использования Context

### 📦 Пример 1: Ручная отмена операции

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 🔑 Создаём контекст с возможностью ручной отмены
    ctx, cancel := context.WithCancel(context.Background())
    
    // 📋 Создаём длительную операцию
    processData := func(ctx context.Context) {
        for i := 1; i <= 5; i++ {
            // 🔍 Проверяем, не отменён ли контекст
            select {
            case <-ctx.Done():
                fmt.Println("❌ Обработка прервана на шаге", i)
                return
            default:
                // Продолжаем работу
                fmt.Printf("✅ Шаг %d из 5 выполнен\n", i)
                time.Sleep(500 * time.Millisecond)
            }
        }
        fmt.Println("🎉 Обработка успешно завершена!")
    }
    
    // 🧵 Запускаем обработку в отдельной горутине
    go processData(ctx)
    
    // ⏱️ Ждём 1.7 секунды и отменяем операцию
    time.Sleep(1700 * time.Millisecond)
    fmt.Println("🚫 Отменяем операцию...")
    cancel()
    
    // Даём время горутине обработать отмену и завершиться
    time.Sleep(100 * time.Millisecond)
}
```

Результат:

```
✅ Шаг 1 из 5 выполнен
✅ Шаг 2 из 5 выполнен
✅ Шаг 3 из 5 выполнен
🚫 Отменяем операцию...
❌ Обработка прервана на шаге 4
```

### 📦 Пример 2: Таймаут для HTTP-запроса

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func fetchURL(ctx context.Context, url string) (string, error) {
    // 🔑 Создаём HTTP-запрос с контекстом
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    // 🚀 Выполняем запрос
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    // 📦 Читаем ответ
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    return string(body[:100]) + "...", nil // Возвращаем первые 100 символов
}

func main() {
    // ⏱️ Создаём контекст с таймаутом 2 секунды
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    fmt.Println("🚀 Отправляем запрос с таймаутом 2 секунды...")
    
    // Попробуйте заменить на медленный URL, чтобы увидеть таймаут
    result, err := fetchURL(ctx, "https://example.com")
    
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("⏱️ Превышен таймаут запроса!")
        } else {
            fmt.Println("❌ Ошибка:", err)
        }
    } else {
        fmt.Println("✅ Получен ответ:", result)
    }
}
```

### 📦 Пример 3: Иерархия контекстов

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 👨‍👩‍👧 Создаём иерархию контекстов
    
    // 👴 Родительский контекст с таймаутом 3 секунды
    parentCtx, parentCancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer parentCancel()
    
    // 👧 Дочерний контекст с таймаутом 1 секунда
    childCtx, childCancel := context.WithTimeout(parentCtx, 1*time.Second)
    defer childCancel()
    
    // 📋 Запускаем две горутины с разными контекстами
    go func() {
        <-parentCtx.Done()
        fmt.Println("👴 Родительский контекст отменён:", parentCtx.Err())
    }()
    
    go func() {
        <-childCtx.Done()
        fmt.Println("👧 Дочерний контекст отменён:", childCtx.Err())
    }()
    
    // ⏱️ Ждём...
    time.Sleep(2 * time.Second)
    fmt.Println("🕒 Прошло 2 секунды")
}
```

Результат:

```
👧 Дочерний контекст отменён: context deadline exceeded
🕒 Прошло 2 секунды
```

Через 3 секунды также отменится и родительский контекст.

## 🔍 Ключевые правила при работе с Context

1. **Всегда передавайте контекст первым аргументом** 📋

   ```go
   // Хорошо
   func DoSomething(ctx context.Context, arg string) error {
       // ...
   }
   
   // Не очень
   func DoSomething(arg string, ctx context.Context) error {
       // ...
   }
   ```

2. **Всегда вызывайте cancel() через defer** 🧹

   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel() // Гарантированная очистка ресурсов
   ```

3. **Не храните контекст в структурах** 🚫

   ```go
   // Неправильно
   type Service struct {
       ctx context.Context
   }
   
   // Правильно
   type Service struct {
       // Без хранения контекста
   }
   
   func (s *Service) DoWork(ctx context.Context) {
       // Используем переданный контекст
   }
   ```

4. **Проверяйте отмену контекста в циклах** 🔄

   ```go
   for {
       select {
       case <-ctx.Done():
           return ctx.Err()
       default:
           // Продолжаем работу
       }
       // ...
   }
   ```

## 📊 Когда использовать Context

| Ситуация | Решение | Пример |
|----------|---------|--------|
| HTTP-запросы | `WithTimeout` | Ограничение времени ожидания ответа от API |
| Длительные операции | `WithCancel` | Возможность прервать долгую обработку |
| Цепочка вызовов | Передача контекста | Отмена сразу всей цепочки операций |
| Крайний срок | `WithDeadline` | Операция должна завершиться до определённого времени |

## 🎯 Преимущества использования Context

1. **Стандартный подход** 📋
   - Используется во всей стандартной библиотеке Go
   - Понятен любому Go-разработчику

2. **Предотвращение утечек ресурсов** 🧹
   - Помогает корректно остановить все горутины
   - Вовремя освобождает ресурсы

3. **Гибкость отмены** 🔄
   - Можно отменить вручную
   - Можно задать таймаут
   - Можно задать точное время истечения срока

4. **Простая интеграция в цепочке вызовов** 🔗
   - Просто передаётся из функции в функцию
   - Автоматически распространяет сигнал отмены

## 🌟 Заключение

Context в Go - это мощный инструмент для управления жизненным циклом операций. Он особенно полезен в многопоточных программах, где нужно координировать отмену операций между разными горутинами.

Используя контекст правильно, вы сможете создавать более надёжные и эффективные программы, которые корректно освобождают ресурсы и не оставляют "висящих" горутин.
