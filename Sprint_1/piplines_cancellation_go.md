# Go Concurrency Patterns: Pipelines and cancellation

Ссылка на документацию паттернов - [ссылка](https://go.dev/blog/pipelines)

Ссылка на паттерн **fan-out/fan-in** - [ссылка](https://dzen.ru/a/ZUkC52CW8UQlC7vp)

---

Go каналы и горутины создают эффективные конвейеры данных (pipelines) через последовательные этапы обработки, соединенные каналами. Ключевые паттерны включают: fan-out (распараллеливание работы), fan-in (объединение результатов) и корректная обработка отмены.

## 🚀 Что такое конвейеры (pipelines) в Go?

Конвейер в Go - это просто набор этапов (stages), соединенных каналами. Каждый этап - это группа горутин, выполняющих одну и ту же функцию:

1. Получают данные из входного канала
2. Обрабатывают эти данные
3. Отправляют результат в выходной канал

Давайте разберем на простом примере:

```go
// Простой конвейер для обработки чисел
package main

import (
    "fmt"
    "time"
)

// 📥 Этап 1: Генерация чисел
func generateNumbers() <-chan int {
    out := make(chan int)
    go func() {
        // Используем defer для закрытия канала при завершении
        defer close(out)
        for i := 1; i <= 5; i++ {
            fmt.Printf("⏩ Отправляю число: %d\n", i)
            out <- i
            time.Sleep(100 * time.Millisecond) // Имитация работы
        }
    }()
    return out
}

// 🔄 Этап 2: Умножение на 2
func multiplyByTwo(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for number := range in {
            result := number * 2
            fmt.Printf("✖️ Умножаю %d на 2: %d\n", number, result)
            out <- result
            time.Sleep(200 * time.Millisecond) // Имитация работы
        }
    }()
    return out
}

// 📊 Этап 3: Прибавление 10
func addTen(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for number := range in {
            result := number + 10
            fmt.Printf("➕ Прибавляю к %d число 10: %d\n", number, result)
            out <- result
            time.Sleep(150 * time.Millisecond) // Имитация работы
        }
    }()
    return out
}

func main() {
    // Настраиваем конвейер
    numbers := generateNumbers()
    doubled := multiplyByTwo(numbers)
    results := addTen(doubled)

    // Потребляем результаты
    for r := range results {
        fmt.Printf("📝 Получен результат: %d\n", r)
    }

    fmt.Println("🏁 Программа завершена")
}
```

## 🔄 Что здесь происходит?

1. **Генерация чисел**: Первый этап создаёт числа от 1 до 5 и отправляет их в канал.
2. **Умножение на 2**: Второй этап получает числа, умножает их на 2 и отправляет дальше.
3. **Прибавление 10**: Третий этап получает удвоенные числа и прибавляет к ним 10.
4. **Потребление результатов**: В `main()` мы читаем и выводим финальные результаты.

## 🌟 Ключевые принципы конвейеров

1. **Закрытие каналов**: Каждый этап закрывает свой выходной канал, когда заканчивает работу (`defer close(out)`).
2. **Получение данных до закрытия**: Каждый этап продолжает получать данные, пока входной канал не закроется (используя `for number := range in`).
3. **Горутины для параллельности**: Каждый этап запускается в своей горутине, что позволяет всем этапам работать одновременно.

## 🔀 Fan-out и Fan-in

Теперь давайте разберем две важные концепции: fan-out и fan-in.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Генерация чисел - как в предыдущем примере
func generateNumbers() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 1; i <= 10; i++ {
            fmt.Printf("⏩ Отправляю число: %d\n", i)
            out <- i
            time.Sleep(100 * time.Millisecond)
        }
    }()
    return out
}

// Возведение в квадрат
func squareNumber(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for number := range in {
            time.Sleep(500 * time.Millisecond) // Имитация тяжелой работы
            result := number * number
            fmt.Printf("🔢 %d в квадрате: %d\n", number, result)
            out <- result
        }
    }()
    return out
}

// 🌟 Fan-in (объединение каналов)
func mergeChannels(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    merged := make(chan int)
    
    // Функция для копирования данных из входного канала в выходной
    copyValues := func(ch <-chan int) {
        defer wg.Done()
        for val := range ch {
            merged <- val
        }
    }
    
    // Запускаем горутину-копировщик для каждого входного канала
    wg.Add(len(channels))
    for _, ch := range channels {
        go copyValues(ch)
    }
    
    // Горутина для закрытия выходного канала после завершения всех копировщиков
    go func() {
        wg.Wait()
        close(merged)
    }()
    
    return merged
}

func main() {
    numbers := generateNumbers()
    
    // 🌟 Fan-out: запускаем 3 "обработчика" с одним входным каналом
    worker1 := squareNumber(numbers)
    worker2 := squareNumber(numbers)
    worker3 := squareNumber(numbers)
    
    // 🌟 Fan-in: объединяем результаты всех обработчиков в один канал
    results := mergeChannels(worker1, worker2, worker3)
    
    // Потребляем объединенные результаты
    for result := range results {
        fmt.Printf("📝 Получен результат: %d\n", result)
    }
    
    fmt.Println("🏁 Программа завершена")
}
```

---

## 🔍 Что такое Fan-out и Fan-in?

1. **Fan-out (распределение)** - это когда несколько горутин читают из одного канала. Это позволяет параллельно обрабатывать данные из одного источника. В примере: `worker1`, `worker2`, `worker3` - все читают из одного канала `numbers`.

2. **Fan-in (объединение)** - это когда один канал получает данные из нескольких источников. Функция `mergeChannels` берет несколько входных каналов и объединяет их в один выходной.

## ⚠️ Обработка отмены

Иногда нам нужно прервать конвейер досрочно. Представим, что пользователь нажал Ctrl+C или возникла ошибка:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Генерация чисел с поддержкой отмены
func generateNumbers(cancel <-chan struct{}) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 1; ; i++ {
            select {
            case <-cancel:
                fmt.Println("❌ Генерация чисел отменена")
                return
            case out <- i:
                fmt.Printf("⏩ Отправлено число: %d\n", i)
                time.Sleep(200 * time.Millisecond)
            }
        }
    }()
    return out
}

// Возведение в квадрат с поддержкой отмены
func squareNumber(cancel <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for number := range in {
            select {
            case <-cancel:
                fmt.Println("❌ Возведение в квадрат отменено")
                return
            case out <- number * number:
                fmt.Printf("🔢 %d в квадрате: %d\n", number, number*number)
                time.Sleep(300 * time.Millisecond)
            }
        }
    }()
    return out
}

func main() {
    // Канал для сигнала отмены
    cancel := make(chan struct{})
    
    numbers := generateNumbers(cancel)
    squares := squareNumber(cancel, numbers)
    
    // Читаем только первые 5 значений
    for i := 0; i < 5; i++ {
        fmt.Printf("📝 Получен результат: %d\n", <-squares)
    }
    
    // Отменяем конвейер
    fmt.Println("🛑 Отправляем сигнал отмены")
    close(cancel)
    
    // Даем горутинам время завершиться корректно
    time.Sleep(1 * time.Second)
    fmt.Println("🏁 Программа завершена")
}
```

## 🔑 Ключевые моменты отмены

1. **Канал отмены**: Мы создаем специальный канал `cancel`, который используется для сигнализации о необходимости завершения работы.

2. **Использование select**: В каждой горутине мы используем `select` для проверки как основного канала, так и канала отмены.

3. **Закрытие для сигнала**: Мы закрываем канал отмены (`close(cancel)`), а не отправляем в него значение. Когда канал закрыт, все операции получения (`<-cancel`) сразу же завершаются, что позволяет сигнализировать всем горутинам одновременно.

## 🎯 Когда использовать эти паттерны?

1. **Конвейеры** идеальны, когда у вас есть последовательность операций обработки данных, особенно если разные этапы могут выполняться параллельно.

2. **Fan-out** полезен, когда обработка каждого элемента занимает много времени и может выполняться независимо. Например, обработка изображений или запросы к API.

3. **Fan-in** необходим, когда данные приходят из нескольких источников, и вы хотите обрабатывать их через единый поток.

4. **Обработка отмены** критична для долгоживущих процессов, чтобы избежать утечки горутин и позволить системе корректно завершаться.

## 💡 Практические советы

1. **Всегда закрывайте каналы**: Использование `defer close(out)` гарантирует, что канал закроется даже при возникновении ошибок.

2. **Используйте sync.WaitGroup** для отслеживания активных горутин.

3. **Продумывайте обработку ошибок**: Можно создавать специальные структуры, содержащие результат и возможную ошибку.

4. **Не забывайте о контексте отмены**: В реальных приложениях используйте `context.Context` вместо самодельных каналов отмены.
