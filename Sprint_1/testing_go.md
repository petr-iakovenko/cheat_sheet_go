# Универсальное руководство по написанию юнит-тестов в Go с использованием Testify

## Содержание

1. [Введение в тестирование Go](#1-введение-в-тестирование-go) [⭐]
2. [Базовые концепции библиотеки Testify](#2-базовые-концепции-библиотеки-testify) [⭐]
3. [Структура юнит-теста](#3-структура-юнит-теста) [⭐]
4. [Табличные тесты](#4-табличные-тесты) [⭐⭐]
5. [Типы утверждений (assertions) и когда их использовать](#5-типы-утверждений-assertions-и-когда-их-использовать) [⭐⭐]
6. [Тестирование функций, возвращающих ошибки](#6-тестирование-функций-возвращающих-ошибки) [⭐⭐]
7. [Граничные случаи и нетривиальные сценарии](#7-граничные-случаи-и-нетривиальные-сценарии) [⭐⭐⭐]
8. [Тестирование сложных структур данных](#8-тестирование-сложных-структур-данных) [⭐⭐⭐]
9. [Мокирование зависимостей](#9-мокирование-зависимостей) [⭐⭐⭐⭐]
10. [Работа с временем и таймерами](#10-работа-с-временем-и-таймерами) [⭐⭐⭐]
11. [Параллельное выполнение тестов](#11-параллельное-выполнение-тестов) [⭐⭐⭐]
12. [Анализ покрытия тестами](#12-анализ-покрытия-тестами) [⭐⭐]
13. [Отладка тестов](#13-отладка-тестов) [⭐⭐⭐]
14. [Чек-лист для оценки качества тестов](#14-чек-лист-для-оценки-качества-тестов) [⭐]
15. [Частые ошибки и пути их решения](#15-частые-ошибки-и-пути-их-решения) [⭐⭐⭐]

## 1. Введение в тестирование Go

[⭐] _Базовый уровень_

Golang имеет встроенную поддержку тестирования через пакет `testing`. Тесты в Go - это обычные функции, следующие определенным соглашениям:

### Соглашения о именовании и размещении

- Файлы с тестами должны иметь суффикс `_test.go`
- Функции тестирования должны начинаться с префикса `Test`
- Функции тестирования должны принимать параметр `t *testing.T`

### Пример простейшего теста

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import (
    "testing"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### Запуск тестов

```bash
go test ./...       # Запуск всех тестов в проекте
go test ./pkg/...   # Запуск тестов в определенной директории
go test -v          # Подробный вывод
go test -run=TestAdd # Запуск конкретного теста
```

## 2. Базовые концепции библиотеки Testify

[⭐] _Базовый уровень_

Testify - это популярная библиотека для тестирования в Go, которая делает утверждения (assertions) более выразительными и информативными.

### Установка Testify

```bash
go get github.com/stretchr/testify
```

### Основные пакеты Testify

1. **assert** - содержит функции для проверки утверждений
2. **require** - похож на assert, но останавливает тест при первой неудачной проверке
3. **mock** - предоставляет инструменты для создания моков
4. **suite** - позволяет организовывать тесты в наборы

### Импорт пакетов

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

### Разница между assert и require

```go
// assert продолжит выполнение теста, даже если проверка не прошла
assert.Equal(t, 5, Add(2, 3))
assert.Equal(t, 7, Add(3, 4))  // Выполнится даже если предыдущий assert не прошел

// require остановит тест при первой неудачной проверке
require.Equal(t, 5, Add(2, 3))
require.Equal(t, 7, Add(3, 4))  // Не выполнится, если предыдущий require не прошел
```

## 3. Структура юнит-теста

[⭐] _Базовый уровень_

Хорошая структура юнит-теста должна следовать паттерну AAA (Arrange-Act-Assert):

### Паттерн AAA

```go
func TestSomeFunction(t *testing.T) {
    // Arrange (Подготовка)
    input1 := "test"
    input2 := 42
    expected := "result"
    
    // Act (Действие)
    actual := SomeFunction(input1, input2)
    
    // Assert (Проверка)
    assert.Equal(t, expected, actual, "Функция должна вернуть правильный результат")
}
```

### Именование тестов

Следуйте соглашению:

```
Test<Имя_Функции>_<Сценарий>
```

Примеры:

```go
TestAdd_PositiveNumbers
TestAdd_NegativeNumbers
TestAdd_ZeroValues
```

### Структурирование сложных тестов

Используйте подтесты с t.Run():

```go
func TestComplex(t *testing.T) {
    t.Run("Сценарий 1", func(t *testing.T) {
        // Тестовый код для сценария 1
    })
    
    t.Run("Сценарий 2", func(t *testing.T) {
        // Тестовый код для сценария 2
    })
}
```

## 4. Табличные тесты

[⭐⭐] _Средний уровень_

Табличные тесты - идиоматический способ в Go для тестирования функции с разными входными данными.

### Базовая структура табличного теста

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input1   string
        input2   int
        expected string
        wantErr  bool
    }{
        {
            name:     "Базовый случай",
            input1:   "hello",
            input2:   42,
            expected: "hello 42",
            wantErr:  false,
        },
        {
            name:     "Граничный случай",
            input1:   "",
            input2:   0,
            expected: " 0",
            wantErr:  false,
        },
        // Больше тестовых случаев...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            result, err := MyFunction(tt.input1, tt.input2)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

### Проверка ошибок в табличных тестах

Для функций, возвращающих ошибки, добавьте поля для проверки:

```go
tests := []struct {
    name          string
    input         string
    expected      string
    wantErr       bool
    expectedError string
}{
    {
        name:          "Ошибка валидации",
        input:         "",
        expected:      "",
        wantErr:       true,
        expectedError: "input cannot be empty",
    },
}

// В цикле проверки:
if tt.wantErr {
    assert.Error(t, err)
    if err != nil {
        assert.Contains(t, err.Error(), tt.expectedError)
    }
} else {
    assert.NoError(t, err)
    assert.Equal(t, tt.expected, result)
}
```

### Организация табличных тестов для лучшей читаемости

Для больших наборов тестовых случаев:

```go
// Определение типа для тестового случая
type addTestCase struct {
    name     string
    a, b     int
    expected int
}

// Группировка тестовых случаев по категориям
var positiveTests = []addTestCase{
    {"Оба положительные", 2, 3, 5},
    {"Большие числа", 1000, 2000, 3000},
}

var negativeTests = []addTestCase{
    {"Оба отрицательные", -2, -3, -5},
    {"Отрицательный и положительный", -5, 3, -2},
}

func TestAdd(t *testing.T) {
    // Запуск положительных тестов
    for _, tt := range positiveTests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, Add(tt.a, tt.b))
        })
    }
    
    // Запуск отрицательных тестов
    for _, tt := range negativeTests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, Add(tt.a, tt.b))
        })
    }
}
```

## 5. Типы утверждений (assertions) и когда их использовать

[⭐⭐] _Средний уровень_

Testify предоставляет множество функций для проверки утверждений. Вот основные из них и рекомендации по их использованию:

### Проверки равенства

```go
// Базовое сравнение (==)
assert.Equal(t, expected, actual, "Сообщение")

// Сравнение с приведением типов
assert.EqualValues(t, expected, actual, "Сообщение")

// Проверка неравенства
assert.NotEqual(t, notExpected, actual, "Сообщение")

// Глубокое сравнение для сложных типов
assert.ElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1}, "Порядок не важен")
assert.Subset(t, []int{1, 2, 3, 4}, []int{2, 4}, "Подмножество")
```

### Проверки ошибок

```go
// Проверка на наличие ошибки
assert.Error(t, err, "Должна быть ошибка")
assert.NoError(t, err, "Не должно быть ошибки")

// Проверка типа ошибки
assert.IsType(t, &MyErrorType{}, err, "Должна быть ошибка определенного типа")

// Проверка содержимого ошибки
assert.Contains(t, err.Error(), "expected text", "Ошибка должна содержать текст")
assert.ErrorContains(t, err, "expected text", "Ошибка должна содержать текст")
```

### Проверки значений nil/not nil

```go
assert.Nil(t, value, "Значение должно быть nil")
assert.NotNil(t, value, "Значение не должно быть nil")
```

### Проверки булевых значений

```go
assert.True(t, value, "Значение должно быть true")
assert.False(t, value, "Значение должно быть false")
```

### Проверка частичного соответствия

```go
// Строки
assert.Contains(t, "Hello World", "World", "Строка должна содержать подстроку")
assert.Prefix(t, "Hello", "Hell", "Строка должна начинаться с префикса")
assert.Suffix(t, "Hello", "llo", "Строка должна заканчиваться суффиксом")

// Коллекции
assert.Contains(t, []string{"a", "b", "c"}, "b", "Слайс должен содержать элемент")
assert.Contains(t, map[string]int{"a": 1, "b": 2}, "a", "Карта должна содержать ключ")
```

### Проверки для чисел

```go
assert.Greater(t, 5, 3, "5 должно быть больше 3")
assert.GreaterOrEqual(t, 5, 5, "5 должно быть больше или равно 5")
assert.Less(t, 3, 5, "3 должно быть меньше 5")
assert.LessOrEqual(t, 5, 5, "5 должно быть меньше или равно 5")
```

### Когда использовать require вместо assert

Используйте require, когда:

- Последующие проверки не имеют смысла, если эта проверка не прошла
- Продолжение теста может привести к панике
- Вы тестируете предусловие, необходимое для основного теста

```go
// Пример использования require
result, err := SomeFunction()
require.NoError(t, err, "Функция не должна возвращать ошибку")
require.NotNil(t, result, "Результат не должен быть nil")

// Теперь можно безопасно работать с result
assert.Equal(t, 5, result.Count)
```

## 6. Тестирование функций, возвращающих ошибки

[⭐⭐] _Средний уровень_

Go имеет идиоматический подход к обработке ошибок, и тесты должны учитывать возможные ошибки.

### Базовый шаблон для тестирования функций с ошибками

```go
func TestFunctionWithError(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result, err := FunctionWithError(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "expected", result)
}
```

### Проверка конкретных ошибок

```go
// Определенный тип ошибки
var ErrNotFound = errors.New("not found")

func TestWithSpecificError(t *testing.T) {
    _, err := FunctionThatReturnsSpecificError()
    
    assert.Error(t, err)
    assert.Equal(t, ErrNotFound, err)
    // Или
    assert.ErrorIs(t, err, ErrNotFound)
}

// Пользовательский тип ошибки
type ValidationError struct {
    Field string
    Msg   string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func TestWithCustomError(t *testing.T) {
    _, err := FunctionWithCustomError()
    
    assert.Error(t, err)
    
    // Проверка типа ошибки
    var validationErr ValidationError
    assert.True(t, errors.As(err, &validationErr))
    
    // Проверка полей ошибки
    assert.Equal(t, "username", validationErr.Field)
    assert.Equal(t, "cannot be empty", validationErr.Msg)
}
```

### Табличные тесты для функций с ошибками

```go
func TestWithErrorCases(t *testing.T) {
    tests := []struct {
        name           string
        input          string
        expectedResult string
        expectError    bool
        errorContains  string
    }{
        {
            name:           "Успешный случай",
            input:          "valid",
            expectedResult: "processed valid",
            expectError:    false,
        },
        {
            name:          "Пустой ввод",
            input:         "",
            expectError:   true,
            errorContains: "empty input",
        },
        {
            name:          "Недопустимые символы",
            input:         "inv@lid",
            expectError:   true,
            errorContains: "invalid characters",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ProcessWithError(tt.input)
            
            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorContains)
                assert.Empty(t, result)  // Проверка, что результат пуст при ошибке
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedResult, result)
            }
        })
    }
}
```

### Проверка взаимоисключающих ошибок и результатов

```go
// Функция либо возвращает результат, либо ошибку, никогда оба
func TestMutuallyExclusive(t *testing.T) {
    result, err := FunctionWithError(input)
    
    // Если ошибка есть, результат должен быть "нулевым"
    if err != nil {
        assert.Empty(t, result)
        // Проверки ошибки...
    } else {
        assert.NotEmpty(t, result)
        // Проверки результата...
    }
}
```

## 7. Граничные случаи и нетривиальные сценарии

[⭐⭐⭐] _Продвинутый уровень_

Хорошие тесты проверяют не только типичные случаи, но и граничные условия, которые могут вызвать проблемы.

### Типы граничных случаев для тестирования

1. **Пустые значения**:

   ```go
   t.Run("Empty string", func(t *testing.T) {
       result := ProcessString("")
       assert.Equal(t, "default", result)
   })
   ```

2. **Нулевые значения**:

   ```go
   t.Run("Zero values", func(t *testing.T) {
       result := Calculate(0, 0)
       assert.Equal(t, 0, result)
   })
   ```

3. **Максимальные/минимальные значения**:

   ```go
   t.Run("Max int", func(t *testing.T) {
       result, err := Process(math.MaxInt32)
       assert.NoError(t, err)
       assert.Equal(t, expected, result)
   })
   
   t.Run("Min int", func(t *testing.T) {
       result, err := Process(math.MinInt32)
       assert.NoError(t, err)
       assert.Equal(t, expected, result)
   })
   ```

4. **Очень большие коллекции**:

   ```go
   t.Run("Large slice", func(t *testing.T) {
       // Создаем слайс с 10000 элементами
       largeSlice := make([]int, 10000)
       for i := range largeSlice {
           largeSlice[i] = i
       }
       
       result := ProcessSlice(largeSlice)
       assert.Equal(t, expected, result)
   })
   ```

5. **Специальные символы**:

   ```go
   t.Run("Special characters", func(t *testing.T) {
       result := ProcessString("!@#$%^&*()")
       assert.Equal(t, expected, result)
   })
   
   t.Run("Unicode characters", func(t *testing.T) {
       result := ProcessString("こんにちは世界") // Hello world на японском
       assert.Equal(t, expected, result)
   })
   ```

6. **Недопустимые входные данные**:

   ```go
   t.Run("Invalid format", func(t *testing.T) {
       _, err := ParseDate("not-a-date")
       assert.Error(t, err)
       assert.Contains(t, err.Error(), "invalid date format")
   })
   ```

### Нетривиальные сценарии

1. **Конкурентный доступ**:

   ```go
   t.Run("Concurrent access", func(t *testing.T) {
       cache := NewCache()
       const goroutines = 100
       
       var wg sync.WaitGroup
       wg.Add(goroutines)
       
       for i := 0; i < goroutines; i++ {
           go func(id int) {
               defer wg.Done()
               key := fmt.Sprintf("key-%d", id)
               cache.Set(key, id)
               val, ok := cache.Get(key)
               assert.True(t, ok)
               assert.Equal(t, id, val)
           }(i)
       }
       
       wg.Wait()
   })
   ```

2. **Циклические зависимости**:

   ```go
   t.Run("Cyclic dependencies", func(t *testing.T) {
       graph := NewGraph()
       graph.AddEdge("A", "B")
       graph.AddEdge("B", "C")
       graph.AddEdge("C", "A")
       
       hasCycle, cycle := graph.DetectCycle()
       assert.True(t, hasCycle)
       assert.Contains(t, cycle, "A")
       assert.Contains(t, cycle, "B")
       assert.Contains(t, cycle, "C")
   })
   ```

3. **Повторное использование ресурсов**:

   ```go
   t.Run("Resource reuse", func(t *testing.T) {
       // Создаем ресурс
       resource := NewResource()
       defer resource.Close()
       
       // Используем ресурс в первый раз
       result1, err := resource.Process("input1")
       assert.NoError(t, err)
       assert.Equal(t, "output1", result1)
       
       // Используем повторно
       result2, err := resource.Process("input2")
       assert.NoError(t, err)
       assert.Equal(t, "output2", result2)
   })
   ```

## 8. Тестирование сложных структур данных

[⭐⭐⭐] _Продвинутый уровень_

При тестировании функций, работающих со сложными структурами данных, нужен особый подход.

### Тестирование функций, работающих с JSON

```go
func TestParseJSON(t *testing.T) {
    // JSON строка
    jsonStr := `{"name":"John","age":30,"address":{"city":"New York","country":"USA"}}`
    
    // Ожидаемая структура
    expected := Person{
        Name: "John",
        Age:  30,
        Address: Address{
            City:    "New York",
            Country: "USA",
        },
    }
    
    // Парсинг
    var actual Person
    err := json.Unmarshal([]byte(jsonStr), &actual)
    
    // Проверки
    assert.NoError(t, err)
    assert.Equal(t, expected, actual)
    
    // Для более сложных структур используйте DeepEqual
    assert.True(t, reflect.DeepEqual(expected, actual))
}
```

### Тестирование функций, работающих с деревьями

```go
func TestTreeOperations(t *testing.T) {
    // Создаем тестовое дерево
    //      1
    //     / \
    //    2   3
    //   / \
    //  4   5
    root := &Node{
        Value: 1,
        Left: &Node{
            Value: 2,
            Left:  &Node{Value: 4},
            Right: &Node{Value: 5},
        },
        Right: &Node{Value: 3},
    }
    
    // Тест in-order traversal
    result := InOrderTraversal(root)
    assert.Equal(t, []int{4, 2, 5, 1, 3}, result)
    
    // Тест поиска
    found := FindNode(root, 5)
    assert.NotNil(t, found)
    assert.Equal(t, 5, found.Value)
    
    // Тест поиска несуществующего значения
    notFound := FindNode(root, 6)
    assert.Nil(t, notFound)
}
```

### Тестирование функций, работающих с графами

```go
func TestGraph(t *testing.T) {
    // Создаем тестовый граф
    //    A --- B
    //    |     |
    //    |     |
    //    C --- D
    g := NewGraph()
    g.AddVertex("A")
    g.AddVertex("B")
    g.AddVertex("C")
    g.AddVertex("D")
    
    g.AddEdge("A", "B")
    g.AddEdge("A", "C")
    g.AddEdge("B", "D")
    g.AddEdge("C", "D")
    
    // Тест поиска пути
    path, found := g.FindPath("A", "D")
    assert.True(t, found)
    assert.Contains(t, [][]string{
        {"A", "B", "D"},
        {"A", "C", "D"},
    }, path)
    
    // Тест на отсутствие пути
    g = NewGraph()
    g.AddVertex("A")
    g.AddVertex("B")
    
    _, found = g.FindPath("A", "B")
    assert.False(t, found)
}
```

### Тестирование сортировки и работы с коллекциями

```go
func TestSort(t *testing.T) {
    // Тест с разными входными данными
    for _, tc := range []struct {
        name     string
        input    []int
        expected []int
    }{
        {"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
        {"Reverse order", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
        {"Random order", []int{3, 1, 5, 2, 4}, []int{1, 2, 3, 4, 5}},
        {"Duplicates", []int{3, 1, 3, 2, 1}, []int{1, 1, 2, 3, 3}},
        {"Empty slice", []int{}, []int{}},
        {"Single element", []int{1}, []int{1}},
    } {
        t.Run(tc.name, func(t *testing.T) {
            result := Sort(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func TestMapOperations(t *testing.T) {
    // Создаем тестовую мапу
    m := make(map[string]int)
    m["a"] = 1
    m["b"] = 2
    
    // Тест операций с мапой
    result := ProcessMap(m)
    
    // В зависимости от ожидаемого результата:
    assert.Equal(t, 3, result)
    // или
    assert.Equal(t, map[string]int{"a": 2, "b": 4}, result)
}
```

## 9. Мокирование зависимостей

[⭐⭐⭐⭐] _Экспертный уровень_

При тестировании кода с внешними зависимостями часто требуется "мокировать" эти зависимости.

### Тестирование с использованием интерфейсов

```go
// Определение интерфейса
type DataStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

// Функция, использующая интерфейс
func ProcessData(store DataStore, key string) (string, error) {
    data, err := store.Get(key)
    if err != nil {
        return "", err
    }
    
    processed := strings.ToUpper(data)
    return processed, nil
}

// Мок-реализация для тестирования
type MockDataStore struct {
    mock.Mock
}

func (m *MockDataStore) Get(key string) (string, error) {
    args := m.Called(key)
    return args.String(0), args.Error(1)
}

func (m *MockDataStore) Set(key, value string) error {
    args := m.Called(key, value)
    return args.Error(0)
}

// Тест с использованием мока
func TestProcessData(t *testing.T) {
    // Создаем мок
    mockStore := new(MockDataStore)
    
    // Настраиваем поведение мока
    mockStore.On("Get", "existingKey").Return("data", nil)
    mockStore.On("Get", "missingKey").Return("", errors.New("key not found"))
    
    // Тест успешного случая
    t.Run("Existing key", func(t *testing.T) {
        result, err := ProcessData(mockStore, "existingKey")
        assert.NoError(t, err)
        assert.Equal(t, "DATA", result)
    })
    
    // Тест случая с ошибкой
    t.Run("Missing key", func(t *testing.T) {
        result, err := ProcessData(mockStore, "missingKey")
        assert.Error(t, err)
        assert.Equal(t, "key not found", err.Error())
        assert.Empty(t, result)
    })
    
    // Проверяем, что все ожидаемые вызовы произошли
    mockStore.AssertExpectations(t)
}

### Использование встроенных интерфейсов Go

Часто удобнее создавать собственные простые моки, особенно для стандартных интерфейсов:

```go
// Мокирование io.Reader
type MockReader struct {
    ReadFunc func(p []byte) (n int, err error)
}

func (m MockReader) Read(p []byte) (n int, err error) {
    return m.ReadFunc(p)
}

func TestProcessReader(t *testing.T) {
    // Создаем мок, возвращающий данные
    mock := MockReader{
        ReadFunc: func(p []byte) (int, error) {
            data := []byte("test data")
            copy(p, data)
            return len(data), nil
        },
    }
    
    result, err := ProcessReader(mock)
    assert.NoError(t, err)
    assert.Equal(t, "TEST DATA", result)
    
    // Создаем мок, возвращающий ошибку
    mockError := MockReader{
        ReadFunc: func(p []byte) (int, error) {
            return 0, errors.New("read error")
        },
    }
    
    result, err = ProcessReader(mockError)
    assert.Error(t, err)
    assert.Empty(t, result)
}
```

### Мокирование HTTP запросов

```go
// Мокирование http.Client
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

func TestFetchData(t *testing.T) {
    // Создаем успешный ответ
    mockClient := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            // Проверяем URL запроса
            assert.Equal(t, "https://api.example.com/data", req.URL.String())
            
            // Создаем мок-ответ
            responseBody := `{"status": "ok", "data": "test"}`
            return &http.Response{
                StatusCode: 200,
                Body:       ioutil.NopCloser(strings.NewReader(responseBody)),
            }, nil
        },
    }
    
    service := NewService(mockClient)
    data, err := service.FetchData()
    
    assert.NoError(t, err)
    assert.Equal(t, "test", data)
    
    // Тест ошибки запроса
    mockErrorClient := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return nil, errors.New("network error")
        },
    }
    
    service = NewService(mockErrorClient)
    data, err = service.FetchData()
    
    assert.Error(t, err)
    assert.Empty(t, data)
}
```

### Использование httptest для тестирования HTTP серверов

```go
func TestHTTPHandler(t *testing.T) {
    // Создаем тестовый HTTP сервер
    handler := http.HandlerFunc(MyHandler)
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Отправляем запрос к тестовому серверу
    resp, err := http.Get(server.URL + "/test")
    assert.NoError(t, err)
    
    // Проверяем статус ответа
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Проверяем тело ответа
    body, err := ioutil.ReadAll(resp.Body)
    assert.NoError(t, err)
    assert.Equal(t, "Expected response", string(body))
}
```

## 10. Работа с временем и таймерами

[⭐⭐⭐] _Продвинутый уровень_

Тестирование кода, который зависит от времени, может быть сложным. Есть несколько подходов:

### Инъекция зависимости времени

```go
// Вместо прямого использования time.Now()
type timeProvider interface {
    Now() time.Time
}

type realTimeProvider struct{}

func (r realTimeProvider) Now() time.Time {
    return time.Now()
}

type Service struct {
    timeProvider timeProvider
}

func NewService(tp timeProvider) *Service {
    return &Service{timeProvider: tp}
}

func (s *Service) IsDuringBusinessHours() bool {
    now := s.timeProvider.Now()
    hour := now.Hour()
    return hour >= 9 && hour < 17
}

// В тесте
func TestBusinessHours(t *testing.T) {
    // Мок времени
    mockTime := struct{ timeProvider }{}
    
    // Устанавливаем разное время для тестов
    tests := []struct {
        name     string
        mockTime time.Time
        expected bool
    }{
        {
            name:     "During business hours",
            mockTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
            expected: true,
        },
        {
            name:     "Before business hours",
            mockTime: time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC),
            expected: false,
        },
        {
            name:     "After business hours",
            mockTime: time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC),
            expected: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockTime := struct {
                timeProvider
            }{
                timeProvider: &struct {
                    timeProvider
                }{
                    Now: func() time.Time {
                        return tt.mockTime
                    },
                },
            }
            
            service := NewService(mockTime)
            result := service.IsDuringBusinessHours()
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Использование подменных функций для тестирования таймаутов

```go
func TestWithTimeout(t *testing.T) {
    // Заменяем реальный sleep на мгновенный
    originalSleep := sleep
    defer func() { sleep = originalSleep }()
    
    sleep = func(d time.Duration) {
        // Ничего не делаем, не спим
    }
    
    start := time.Now()
    result := FunctionWithLongTimeout()
    elapsed := time.Since(start)
    
    assert.True(t, elapsed < time.Second, "Функция должна выполниться мгновенно с замененным sleep")
    assert.Equal(t, expected, result)
}

// Для функций, использующих time.After
func TestWithTimer(t *testing.T) {
    // Заменяем time.After на функцию, возвращающую канал, который сразу закрывается
    originalAfter := timeAfter
    defer func() { timeAfter = originalAfter }()
    
    timeAfter = func(d time.Duration) <-chan time.Time {
        ch := make(chan time.Time, 1)
        ch <- time.Now() // Сразу отправляем время
        return ch
    }
    
    start := time.Now()
    result := FunctionWithTimer()
    elapsed := time.Since(start)
    
    assert.True(t, elapsed < time.Second, "Функция должна выполниться мгновенно с замененным таймером")
    assert.Equal(t, expected, result)
}
```

### Тестирование периодических задач

```go
func TestPeriodicTask(t *testing.T) {
    // Счетчик вызовов
    counter := 0
    
    // Функция, которая увеличивает счетчик
    task := func() {
        counter++
    }
    
    // Создаем ticker с очень маленьким интервалом
    ticker := time.NewTicker(1 * time.Millisecond)
    defer ticker.Stop()
    
    // Запускаем задачу на короткое время
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer cancel()
    
    // Запускаем периодическую задачу
    go RunPeriodic(ctx, ticker, task)
    
    // Ждем завершения контекста
    <-ctx.Done()
    
    // Проверяем, что задача выполнялась несколько раз
    assert.Greater(t, counter, 1, "Задача должна выполниться несколько раз")
}
```

## 11. Параллельное выполнение тестов

[⭐⭐⭐] _Продвинутый уровень_

Go позволяет запускать тесты параллельно для ускорения выполнения.

### Маркировка тестов как параллельных

```go
func TestParallelFunction1(t *testing.T) {
    // Помечаем тест как параллельный
    t.Parallel()
    
    // Тест код...
}

func TestParallelFunction2(t *testing.T) {
    t.Parallel()
    
    // Тест код...
}
```

### Параллельные подтесты

```go
func TestParallelSubtests(t *testing.T) {
    // Тесты верхнего уровня не должны быть параллельными
    // если они содержат параллельные подтесты
    
    tests := []struct {
        name  string
        input int
    }{
        {"Case 1", 1},
        {"Case 2", 2},
        {"Case 3", 3},
    }
    
    for _, tt := range tests {
        tt := tt // Важно: создаем локальную копию переменной цикла
        
        t.Run(tt.name, func(t *testing.T) {
            // Помечаем подтест как параллельный
            t.Parallel()
            
            // Тестовый код...
            result := SomeFunction(tt.input)
            assert.Equal(t, tt.input*2, result)
        })
    }
}
```

### Синхронизация в параллельных тестах

```go
func TestConcurrentAccess(t *testing.T) {
    // Создаем объект для тестирования
    cache := NewThreadSafeCache()
    
    // Количество горутин
    const workers = 100
    
    // Для синхронизации
    var wg sync.WaitGroup
    wg.Add(workers)
    
    // Для трекинга ошибок
    errs := make(chan error, workers)
    
    // Запускаем конкурентные операции
    for i := 0; i < workers; i++ {
        go func(id int) {
            defer wg.Done()
            
            key := fmt.Sprintf("key-%d", id)
            value := fmt.Sprintf("value-%d", id)
            
            // Выполняем операции с кешем
            cache.Set(key, value)
            
            // Проверяем результат
            got, ok := cache.Get(key)
            if !ok {
                errs <- fmt.Errorf("key %s not found", key)
                return
            }
            
            if got != value {
                errs <- fmt.Errorf("expected %s, got %s", value, got)
                return
            }
        }(i)
    }
    
    // Ждем завершения всех горутин
    wg.Wait()
    close(errs)
    
    // Проверяем ошибки
    var errList []string
    for err := range errs {
        errList = append(errList, err.Error())
    }
    
    assert.Empty(t, errList, "Должно быть 0 ошибок, получено: %v", errList)
}
```

### Предостережения при использовании параллельных тестов

- Избегайте общего состояния между параллельными тестами
- Используйте t.TempDir() для создания временных директорий
- Не забывайте о том, что t.Cleanup() выполняется после завершения теста
- Помните о возможных гонках данных

## 12. Анализ покрытия тестами

[⭐⭐] _Средний уровень_

Go предоставляет встроенный инструмент для анализа покрытия кода тестами.

### Запуск тестов с анализом покрытия

```bash
# Запуск тестов с покрытием
go test -cover ./...

# Генерация файла с данными о покрытии
go test -coverprofile=coverage.out ./...

# Просмотр отчета в HTML формате
go tool cover -html=coverage.out

# Просмотр отчета в терминале
go tool cover -func=coverage.out
```

### Использование HTML-отчета о покрытии

HTML-отчет предоставляет визуальное представление о покрытии кода:

1. **Запуск генерации HTML-отчета**:

   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

2. **Интерпретация цветовой схемы**:
   - **Зеленый**: Строки кода, которые были выполнены во время тестов
   - **Красный**: Строки кода, которые не были выполнены
   - **Серый**: Неисполняемый код (комментарии, пустые строки)

3. **Анализ информации в отчете**:
   - **Общий процент покрытия**: отображается в верхней части отчета
   - **Покрытие по файлам**: можно переключаться между файлами проекта
   - **Детальное покрытие функций**: показывает какие ветви условных операторов выполнялись

4. **Использование отчета для улучшения тестов**:
   - Выявление непокрытых участков кода
   - Определение сложных условных ветвлений, требующих дополнительных тестов
   - Обнаружение недостаточно протестированных функций

### Пример анализа HTML-отчета

Для функции Hello из нашего примера HTML-отчет мог бы показать:

```go
func Hello(name string, language string) (string, error) {
    if name == "" {               // Покрыто (зеленый)
        name = "World"            // Покрыто (зеленый)
    }
    
    prefix := ""                  // Покрыто (зеленый)
    
    switch language {             // Покрыто (зеленый)
    case "english":               // Покрыто (зеленый)
        prefix = "Hello"          // Покрыто (зеленый)
    case "spanish":               // Не покрыто (красный) - нет теста для испанского
        prefix = "Hola"           // Не покрыто (красный)
    case "german":                // Не покрыто (красный) - нет теста для немецкого
        prefix = "Hallo"          // Не покрыто (красный)
    default:                      // Покрыто (зеленый)
        return "", errors.New("need to provide a supported language") // Покрыто (зеленый)
    }
    
    return prefix + " " + name, nil // Покрыто (зеленый)
}
```

На основе такого отчета мы могли бы добавить тесты для непокрытых веток языков:

```go
// Добавить в тестовые случаи:
{"test5", "Juan", "spanish", "Hola Juan", false, ""},
{"test6", "Hans", "german", "Hallo Hans", false, ""},
```

### Интеграция с IDE (VSCode)

В VSCode можно настроить визуализацию покрытия кода:

1. **Установка расширения**:
   - Установите расширение "Go" от официальной команды Go
   - Или расширение "Coverage Gutters" для более широких возможностей

2. **Активация покрытия в VSCode**:
   - Запустите тесты с покрытием: `go test -coverprofile=coverage.out ./...`
   - Откройте палитру команд: `Ctrl+Shift+P` (Windows/Linux) или `Cmd+Shift+P` (Mac)
   - Выполните команду: "Go: Toggle Coverage in Current Package"

3. **Использование визуального покрытия**:
   - Просматривайте покрытие прямо в редакторе кода
   - Непокрытые строки будут подсвечены красным цветом
   - Покрытые строки будут подсвечены зеленым цветом

4. **Анализ покрытия в режиме реального времени**:
   - Вносите изменения в тесты
   - Перезапускайте генерацию отчета
   - Наблюдайте, как меняется покрытие

### Целевые показатели покрытия

Рекомендуемые целевые показатели покрытия кода:

- **80-90%** - отличное покрытие для большинства проектов
- **70-80%** - хорошее покрытие
- **>90%** - может потребоваться для критически важного кода

Однако помните, что 100% покрытие не гарантирует отсутствие ошибок. Важно не только количество, но и качество тестов.

### Настройка покрытия в CI/CD

Добавьте проверку покрытия в ваш CI/CD пайплайн:

```yaml
# .github/workflows/go.yml пример для GitHub Actions
name: Go Test and Coverage

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: '1.20'
    - name: Test with coverage
      run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
    - name: Check coverage threshold
      run: |
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        if (( $(echo "$coverage < 70" | bc -l) )); then
          echo "Code coverage is below threshold: $coverage% < 70%"
          exit 1
        fi
```

### Интерпретация отчета о покрытии

```
github.com/example/package/file.go:  SomeFunc          75.0%
github.com/example/package/file.go:  AnotherFunc       100.0%
github.com/example/package/file.go:  ComplexFunc       45.5%
total:                               (statements)      73.8%
```

### Стратегии увеличения покрытия

1. **Сначала покрывайте критические пути**:

   ```go
   // Пример функции с условными ветвлениями
   func ProcessOrder(order Order) (string, error) {
       // Проверка базовых условий
       if order.ID == "" {
           return "", errors.New("order ID is required")
       }
       
       if order.Amount <= 0 {
           return "", errors.New("order amount must be positive")
       }
       
       // Обработка заказа
       return fmt.Sprintf("Processed order %s for $%.2f", order.ID, order.Amount), nil
   }
   
   // Тесты
   func TestProcessOrder(t *testing.T) {
       tests := []struct {
           name    string
           order   Order
           want    string
           wantErr bool
       }{
           {
               "Valid order",
               Order{ID: "123", Amount: 99.99},
               "Processed order 123 for $99.99",
               false,
           },
           {
               "Empty ID",
               Order{ID: "", Amount: 99.99},
               "",
               true,
           },
           {
               "Zero amount",
               Order{ID: "123", Amount: 0},
               "",
               true,
           },
           {
               "Negative amount",
               Order{ID: "123", Amount: -10},
               "",
               true,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got, err := ProcessOrder(tt.order)
               
               if tt.wantErr {
                   assert.Error(t, err)
               } else {
                   assert.NoError(t, err)
                   assert.Equal(t, tt.want, got)
               }
           })
       }
   }
   ```

2. **Проверяйте обработку ошибок**:

   ```go
   // Проверяем, что функция корректно обрабатывает ошибки от зависимостей
   func TestHandleRepositoryError(t *testing.T) {
       // Создаем мок-репозиторий, который возвращает ошибку
       mockRepo := new(MockRepository)
       mockRepo.On("FindByID", mock.Anything).Return(nil, errors.New("database error"))
       
       service := NewService(mockRepo)
       
       // Проверяем, что сервис правильно обрабатывает ошибку репозитория
       result, err := service.GetItemByID("123")
       
       assert.Error(t, err)
       assert.Contains(t, err.Error(), "failed to fetch item")
       assert.Nil(t, result)
   }
   ```

3. **Включайте тестирование конкурентных операций**:

   ```go
   func TestConcurrentOperations(t *testing.T) {
       counter := NewAtomicCounter()
       
       var wg sync.WaitGroup
       const numOperations = 1000
       
       // Увеличиваем счетчик из нескольких горутин
       for i := 0; i < numOperations; i++ {
           wg.Add(1)
           go func() {
               defer wg.Done()
               counter.Increment()
           }()
       }
       
       wg.Wait()
       
       // Проверяем, что все операции были учтены
       assert.Equal(t, numOperations, counter.Value())
   }
   ```

## 13. Отладка тестов

[⭐⭐⭐] _Продвинутый уровень_

Иногда тесты требуют отладки. Вот несколько полезных техник:

### Использование t.Log для отладочной информации

```go
func TestComplexFunction(t *testing.T) {
    input := "test input"
    
    t.Logf("Testing with input: %s", input)
    
    result, err := ComplexFunction(input)
    
    t.Logf("Got result: %v, err: %v", result, err)
    
    assert.NoError(t, err)
    assert.Equal(t, "expected output", result)
}
```

### Запуск тестов с повышенным уровнем вывода

```bash
go test -v ./...
```

### Запуск конкретного теста или подтеста

```bash
# Запуск конкретного теста
go test -run TestMyFunction

# Запуск конкретного подтеста
go test -run TestMyFunction/Subtest_Name

# Запуск тестов, соответствующих регулярному выражению
go test -run "TestMy.*"
```

### Использование delve для отладки тестов

```bash
# Установка delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Отладка теста
dlv test github.com/example/package -- -test.run TestMyFunction

# В консоли delve:
(dlv) break TestMyFunction
(dlv) continue
(dlv) next
(dlv) print variable
```

### Отладка с помощью testify/assert/require

require.Equal и другие методы из пакета require прекращают выполнение теста при ошибке, что помогает быстрее найти проблему:

```go
func TestWithRequire(t *testing.T) {
    // Подготовка
    input := "test"
    
    // Шаг 1: Проверяем предусловие
    result1 := Step1(input)
    require.NotNil(t, result1, "Шаг 1 не должен возвращать nil")
    
    // Шаг 2: Дальнейшая обработка
    result2 := Step2(result1)
    require.NotEmpty(t, result2, "Шаг 2 не должен возвращать пустой результат")
    
    // Шаг 3: Финальная проверка
    assert.Equal(t, "expected", result2)
}
```

### Использование строк Go для отладки сложных структур данных

```go
func TestComplexStructure(t *testing.T) {
    got := ComplexFunction()
    want := ComplexStruct{/* ... */}
    
    // При несовпадении вывести обе структуры для сравнения
    if !reflect.DeepEqual(got, want) {
        t.Errorf("Структуры не совпадают:\nПолучено: %#v\nОжидалось: %#v", got, want)
    }
}
```

## 14. Чек-лист для оценки качества тестов

[⭐] _Базовый уровень_

Используйте этот чек-лист для проверки качества ваших тестов:

### Основной чек-лист

- [ ] Все ли публичные функции покрыты тестами?
- [ ] Тесты проверяют как "happy path", так и обработку ошибок?
- [ ] Тесты проверяют граничные случаи (пустые строки, нулевые значения и т.д.)?
- [ ] Тесты понятно именованы и структурированы?
- [ ] Используются ли подтесты (t.Run) для группировки связанных тестовых случаев?
- [ ] Тесты не зависят друг от друга?
- [ ] Тесты достаточно быстрые?
- [ ] Тесты детерминированы (всегда дают одинаковый результат)?
- [ ] Тесты не содержат логику, которая требует своих собственных тестов?
- [ ] Покрытие кода тестами удовлетворительное?

### Продвинутый чек-лист

- [ ] Функции с побочными эффектами тестируются с использованием моков?
- [ ] Тесты для конкурентного кода проверяют гонки данных?
- [ ] Сложные сценарии разбиты на понятные подтесты?
- [ ] Тесты не зависят от окружения (временные зоны, локаль и т.д.)?
- [ ] Тесты используют assert/require для понятных сообщений об ошибках?
- [ ] Тесты запускаются в CI/CD пайплайне?
- [ ] Тесты соответствуют паттерну Arrange-Act-Assert?
- [ ] Табличные тесты используются для проверки различных входных данных?
- [ ] Тесты с зависимостью от времени используют мокирование?
- [ ] Функции, работающие с сетью или файловой системой, тестируются с использованием моков или фикстур?

## 15. Частые ошибки и пути их решения

[⭐⭐⭐] _Продвинутый уровень_

### Таблица частых ошибок и их решений

| Ошибка | Симптом | Решение |
|--------|---------|---------|
| Несоответствие между ожиданием ошибки и ее проверкой | `Expected nil, but got error` или наоборот | Проверьте логику проверки ошибок, используйте `assert.Error` или `assert.NoError` соответственно |
| Проверка только наличия ошибки, без проверки ее содержимого | Тест проходит с неправильным сообщением об ошибке | Добавьте `assert.Equal(t, "expected message", err.Error())` или `assert.Contains(t, err.Error(), "expected part")` |
| Зависимость тестов друг от друга | Тесты проходят при индивидуальном запуске, но не при запуске всех тестов | Избегайте глобальных состояний, используйте `t.Run` с индивидуальной подготовкой |
| Неправильное использование моков | Тесты с моками проходят, но реальный код падает | Убедитесь, что моки точно имитируют реальное поведение |
| Отсутствие проверки граничных случаев | Тесты проходят, но реальный код падает на граничных случаях | Добавьте тесты для пустых значений, максимальных/минимальных значений, специальных символов |
| Проблемы с конкурентным кодом | Случайные падения тестов | Используйте `-race` флаг, добавьте синхронизацию, учитывайте гонки данных |
| Нестабильные тесты с временем | Тесты иногда проходят, иногда падают | Замените `time.Now()` на инъекцию времени, моки для таймеров |
| Невоспроизводимые тесты | Тесты проходят на одной машине, но падают на другой | Исключите зависимости от окружения, времени, локали и т.д. |
| Слишком медленные тесты | Тесты выполняются долго | Используйте моки вместо реальных сервисов, параллельное выполнение тестов |
| Запутанные сообщения об ошибках | Трудно понять, почему тест падает | Добавьте информативные сообщения к assert/require вызовам |

### Детальные решения для типичных проблем

#### 1. Проблема с проверкой ошибок

Неправильно:

```go
// Ожидается ошибка, но проверяем, что ее нет
assert.NoError(t, err)
```

Правильно:

```go
// Если ожидается ошибка
if tt.expectError {
    assert.Error(t, err)
    if err != nil {
        assert.Contains(t, err.Error(), tt.errorMessage)
    }
} else {
    assert.NoError(t, err)
}
```

#### 2. Проблема с нестабильными тестами с временем

Нестабильный код:

```go
func IsExpired() bool {
    return time.Now().After(expiryTime)
}

func TestIsExpired(t *testing.T) {
    // Зависит от текущего времени
    result := IsExpired()
    assert.False(t, result)
}
```

Стабильное решение:

```go
func IsExpired(now time.Time) bool {
    return now.After(expiryTime)
}

func TestIsExpired(t *testing.T) {
    // Фиксированное время для теста
    fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
    expiryTime := fixedTime.Add(1 * time.Hour)
    
    // Проверка до истечения срока
    result := IsExpired(fixedTime, expiryTime)
    assert.False(t, result)
    
    // Проверка после истечения срока
    afterExpiry := fixedTime.Add(2 * time.Hour)
    result = IsExpired(afterExpiry, expiryTime)
    assert.True(t, result)
}
```

#### 3. Проблема с зависимостями между тестами

Проблемный код:

```go
// Глобальная переменная
var globalCounter int

func TestIncrementGlobal(t *testing.T) {
    globalCounter = 0 // Сброс
    IncrementGlobal()
    assert.Equal(t, 1, globalCounter)
}

func TestDoubleIncrementGlobal(t *testing.T) {
    // Зависит от начального значения globalCounter
    IncrementGlobal()
    IncrementGlobal()
    assert.Equal(t, 2, globalCounter) // Может быть 3, если первый тест уже выполнился
}
```

Решение:

```go
func TestIncrementGlobal(t *testing.T) {
    // Локальная копия для теста
    counter := 0
    Increment(&counter)
    assert.Equal(t, 1, counter)
}

func TestDoubleIncrementGlobal(t *testing.T) {
    // Независимая копия
    counter := 0
    Increment(&counter)
    Increment(&counter)
    assert.Equal(t, 2, counter)
}
```

#### 4. Проблема с конкурентным кодом

Проблемный код:

```go
func TestConcurrentMap(t *testing.T) {
    // Обычная карта не является потокобезопасной
    m := make(map[string]int)
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("key-%d", i)
            m[key] = i // Гонка данных!
        }(i)
    }
    
    wg.Wait()
    // Тест может упасть или пройти непредсказуемо
}
```

Решение:

```go
func TestConcurrentMap(t *testing.T) {
    // Используем sync.Map или мьютекс
    var m sync.Map
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            key := fmt.Sprintf("key-%d", i)
            m.Store(key, i) // Потокобезопасно
        }(i)
    }
    
    wg.Wait()
    
    // Проверяем результаты
    for i := 0; i < 100; i++ {
        key := fmt.Sprintf("key-%d", i)
        value, ok := m.Load(key)
        assert.True(t, ok)
        assert.Equal(t, i, value)
    }
}
```

#### 5. Проблема с мокированием

Недостаточный мок:

```go
// Мок только для успешного случая
mockRepo.On("FindByID", "123").Return(item, nil)

// Тест падает, если вызвать с другим ID
service.GetItem("456") // Паника или неожиданное поведение
```

Надежный мок:

```go
// Мок для любого ID с соответствующим поведением
mockRepo.On("FindByID", mock.AnythingOfType("string")).Return(
    func(id string) *Item {
        if id == "123" {
            return &Item{ID: id, Name: "Test Item"}
        }
        return nil
    },
    func(id string) error {
        if id == "123" {
            return nil
        }
        return errors.New("item not found")
    },
)

// Теперь работает для любого ID
service.GetItem("123") // Возвращает item, nil
service.GetItem("456") // Возвращает nil, error
```

### Советы для повышения качества тестов

1. **Используйте Helper-функции для повторяющейся логики**:

   ```go
   func TestMultipleCases(t *testing.T) {
       // Вспомогательная функция
       assertValidation := func(t *testing.T, input string, expectValid bool) {
           t.Helper() // Помечает как helper для лучшей трассировки
           
           isValid := Validate(input)
           assert.Equal(t, expectValid, isValid)
       }
       
       // Использование
       t.Run("Valid input", func(t *testing.T) {
           assertValidation(t, "valid_input", true)
       })
       
       t.Run("Invalid input", func(t *testing.T) {
           assertValidation(t, "", false)
       })
   }
   ```

2. **Создавайте функции-генераторы для тестовых данных**:

   ```go
   // Генератор пользователей с опциями
   func createTestUser(opts ...func(*User)) User {
       // Пользователь по умолчанию
       u := User{
           ID:    "test-1",
           Name:  "Test User",
           Email: "test@example.com",
           Age:   30,
       }
       
       // Применяем опции
       for _, opt := range opts {
           opt(&u)
       }
       
       return u
   }
   
   // Использование в тестах
   func TestUserValidation(t *testing.T) {
       // Валидный пользователь
       valid := createTestUser()
       assert.True(t, valid.IsValid())
       
       // Невалидный пользователь (без email)
       invalid := createTestUser(func(u *User) {
           u.Email = ""
       })
       assert.False(t, invalid.IsValid())
   }
   ```

3. **Документируйте тестовые случаи**:

   ```go
   // Хорошо документированный тест
   func TestCalculateTax(t *testing.T) {
       tests := []struct {
           name        string
           description string // Детальное описание сценария
           income      float64
           deductions  float64
           expected    float64
       }{
           {
               name:        "Standard tax calculation",
               description: "Стандартный расчет налога для дохода в средней категории с типичными вычетами",
               income:      50000,
               deductions:  5000,
               expected:    9000,
           },
           {
               name:        "Zero income",
               description: "Если доход нулевой, налог должен быть тоже нулевой",
               income:      0,
               deductions:  0,
               expected:    0,
           },
           // Другие тесты...
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Logf для документирования в выводе теста
               t.Logf("Test case: %s", tt.description)
               
               result := CalculateTax(tt.income, tt.deductions)
               assert.Equal(t, tt.expected, result)
           })
       }
   }
   ```

---

Это универсальное руководство охватывает все основные аспекты написания юнит-тестов в Go с использованием библиотеки Testify. Применяя эти практики, вы сможете создавать надежные, понятные и эффективные тесты для любых функций и компонентов вашего Go-приложения.
