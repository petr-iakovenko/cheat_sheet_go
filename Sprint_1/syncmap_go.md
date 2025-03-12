# Когда использовать sync.Map

Ссылка на статью хабра - [sync.Map в Go 1.9](https://habr.com/ru/articles/338718/)

---

`sync.Map` в Go используется когда нужен потокобезопасный доступ к карте с высокой нагрузкой на чтение и редкими записями, в отличие от `map` с мьютексами, которые блокируют всю структуру при изменениях.

---

## Что такое sync.Map и когда её использовать 🤔

`sync.Map` — это специальная реализация карты (словаря) в Go, которая обеспечивает безопасность при одновременном доступе из нескольких горутин без необходимости добавлять внешние блокировки.

### Когда стоит использовать sync.Map: 🚀

1. Когда карта используется **несколькими горутинами** одновременно
2. Когда **операций чтения намного больше**, чем операций записи
3. При работе с картой, где **ключи добавляются один раз**, но потом многократно читаются
4. Когда вы хотите избежать явных блокировок (`mutex`)

### Простое сравнение с обычной картой: 📊

```go
// Обычная карта с мьютексом
var regularMap = make(map[string]int)
var mutex sync.Mutex

// Чтение и запись требуют блокировки всей карты
func updateRegularMap(key string, value int) {
    mutex.Lock()
    regularMap[key] = value
    mutex.Unlock()
}

// sync.Map - встроенная потокобезопасность
var syncMap sync.Map

// Чтение и запись без явных блокировок
func updateSyncMap(key string, value int) {
    syncMap.Store(key, value)
}
```

## Основные методы sync.Map с примерами 🛠️

### 1. Store(key, value) - сохранение значения

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var userScores sync.Map
    
    // Сохраняем баллы пользователей
    // Несколько горутин могут вызывать Store безопасно
    userScores.Store("Алексей", 100)  // Добавление нового ключа
    userScores.Store("Мария", 85)     // Добавление нового ключа
    userScores.Store("Алексей", 110)  // Обновление существующего ключа
    
    // Проверяем результат
    score, ok := userScores.Load("Алексей")
    if ok {
        fmt.Printf("Баллы Алексея: %d\n", score)  // Выведет: Баллы Алексея: 110
    }
}
```

### 2. Load(key) - чтение значения

```go
func main() {
    var userCache sync.Map
    
    // Добавим данные
    userCache.Store("user:1", "Иван Петров")
    
    // Получаем данные
    userName, exists := userCache.Load("user:1")
    if exists {
        fmt.Printf("Пользователь: %s\n", userName)  // Пользователь: Иван Петров
    } else {
        fmt.Println("Пользователь не найден")
    }
    
    // Попытка загрузки несуществующего ключа
    adminName, exists := userCache.Load("admin:1")
    if !exists {
        fmt.Println("Админ не найден в кэше")  // Будет выведено
    }
}
```

### 3. LoadOrStore(key, value) - атомарное "загрузить или сохранить"

```go
func main() {
    var counters sync.Map
    
    // LoadOrStore пытается загрузить ключ, если его нет - сохраняет новое значение
    // Возвращает: (сохраненное/загруженное значение, true если ключ уже существовал)
    
    // Первый вызов - ключа нет, сохраняем 1
    actual, loaded := counters.LoadOrStore("visitors", 1)
    fmt.Printf("Значение: %v, Существовал: %v\n", actual, loaded)
    // Выведет: Значение: 1, Существовал: false
    
    // Второй вызов - ключ уже есть, значение не изменится, загружаем существующее
    actual, loaded = counters.LoadOrStore("visitors", 100)
    fmt.Printf("Значение: %v, Существовал: %v\n", actual, loaded)
    // Выведет: Значение: 1, Существовал: true
    
    // Это идеально для инициализации кэша или счетчиков
}
```

### 4. Delete(key) - удаление значения

```go
func main() {
    var sessions sync.Map
    
    // Добавляем сессии
    sessions.Store("session:123", "user:456")
    
    // Проверяем существование
    _, exists := sessions.Load("session:123")
    fmt.Println("Сессия существует:", exists)  // true
    
    // Удаляем сессию
    sessions.Delete("session:123")
    
    // Проверяем после удаления
    _, exists = sessions.Load("session:123")
    fmt.Println("Сессия существует:", exists)  // false
}
```

### 5. Range(функция) - перебор всех значений

```go
func main() {
    var inventory sync.Map
    
    // Добавляем товары
    inventory.Store("яблоки", 150)
    inventory.Store("апельсины", 75)
    inventory.Store("бананы", 120)
    
    // Перебираем все товары и их количество
    fmt.Println("Товары на складе:")
    inventory.Range(func(key, value interface{}) bool {
        // Для каждой пары ключ-значение вызывается эта функция
        fmt.Printf("- %s: %d шт.\n", key, value)
        return true  // return true = продолжаем обход
        // return false = прерываем обход
    })
}
```

## Реальный пример использования sync.Map 🌟

Давайте рассмотрим пример кэша для веб-сервера, который обрабатывает запросы из разных горутин:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Имитация базы данных
func fetchUserFromDB(userID string) string {
    // В реальности здесь был бы запрос к базе данных
    time.Sleep(200 * time.Millisecond) // Имитация задержки БД
    return "Пользователь_" + userID
}

func main() {
    var userCache sync.Map // 🔥 Кэш пользователей, доступный из всех горутин
    var wg sync.WaitGroup
    
    // Функция получения пользователя (из кэша или БД)
    getUser := func(userID string) string {
        // Пытаемся взять из кэша
        if cachedUser, found := userCache.Load(userID); found {
            fmt.Printf("✅ Пользователь %s найден в кэше\n", userID)
            return cachedUser.(string)
        }
        
        // Если нет в кэше, идем в "базу данных"
        fmt.Printf("🔍 Загружаем пользователя %s из БД...\n", userID)
        user := fetchUserFromDB(userID)
        
        // Сохраняем в кэш для следующих запросов
        userCache.Store(userID, user)
        return user
    }
    
    // Симулируем 10 параллельных запросов к одним и тем же пользователям
    userIDs := []string{"1", "2", "3", "1", "2", "3", "1", "2", "3", "4"}
    
    wg.Add(len(userIDs))
    for _, id := range userIDs {
        go func(userID string) {
            defer wg.Done()
            
            user := getUser(userID)
            fmt.Printf("👤 Получен пользователь: %s\n", user)
            
            // Имитация работы с данными пользователя
            time.Sleep(100 * time.Millisecond)
        }(id)
    }
    
    wg.Wait()
    
    // Вывод статистики кэша
    fmt.Println("\n📊 Состояние кэша:")
    userCache.Range(func(key, value interface{}) bool {
        fmt.Printf("- ID: %s, Значение: %s\n", key, value)
        return true
    })
}
```

## Сравнение sync.Map с использованием обычной map и mutex ⚔️

```go
package main

import (
    "fmt"
    "sync"
    "testing"
)

func BenchmarkMapWithMutex(b *testing.B) {
    // Обычная карта с мьютексом
    m := make(map[int]int)
    var mutex sync.Mutex
    
    // Запускаем 100 горутин, каждая делает 1000 операций
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            // Каждая горутина делает операции чтения и редкие операции записи
            for j := 0; j < 1000; j++ {
                if j % 100 == 0 {  // Пишем редко (1% операций)
                    mutex.Lock()
                    m[workerID*1000+j] = j  // Запись
                    mutex.Unlock()
                } else {  // Читаем часто (99% операций)
                    mutex.Lock()
                    _ = m[workerID]  // Чтение
                    mutex.Unlock()
                }
            }
        }(i)
    }
    wg.Wait()
    
    fmt.Println("Map с мьютексом: записей =", len(m))
}

func BenchmarkSyncMap(b *testing.B) {
    // sync.Map
    var sm sync.Map
    
    // Та же нагрузка, что и выше
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            for j := 0; j < 1000; j++ {
                if j % 100 == 0 {  // Пишем редко (1% операций)
                    sm.Store(workerID*1000+j, j)  // Запись
                } else {  // Читаем часто (99% операций)
                    sm.Load(workerID)  // Чтение
                }
            }
        }(i)
    }
    wg.Wait()
    
    // Подсчет записей
    count := 0
    sm.Range(func(_, _ interface{}) bool {
        count++
        return true
    })
    fmt.Println("sync.Map: записей =", count)
}
```

## Ключевые моменты для понимания sync.Map 🔑

1. **Внутреннее устройство**: `sync.Map` использует две внутренние карты - одну для "только чтения" и другую для "грязных" записей, что обеспечивает высокую производительность при частых чтениях.

2. **Нет размера напрямую**: У `sync.Map` нет метода `len()`. Для подсчета элементов нужно использовать `Range()`.

3. **Приведение типов**: Поскольку `sync.Map` работает с `interface{}`, вам нужно приводить типы при извлечении значений:

   ```go
   value, ok := myMap.Load("key")
   if ok {
       stringValue := value.(string) // Приведение типа
   }
   ```

4. **Производительность**: `sync.Map` эффективнее обычной карты с мьютексом только при определенной нагрузке (много чтений, мало записей).

5. **Тест нагрузки**: Перед выбором `sync.Map` рекомендуется провести тест производительности для вашего конкретного сценария - иногда обычная карта с мьютексом может быть эффективнее.

## Преимущества sync.Map перед обычной картой с мьютексом 🏆

1. **Меньше блокировок**: `sync.Map` использует более тонкую стратегию блокировок, что позволяет нескольким горутинам читать одновременно без ожидания.

2. **Оптимизация для особых случаев**: Особенно эффективна для кэшей, где запись происходит один раз, а чтение - многократно.

3. **Простота использования**: Не требует ручного управления мьютексами, уменьшая риск ошибок (забытая блокировка или разблокировка).

4. **Атомарные операции**: Методы вроде `LoadOrStore` позволяют атомарно выполнять сложные операции без дополнительной синхронизации.

## Когда НЕ использовать sync.Map ⛔

1. **Большое количество записей**: Если ваша программа часто изменяет данные в карте, лучше использовать обычную карту с мьютексом.

2. **Одна горутина**: Если карта используется только в одной горутине, используйте обычную карту без синхронизации.

3. **Необходим len()**: Если вам часто нужно знать размер карты, `sync.Map` будет неудобна из-за отсутствия прямого метода `len()`.

4. **Высокая нагрузка записи**: В ситуациях, где количество операций записи сопоставимо с количеством чтений, обычная карта с мьютексом может быть эффективнее.