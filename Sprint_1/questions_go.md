# Вопросы к собесам - Сети 

## 🔹 **Основные различия между TCP и UDP**  

**TCP** (Transmission Control Protocol) — надёжный, ориентированный на соединение протокол с гарантированной доставкой данных. 

**UDP** (User Datagram Protocol) — более быстрый, но ненадёжный протокол без установки соединения и подтверждения доставки.  

---

| Характеристика  | TCP | UDP |
|---------------|----|----|
| **Тип протокола** | Ориентированный на соединение (нужно установить соединение перед передачей) | Без установления соединения (отправка без предварительных договорённостей) |
| **Гарантия доставки** | Да (использует подтверждения и повторную отправку) | Нет (пакеты могут потеряться или прийти не в том порядке) |
| **Скорость** | Медленнее из-за контроля ошибок и подтверждений | Быстрее, так как нет контроля доставки |
| **Последовательность пакетов** | Гарантированная (данные приходят в том порядке, в котором отправлены) | Нет гарантии (могут прийти в любом порядке) |
| **Контроль ошибок** | Да (использует контрольную сумму, механизмы проверки потерь) | Минимальный (только контрольная сумма, но нет повторной передачи) |
| **Использование** | Надёжные соединения: HTTP, HTTPS, FTP, SSH | Потоковые передачи: VoIP, видео, онлайн-игры |

---

### 🔹 **Когда использовать TCP?**  

✅ **Если важна надёжность**:  

- Передача файлов (FTP, HTTP/S, SFTP)  
- Электронная почта (SMTP, IMAP, POP3)  
- Доступ к серверам (SSH, Telnet)  

✅ **Если важен порядок пакетов**:  

- Передача документов и баз данных  

✅ **Если важно контролировать ошибки**:  

- Банковские системы, онлайн-оплаты  

---

## 🔹 **Когда использовать UDP?**  

✅ **Если важна скорость, а не 100% надёжность**:  

- Видеозвонки и голосовые чаты (VoIP, Zoom, Skype)  
- Онлайн-игры (CS:GO, Dota 2, Fortnite)  
- Потоковое видео (YouTube, Twitch, IPTV)  

✅ **Если потеря нескольких пакетов не критична**:  

- DNS-запросы  
- Системы вещания (broadcast, multicast)  

---

### 🔹 **Простая аналогия**  

- **TCP** — это как почта с уведомлением о доставке (если письмо потерялось, его отправляют заново).  
- **UDP** — это как обычная открытка (если потерялась — ничего страшного, просто отправили дальше).  

---

## HTTP

**HTTP (HyperText Transfer Protocol)** — это протокол передачи данных, используемый для загрузки веб-страниц, API-запросов и взаимодействия между клиентом (браузером) и сервером.

Работает по модели "клиент-сервер":

- Клиент отправляет HTTP-запрос (например, запрашивает веб-страницу).
- Сервер обрабатывает запрос и отправляет HTTP-ответ (HTML, JSON, изображение и т. д.).

---

### 🔹 **Где используется HTTP?**  

✅ Открытие веб-сайтов (браузер).  
✅ REST API и веб-сервисы.  
✅ Запросы к серверам (мобильные приложения, боты).  
✅ Интернет-магазины, соцсети, облачные сервисы.  

---

### 🔹 **Виды HTTP-запросов (Методы HTTP)**  

| Метод  | Описание | Пример использования |
|--------|---------|----------------------|
| **GET** | Запрашивает данные | Открытие страницы (`/index.html`) |
| **POST** | Отправляет данные | Авторизация, регистрация |
| **PUT** | Обновляет данные | Изменение профиля |
| **DELETE** | Удаляет ресурс | Удаление аккаунта |
| **PATCH** | Частичное обновление | Обновление одного поля в базе |
| **HEAD** | Запрос без тела (только заголовки) | Проверка доступности ресурса |

---

### 🔹 **HTTP и HTTPS – в чём разница?**  

**HTTPS** (HyperText Transfer Protocol Secure) → использует **шифрование TLS/SSL**, защищает данные.  

📌 **Пример:**  
🔴 **<http://example.com>** → Небезопасное соединение.  
🟢 **<https://example.com>** → Защищённое соединение.  

---

### 🔹 **Структура HTTP-запроса**  

HTTP-запрос состоит из трёх частей:  

1. **Стартовая строка** — определяет метод, URL и версию протокола.  
2. **Заголовки** — содержат дополнительную информацию (типы данных, кодировки, авторизация).  
3. **Тело (Body)** — присутствует в методах, которые отправляют данные (`POST`, `PUT` и т. д.).  

#### **Пример HTTP-запроса (GET)**

```http
GET /index.html HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)
Accept: text/html,application/xhtml+xml
Connection: keep-alive
```

🔹 **Разбор запроса:**  

- `GET /index.html HTTP/1.1` → **Метод `GET`**, запрашиваемый ресурс `/index.html`, версия протокола `HTTP/1.1`.  
- `Host: example.com` → Указывает, к какому серверу обращаемся.  
- `User-Agent: Mozilla/5.0` → Идентифицирует браузер или клиента.  
- `Accept: text/html` → Какие форматы данных клиент ожидает.  
- `Connection: keep-alive` → Просьба не закрывать соединение сразу.  

---

#### **Пример HTTP-запроса (POST) с телом**  

```http
POST /api/login HTTP/1.1
Host: example.com
Content-Type: application/json
Content-Length: 32

{"username": "user", "password": "123"}
```

🔹 **Разбор:**  

- `POST /api/login HTTP/1.1` → Отправка данных на сервер.  
- `Content-Type: application/json` → Указываем формат тела (JSON).  
- `Content-Length: 32` → Длина тела запроса в байтах.  
- Тело содержит JSON с логином и паролем.  

---

### 🔹 **Структура HTTP-ответа**  

HTTP-ответ также имеет три части:  

1. **Статусная строка** — версия протокола, статус-код, описание.  
2. **Заголовки** — информация о сервере, формате ответа, кодировке.  
3. **Тело (Body)** — данные, которые сервер отправляет клиенту.  

#### **Пример HTTP-ответа (успешный запрос)**

```http
HTTP/1.1 200 OK
Date: Tue, 30 Jan 2025 12:34:56 GMT
Server: Apache/2.4.41 (Ubuntu)
Content-Type: text/html; charset=UTF-8
Content-Length: 1234

<html>
<head><title>Главная</title></head>
<body><h1>Добро пожаловать!</h1></body>
</html>
```

🔹 **Разбор:**  

- `HTTP/1.1 200 OK` → Протокол `HTTP/1.1`, код `200` (успешно).  
- `Date: Tue, 30 Jan 2025 12:34:56 GMT` → Время ответа.  
- `Server: Apache/2.4.41 (Ubuntu)` → Какой сервер отвечает.  
- `Content-Type: text/html` → Формат данных (HTML).  
- `Content-Length: 1234` → Длина тела ответа.  
- **Тело** → HTML-код страницы.  

---

### 🔹 **Статус-коды HTTP** 

- `200 OK` → Успешный запрос.  
- `301 Moved Permanently` → Страница перемещена (редирект).  
- `400 Bad Request` → Ошибочный запрос (неверный формат данных).  
- `401 Unauthorized` → Требуется авторизация.  
- `403 Forbidden` → Доступ запрещён.  
- `404 Not Found` → Ресурс не найден.  
- `500 Internal Server Error` → Ошибка сервера.  

---

### 🔹 **Вывод:**  

- **HTTP-запрос** → содержит метод (`GET`, `POST` и т. д.), заголовки и (иногда) тело.  
- **HTTP-ответ** → состоит из статусного кода, заголовков и (иногда) тела.  
- **Для работы с API** → часто используются JSON (`application/json`) и коды состояния.  
