# Транзакции в SQL: Полное руководство

ссылка на видеоматериал - [Транзакции и блокировки простым языком](https://www.youtube.com/watch?v=e9a4ESSHQ74)

## Что такое транзакции в SQL

Транзакция в SQL — это последовательность операций, которая выполняется как единая логическая единица работы. Проще говоря, это группа связанных операций с базой данных, которые либо все выполняются успешно, либо все отменяются, если хотя бы одна из них не может быть выполнена.

### Свойства транзакций (ACID)

Любая транзакция должна соответствовать принципам ACID:

1. **Атомарность (Atomicity)** — транзакция выполняется полностью или не выполняется вообще. Если какая-либо часть транзакции не может быть выполнена, вся транзакция откатывается (rollback), и база данных возвращается в состояние, в котором она была до начала транзакции.

2. **Согласованность (Consistency)** — транзакция переводит базу данных из одного согласованного состояния в другое. Все правила целостности данных, включая ограничения, каскады и триггеры, должны выполняться для каждой транзакции.

3. **Изолированность (Isolation)** — результаты транзакции не видны другим транзакциям до её завершения. Это предотвращает проблемы одновременного доступа, такие как "грязное чтение" (dirty read), "неповторяемое чтение" (non-repeatable read) и "фантомное чтение" (phantom read).

4. **Долговечность (Durability)** — после завершения транзакции (commit) изменения сохраняются в базе данных и не могут быть потеряны даже в случае сбоя системы.

### Синтаксис транзакций в SQL

```sql
BEGIN TRANSACTION; -- Начало транзакции

-- Операции с базой данных
INSERT INTO Accounts (AccountID, Balance) VALUES (1, 1000);
UPDATE Accounts SET Balance = Balance - 500 WHERE AccountID = 1;
INSERT INTO Transactions (AccountID, Amount) VALUES (1, -500);

COMMIT; -- Успешное завершение транзакции
-- или
ROLLBACK; -- Откат транзакции в случае ошибки
```

## Зачем нужны транзакции

1. **Обеспечение целостности данных** — транзакции гарантируют, что все связанные операции либо выполняются вместе, либо не выполняются вообще, что критически важно для сохранения согласованности данных.

2. **Восстановление после сбоев** — если система выходит из строя во время выполнения транзакции, механизм транзакций обеспечивает возможность восстановления базы данных до согласованного состояния.

3. **Параллельная обработка** — транзакции позволяют нескольким пользователям одновременно работать с базой данных, не мешая друг другу и не нарушая целостность данных.

4. **Управление конкурентным доступом** — изолированность транзакций предотвращает проблемы одновременного доступа к данным.

## Где используются транзакции

1. **Банковские системы** — классический пример: перевод денег между счетами должен быть атомарной операцией; сумма должна быть снята с одного счета и зачислена на другой, либо операция должна быть полностью отменена.

2. **Системы бронирования** — например, при бронировании авиабилетов нужно одновременно зарезервировать место и обновить информацию о наличии мест.

3. **Электронная коммерция** — обработка заказов включает обновление запасов товаров, информации о клиенте и финансовых данных.

4. **ERP-системы** — сложные бизнес-процессы, где несколько таблиц должны быть обновлены согласованно.

5. **Системы управления контентом (CMS)** — создание и публикация контента может затрагивать несколько таблиц и связей.

## Уровни изоляции транзакций

В SQL определены четыре уровня изоляции транзакций, которые определяют, насколько изолированы транзакции друг от друга:

1. **READ UNCOMMITTED** — самый низкий уровень изоляции, допускает "грязное чтение".
2. **READ COMMITTED** — предотвращает "грязное чтение", но допускает "неповторяемое чтение".
3. **REPEATABLE READ** — предотвращает "грязное" и "неповторяемое" чтение, но допускает "фантомное чтение".
4. **SERIALIZABLE** — самый высокий уровень изоляции, предотвращает все проблемы, но может снижать производительность.

## Проблемы конкурентного доступа

1. **Грязное чтение (Dirty Read)** — когда транзакция читает данные, которые были изменены, но еще не зафиксированы другой транзакцией.

2. **Неповторяемое чтение (Non-repeatable Read)** — когда транзакция повторно читает те же данные, но обнаруживает, что они были изменены другой транзакцией.

3. **Фантомное чтение (Phantom Read)** — когда транзакция повторно выполняет запрос, возвращающий набор строк, но обнаруживает, что этот набор изменился из-за другой транзакции.

## Примеры проблем конкурентного доступа

### 1. Грязное чтение (Dirty Read)

**Пример с банковским счетом:**

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN TRANSACTION; | |
| T2 | UPDATE Accounts SET Balance = Balance + 1000 WHERE AccountID = 123; | |
| T3 | | BEGIN TRANSACTION; |
| T4 | | SELECT Balance FROM Accounts WHERE AccountID = 123; <br>/ ***Видит увеличенный баланс 2000 руб., хотя Транзакция 1 еще не завершена*** / |
| T5 | ROLLBACK; / ***Отмена транзакции из-за ошибки*** / | |
| T6 | | /*Транзакция 2 продолжает работать с неверными данными*/ |
| T7 | | COMMIT; |

**Проблема:** Транзакция 2 видит и использует данные, которые фактически никогда не были зафиксированы в базе данных, что приводит к неверным расчетам.

### 2. Неповторяемое чтение (Non-repeatable Read)

**Пример с управлением товарами:**

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN TRANSACTION; | |
| T2 | SELECT Price FROM Products WHERE ProductID = 456; <br>/ ***Цена = 100 руб.*** / | |
| T3 | | BEGIN TRANSACTION; |
| T4 | | UPDATE Products SET Price = 150 WHERE ProductID = 456; |
| T5 | | COMMIT; |
| T6 | SELECT Price FROM Products WHERE ProductID = 456; <br>/ ***Цена уже = 150 руб.*** / | |
| T7 | / ***Транзакция 1 обрабатывает данные, но сталкивается с неконсистентностью, т.к. цена изменилась*** / | |
| T8 | COMMIT; | |

**Проблема:** Транзакция 1 выполнила один и тот же запрос дважды и получила разные результаты, что может нарушить бизнес-логику приложения. Например, при расчете стоимости заказа могут возникнуть несоответствия.

### 3. Фантомное чтение (Phantom Read)

**Пример с системой бронирования:**

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN TRANSACTION; | |
| T2 | SELECT COUNT(*) FROM Rooms WHERE HotelID = 789 AND IsAvailable = 1; <br>/ ***Количество = 5 свободных комнат*** / | |
| T3 | | BEGIN TRANSACTION; |
| T4 | | INSERT INTO Rooms (RoomID, HotelID, IsAvailable) VALUES (1001, 789, 1); <br>/ ***Добавили новую комнату*** / |
| T5 | | COMMIT; |
| T6 | SELECT COUNT(*) FROM Rooms WHERE HotelID = 789 AND IsAvailable = 1; <br>/ ***Теперь количество = 6 свободных комнат*** / | |
| T7 | / ***Транзакция 1 использует противоречивые данные для расчетов*** / | |
| T8 | COMMIT; | |

**Проблема:** Транзакция 1 видит "фантомные" строки, которых не было при первом выполнении запроса. Это может привести к ошибкам в логике приложения, особенно в системах, где критично точное количество записей.

### 4. Потерянное обновление (Lost Update)

Хотя это не входит в основные три проблемы, но часто встречается на практике:

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN TRANSACTION; | |
| T2 | SELECT Quantity FROM Inventory WHERE ProductID = 321; <br>/ ***Количество = 100*** / | |
| T3 | | BEGIN TRANSACTION; |
| T4 | | SELECT Quantity FROM Inventory WHERE ProductID = 321; <br>/ ***Количество = 100*** / |
| T5 | UPDATE Inventory SET Quantity = Quantity - 20 WHERE ProductID = 321; | |
| T6 | | UPDATE Inventory SET Quantity = Quantity - 30 WHERE ProductID = 321; |
| T7 | COMMIT; | |
| T8 | | COMMIT; |
| T9 | / ***Итоговое количество = 70, хотя должно быть 50*** / | |

**Проблема:** Обновление, выполненное Транзакцией 1, "теряется", потому что Транзакция 2 перезаписывает его своим обновлением, не учитывая изменения, внесенные Транзакцией 1.

## Сравнение уровней изоляции

| Уровень изоляции | Грязное чтение | Неповторяемое чтение | Фантомное чтение | Потерянное обновление | Производительность |
|------------------|----------------|----------------------|-------------------|------------------------|-------------------|
| READ UNCOMMITTED | ❌ Возможно | ❌ Возможно | ❌ Возможно | ❌ Возможно | 🟢 Высокая |
| READ COMMITTED | ✅ Предотвращает | ❌ Возможно | ❌ Возможно | ❌ Возможно | 🟢 Высокая |
| REPEATABLE READ | ✅ Предотвращает | ✅ Предотвращает | ❌ Возможно* | ✅ Предотвращает | 🟡 Средняя |
| SERIALIZABLE | ✅ Предотвращает | ✅ Предотвращает | ✅ Предотвращает | ✅ Предотвращает | 🔴 Низкая |

\* *Примечание: В некоторых СУБД (например, MySQL с InnoDB) уровень REPEATABLE READ также предотвращает фантомное чтение.*

### Дополнительная информация по уровням изоляции

1. **READ UNCOMMITTED**:
   - Самый низкий уровень изоляции
   - Транзакции могут видеть незафиксированные изменения других транзакций
   - Используется редко, в основном для операций чтения нечувствительных данных

2. **READ COMMITTED**:
   - Уровень изоляции по умолчанию во многих СУБД (например, PostgreSQL, Oracle)
   - Транзакции видят только зафиксированные изменения других транзакций
   - Хороший баланс между производительностью и согласованностью для большинства приложений

3. **REPEATABLE READ**:
   - Уровень изоляции по умолчанию в MySQL
   - Гарантирует, что повторные чтения в рамках одной транзакции возвращают те же результаты
   - Обычно реализуется с помощью блокировок чтения

4. **SERIALIZABLE**:
   - Самый высокий уровень изоляции
   - Транзакции выполняются так, как если бы они выполнялись последовательно
   - Обычно используется для финансовых операций и других критически важных транзакций

## Способы реализации контроля конкурентного доступа

1. **Блокировки (Locking)**:
   - Пессимистический подход
   - Блокирует ресурсы при доступе к ним
   - Виды: разделяемые (для чтения) и исключительные (для записи)

2. **Многоверсионность (MVCC - Multi-Version Concurrency Control)**:
   - Оптимистический подход
   - Создает "снимки" данных для каждой транзакции
   - Используется в PostgreSQL, Oracle, MySQL (InnoDB)

При выборе уровня изоляции всегда приходится балансировать между согласованностью данных и производительностью системы. Чем выше уровень изоляции, тем больше гарантий согласованности, но ниже параллелизм и выше вероятность блокировок и взаимоблокировок.

## Примеры использования уровней изоляции транзакций

Рассмотрим практические примеры использования различных уровней изоляции на реальных сценариях.

### READ UNCOMMITTED

Этот уровень редко используется в производственных системах из-за отсутствия гарантий, но может быть полезен в некоторых случаях.

```sql
-- Сеанс 1
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
BEGIN;
-- Получение предварительной статистики, где точность не критична
SELECT COUNT(*) FROM orders WHERE status = 'processing';
-- Другие операции чтения
COMMIT;
```

**Пример с отчетностью**:
Предположим, нам нужно получить примерное количество активных пользователей для необязательного отчета.

| Время | Сеанс 1 (READ UNCOMMITTED) | Сеанс 2 |
|-------|---------------------------|---------|
| T1 | BEGIN; | |
| T2 | | BEGIN; |
| T3 | | INSERT INTO active_users VALUES (101, 'user101', NOW()); |
| T4 | SELECT COUNT(*) FROM active_users; <br>/ ***возвращает 1001, включая незафиксированного пользователя*** / | |
| T5 | | ROLLBACK; / ***пользователь фактически не был добавлен*** / |
| T6 | / ***отчет содержит неверные данные*** / | |
| T7 | COMMIT; | |

### READ COMMITTED

Наиболее часто используемый уровень изоляции. Предотвращает грязное чтение, но уязвим для неповторяемого чтения.

```sql
-- Сеанс 1
SET TRANSACTION ISOLATION LEVEL READ COMMITTED; -- В PostgreSQL это уровень по умолчанию
BEGIN;
-- Получение информации о продукте для отображения клиенту
SELECT price FROM products WHERE product_id = 123;
-- Имитация времени, пока клиент принимает решение о покупке
-- ... (пауза)
-- Повторная проверка цены перед покупкой
SELECT price FROM products WHERE product_id = 123; -- цена может отличаться!
-- Создание заказа
INSERT INTO orders (product_id, customer_id, price) VALUES (123, 456, :current_price);
COMMIT;
```

**Пример с обработкой заказов**:

| Время | Сеанс 1 (READ COMMITTED) | Сеанс 2 |
|-------|--------------------------|---------|
| T1 | BEGIN; | |
| T2 | SELECT price FROM products WHERE id = 101; <br>/ ***возвращает 199.99*** / | |
| T3 | | BEGIN; |
| T4 | | UPDATE products SET price = 249.99 WHERE id = 101; |
| T5 | | COMMIT; |
| T6 | SELECT price FROM products WHERE id = 101; <br>/ ***теперь возвращает 249.99*** / | |
| T7 | / ***клиент видит изменение цены в процессе оформления заказа*** / | |
| T8 | COMMIT; | |

### REPEATABLE READ

Обеспечивает более сильные гарантии - предотвращает неповторяемое чтение.

```sql
-- Сеанс 1
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
BEGIN;
-- Генерация отчета за определенный период
SELECT SUM(amount) FROM transactions WHERE date BETWEEN '2023-01-01' AND '2023-01-31';
-- Обработка данных отчета
-- ...
-- Повторный запрос даст тот же результат, даже если другие транзакции добавили новые записи
SELECT SUM(amount) FROM transactions WHERE date BETWEEN '2023-01-01' AND '2023-01-31';
COMMIT;
```

**Пример с формированием отчета**:

| Время | Сеанс 1 (REPEATABLE READ) | Сеанс 2 |
|-------|---------------------------|---------|
| T1 | BEGIN; | |
| T2 | SELECT SUM(revenue) FROM sales WHERE date = CURRENT_DATE; <br>/ ***возвращает 15000*** / | |
| T3 | | BEGIN; |
| T4 | | INSERT INTO sales (id, revenue, date) VALUES (999, 5000, CURRENT_DATE); |
| T5 | | COMMIT; |
| T6 | SELECT SUM(revenue) FROM sales WHERE date = CURRENT_DATE; <br>/ ***все равно возвращает 15000*** / | |
| T7 | COMMIT; | |

### SERIALIZABLE

Самый строгий уровень изоляции, предотвращающий все аномалии.

```sql
-- Сеанс 1
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
BEGIN;
-- Проверка наличия свободного места в самолете
SELECT COUNT(*) FROM airplane_seats WHERE flight_id = 789 AND is_occupied = FALSE;
-- Если место доступно, бронируем его
UPDATE airplane_seats SET is_occupied = TRUE, customer_id = 1001 
WHERE seat_id = (
    SELECT seat_id FROM airplane_seats 
    WHERE flight_id = 789 AND is_occupied = FALSE 
    LIMIT 1
);
COMMIT;
```

**Пример с бронированием билетов**:

| Время | Сеанс 1 (SERIALIZABLE) | Сеанс 2 (SERIALIZABLE) |
|-------|------------------------|------------------------|
| T1 | BEGIN; | |
| T2 | SELECT COUNT(*) FROM seats WHERE flight_id = 101 AND is_booked = FALSE; <br>/ ***возвращает 1 свободное место*** / | |
| T3 | | BEGIN; |
| T4 | | SELECT COUNT(*) FROM seats WHERE flight_id = 101 AND is_booked = FALSE; <br>/ ***тоже возвращает 1 свободное место*** / |
| T5 | UPDATE seats SET is_booked = TRUE, customer_id = 555 WHERE id = 42; | |
| T6 | | UPDATE seats SET is_booked = TRUE, customer_id = 777 WHERE id = 42; <br>/ ***в PostgreSQL эта транзакция будет ждать завершения первой*** / |
| T7 | COMMIT; | |
| T8 | | / ***Попытка обновления завершается ошибкой: "could not serialize access due to concurrent update"*** / |
| T9 | | ROLLBACK; |

## Использование FOR UPDATE для блокировки строк

Конструкция `FOR UPDATE` позволяет реализовать эксклюзивную блокировку выбранных строк независимо от уровня изоляции, предотвращая их изменение другими транзакциями до завершения текущей.

### Пример 1: Система управления дежурствами врачей

```sql
-- Сеанс 1
BEGIN;
-- Блокируем строки для обновления
SELECT * FROM doctors WHERE shift_id = 228 FOR UPDATE;
-- Проверяем текущее состояние
SELECT COUNT(*) FROM doctors WHERE shift_id = 228 AND on_call = TRUE;
-- Выполняем обновление
UPDATE doctors SET on_call = FALSE WHERE name = 'Alice' AND shift_id = 228;
COMMIT;
```

| Время | Сеанс 1 | Сеанс 2 |
|-------|---------|---------|
| T1 | BEGIN; | |
| T2 | SELECT *FROM doctors WHERE shift_id = 228 FOR UPDATE; <br>/ ***Блокирует строки - Alice и Bob, оба on_call = TRUE*** / | |
| T3 | | BEGIN; |
| T4 | | SELECT *FROM doctors WHERE shift_id = 228 FOR UPDATE; <br>/ ***Запрос блокируется и ожидает разблокировки строк*** / |
| T5 | UPDATE doctors SET on_call = FALSE WHERE name = 'Alice' AND shift_id = 228; | |
| T6 | COMMIT; | |
| T7 | | / ***Запрос разблокируется и возвращает обновленные данные*** / <br>/ ***Alice (on_call = FALSE), Bob (on_call = TRUE)*** / |
| T8 | | / ***Теперь Сеанс 2 может выполнять свои операции*** / |

### Пример 2: Управление запасами товаров

Этот пример показывает, как предотвратить потерянное обновление при работе с запасами товаров:

```sql
-- Сеанс 1: Обработка заказа
BEGIN;
-- Блокируем запись о товаре
SELECT * FROM inventory WHERE product_id = 321 FOR UPDATE;
-- Проверяем наличие достаточного количества товара
-- Обновляем запасы
UPDATE inventory SET quantity = quantity - 20 WHERE product_id = 321;
COMMIT;
```

| Время | Сеанс 1 | Сеанс 2 |
|-------|---------|---------|
| T1 | BEGIN; | |
| T2 | SELECT *FROM inventory WHERE product_id = 321 FOR UPDATE; <br>/ *** Блокирует строку, quantity = 100*** / | |
| T3 | | BEGIN; |
| T4 | | SELECT *FROM inventory WHERE product_id = 321 FOR UPDATE; <br>/ ***Ожидает разблокировки строки*** / |
| T5 | UPDATE inventory SET quantity = quantity - 20 WHERE product_id = 321; | |
| T6 | COMMIT; | |
| T7 | | / ***Теперь запрос получает обновленные данные, quantity = 80*** / |
| T8 | | UPDATE inventory SET quantity = quantity - 30 WHERE product_id = 321; <br>/ ***Обновляет с учетом предыдущего изменения*** / |
| T9 | | COMMIT; <br>/ ***Итоговое количество = 50, корректно*** / |

### Пример 3: FOR UPDATE SKIP LOCKED для обработки очередей

В некоторых СУБД (PostgreSQL, Oracle) существует модификация `FOR UPDATE SKIP LOCKED`, которая пропускает уже заблокированные строки вместо ожидания. Это полезно для систем обработки очередей:

```sql
-- Обработчик заданий
BEGIN;
-- Получаем и блокируем первое доступное задание
SELECT * FROM job_queue 
WHERE status = 'pending' 
ORDER BY priority DESC, created_at 
FOR UPDATE SKIP LOCKED
LIMIT 1;

-- Помечаем задание как обрабатываемое
UPDATE job_queue SET status = 'processing', started_at = NOW() WHERE id = :job_id;
COMMIT;
```

Это позволяет нескольким обработчикам параллельно получать разные задания из очереди без блокировок друг друга.

### Пример 4: FOR UPDATE с NOWAIT

Еще одна полезная модификация - `FOR UPDATE NOWAIT`, которая вызывает ошибку немедленно, если строки заблокированы, вместо ожидания:

```sql
BEGIN;
-- Пытаемся заблокировать, но не ждем, если уже заблокировано
SELECT * FROM reservations 
WHERE resource_id = 555 AND date = CURRENT_DATE 
FOR UPDATE NOWAIT;

-- Если дошли до этой точки, значит блокировка получена
UPDATE reservations SET status = 'confirmed' WHERE id = :reservation_id;
COMMIT;
```

Этот подход полезен для интерактивных систем, где ожидание нежелательно, и лучше сразу сообщить пользователю о невозможности выполнения операции.

## Расширенные примеры аномалий конкурентного доступа

Для лучшего понимания проблем конкурентного доступа, рассмотрим еще несколько реальных примеров:

### Грязное чтение в системе складского учета

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN; | |
| T2 | INSERT INTO shipments (id, product_id, quantity, status) VALUES (501, 123, 1000, 'pending'); | |
| T3 | | BEGIN; |
| T4 | | SELECT SUM(quantity) FROM shipments WHERE product_id = 123; <br>/ ***Включает незафиксированную отгрузку в 1000 единиц*** / |
| T5 | | UPDATE inventory SET quantity_reserved = quantity_reserved + 1000 WHERE product_id = 123; <br>/ ***Резервирует товар, основываясь на недостоверных данных*** / |
| T6 | ROLLBACK; / ***Отмена отгрузки из-за ошибки*** / | |
| T7 | | COMMIT; <br>/ ***Привело к некорректному резервированию товара на складе*** / |

### Неповторяемое чтение в системе учета рабочего времени

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN; | |
| T2 | SELECT SUM(hours) FROM time_logs WHERE employee_id = 123 AND week = 25; <br>/ ***Возвращает 38 часов*** / | |
| T3 | / ***Сотрудник HR проверяет, нужно ли оплачивать сверхурочные (>40 часов)*** / | |
| T4 | | BEGIN; |
| T5 | | INSERT INTO time_logs (employee_id, day, hours, week) VALUES (123, 'Friday', 8, 25); <br>/ ***Сотрудник добавил последний день учета*** / |
| T6 | | COMMIT; |
| T7 | SELECT SUM(hours) FROM time_logs WHERE employee_id = 123 AND week = 25; <br>/ ***Теперь возвращает 46 часов*** / | |
| T8 | / ***HR принимает решение о сверхурочных, но данные изменились во время анализа*** / | |
| T9 | COMMIT; | |

### Фантомное чтение в системе бронирования авиабилетов

| Время | Транзакция 1 | Транзакция 2 |
|-------|--------------|--------------|
| T1 | BEGIN; | |
| T2 | SELECT flight_id, COUNT(*) as free_seats FROM seats WHERE flight_id = 555 AND is_booked = FALSE GROUP BY flight_id; <br>/ ***Возвращает flight_id=555, free_seats=3*** / | |
| T3 | / ***Агент сообщает клиенту, что для группы из 3 человек есть места*** / | |
| T4 | | BEGIN; |
| T5 | | INSERT INTO seats (id, flight_id, seat_number, is_booked) VALUES (9876, 555, '24D', TRUE); <br>/ ***Добавление нового места в систему, но уже забронированного*** / |
| T6 | | COMMIT; |
| T7 | SELECT COUNT(*) FROM seats WHERE flight_id = 555 AND is_booked = FALSE; <br>/ ***Все еще показывает 3 свободных места*** / | |
| T8 | / ***Агент приступает к бронированию 3 мест*** / | |
| T9 | SELECT *FROM seats WHERE flight_id = 555 AND is_booked = FALSE LIMIT 3 FOR UPDATE; <br>/ ***Возвращает 3 места, но при попытке бронирования последнего возникнет ошибка*** / | |
| T10 | COMMIT; | |

### Потерянное обновление в многопользовательском редакторе документов

| Время | Пользователь 1 | Пользователь 2 |
|-------|----------------|----------------|
| T1 | BEGIN; | |
| T2 | SELECT content FROM documents WHERE id = 42; <br>/ ***Получает текст документа*** / | |
| T3 | / ***Редактирует документ локально*** / | |
| T4 | | BEGIN; |
| T5 | | SELECT content FROM documents WHERE id = 42; <br>/ ***Тоже получает текст документа*** / |
| T6 | | / ***Также редактирует документ локально*** / |
| T7 | UPDATE documents SET content = 'Версия пользователя 1', updated_at = NOW() WHERE id = 42; | |
| T8 | COMMIT; | |
| T9 | | UPDATE documents SET content = 'Версия пользователя 2', updated_at = NOW() WHERE id = 42; <br>/ ***Перезаписывает изменения первого пользовател*** / |
| T10 | | COMMIT; |
| T11 | / ***Изменения пользователя 1 потеряны*** / | |

## Решение проблем с помощью FOR UPDATE

Многие из описанных выше проблем можно решить с помощью блокировки строк через `FOR UPDATE`. Рассмотрим модифицированный пример системы бронирования авиабилетов:

```sql
-- Правильная реализаяция бронирования мест
BEGIN;
-- Блокируем строки для обновления, получаем точное число свободных мест
SELECT COUNT(*) FROM seats 
WHERE flight_id = 555 AND is_booked = FALSE 
FOR UPDATE;

-- Если мест достаточно, бронируем их
UPDATE seats 
SET is_booked = TRUE, customer_id = 123, booking_time = NOW() 
WHERE id IN (
    SELECT id FROM seats 
    WHERE flight_id = 555 AND is_booked = FALSE 
    LIMIT 3
    FOR UPDATE
);
COMMIT;
```

## Заключение

Транзакции в SQL — это фундаментальный механизм обеспечения целостности и согласованности данных в многопользовательской среде. Правильное использование транзакций, выбор подходящего уровня изоляции и применение блокировок с помощью `FOR UPDATE` критически важны для создания надежных баз данных и предотвращения аномалий конкурентного доступа.

На собеседовании важно продемонстрировать понимание не только базовых принципов ACID, но и проблем конкурентного доступа, компромиссов между уровнями изоляции, а также практических механизмов реализации транзакционной модели в конкретных СУБД.

При проектировании системы всегда следует учитывать требования к целостности данных, производительности и конкурентности, выбирая оптимальный баланс между ними с помощью средств управления транзакциями.
