# from CodeWars

## return "(123) 456-7890"

```go

CreatePhoneNumber([10]uint{1,2,3,4,5,6,7,8,9,0})  // returns "(123) 456-7890"

//

func CreatePhoneNumber(numbers [10]uint) string {
  var result string 
  simbols := "() -"
  for i, v := range(numbers) {
    switch i {
      case 0:
        result = string(simbols[0]) + toString(v)
      case 2:
        result = result + toString(v) + string(simbols[1:3])
      case 5:
        result = result + toString(v) + string(simbols[3])
      default:
        result = result + toString(v)
    }
  }
  return result 
}

func toString(n uint) string {
    return strconv.FormatUint(uint64(n), 10)
}

```

best practices

```go
package kata

import "fmt"

func CreatePhoneNumber(n [10]uint) string {
  return fmt.Sprintf(
    "(%d%d%d) %d%d%d-%d%d%d%d",
    n[0], n[1], n[2], n[3], n[4], n[5], n[6], n[7], n[8], n[9]
  )
}
```

---

## GPT - cкобки

Напишите функцию IsBalanced(s string) bool, которая проверяет, является ли строка s сбалансированной по скобкам.

Условия:

Строка считается сбалансированной, если:

- У нее есть парные открывающие и закрывающие скобки ((), {}, []).
- Скобки правильно вложены (например, "{[()]}" — верно, а "{[(])}" — нет).
- В строке могут быть другие символы, но они игнорируются.

```go
fmt.Println(IsBalanced("{[()]}"))     // true  
fmt.Println(IsBalanced("{[(])}"))     // false  
fmt.Println(IsBalanced("{[()]()}"))   // true  
fmt.Println(IsBalanced("((()))"))     // true  
fmt.Println(IsBalanced("[{]}"))       // false  
fmt.Println(IsBalanced("{Hello}"))    // true  
fmt.Println(IsBalanced("{[(])"))      // false  
```

Ограничения:

- Можно использовать map или stack (слайс).
- Запрещено использовать готовые библиотеки для работы со скобками.

Решение:

```go
func IsBalanced(s string) bool {
 stack := []rune{}
 brackets := map[rune]rune{
  ')': '(', 
  '}': '{', 
  ']': '[',
 }

 for _, ch := range s {
  switch ch {
  case '(', '{', '[':
   stack = append(stack, ch) // Добавляем открывающую скобку в стек
  case ')', '}', ']':
   if len(stack) == 0 || stack[len(stack)-1] != brackets[ch] {
    return false // Несоответствие закрывающей скобки
   }
   stack = stack[:len(stack)-1] // Удаляем последнюю открывающую скобку
  }
 }

 return len(stack) == 0 // Если стек пуст, скобки сбалансированы
}
```

---
