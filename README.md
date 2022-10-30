# Тестовое задание на позицию стажёра-бэкендера

## Микросервис для работы с балансом пользователей

**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 

**При проектировании микросервиса были созданы следующие базы данных:**

CREATE TABLE `test`.`users` (`id` INT NOT NULL AUTO_INCREMENT , `name` VARCHAR(255) NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;

CREATE TABLE `test`.`balance` (`user_id` INT NOT NULL , `balance` DOUBLE NOT NULL ) ENGINE = InnoDB;

CREATE TABLE `test`.`orders` (`id` INT NOT NULL AUTO_INCREMENT , `name` VARCHAR(255) NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;

CREATE TABLE `test`.`reserve_money` (`id` INT NOT NULL AUTO_INCREMENT , `user_id` INT NOT NULL , `service_id` INT NOT NULL , `order_id` INT NOT NULL , `money` DOUBLE NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;

CREATE TABLE `test`.`report` (`id` INT NOT NULL AUTO_INCREMENT , `user_id` INT NOT NULL , `service_id` INT NOT NULL , `order_id` INT NOT NULL , `cash` DOUBLE NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;

*Далее заполнены следующими значениями:*

INSERT INTO `users`(`name`) VALUES ('Мария'),('Иван'),('Tata'),('Игорь'),('Олег Олегов');

INSERT INTO `services`(`name`, `price`) VALUES ('подстричь газон','1000'), ('ламинирование ресниц','1500'),('мастер на час','2000');

INSERT INTO `orders`(`name`) VALUES ('234'), ('12');


**Cервис:**

1. Предоставляет HTTP API с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: Golang.
3. Реляционная СУБД MySQL.
4. Использование docker и docker-compose для поднятия и развертывания dev-среды (не разобралась, поэтому не смогла выполнить эту задачу)

**HTTP API:**

1. Метод начисления средств на баланс. 
http://localhost:800/{id:uint64}/add/{money}
Вызывается метод ChangeMoneyToBalance из файла model.go. 
Данный метод:
- проверяет есть ли средства на балансе (если нет записи в таблицы, то создается),
- добавляет переданное кол-во денежных средств на баланс.

2. Метод резервирования средств с основного баланса на отдельном счете. 
http://localhost:800/user/{id:uint64}/reserve/{service_id:uint}/{order_id:uint}/{money}
Вызывается метод ReserveMoney из файла model.go. 
Данный метод:
- проверяет достаточно ли средств на счете пользователя (исключает возможности появления отрицательного баланса),
- делает запись в таблицу с зарезервированными счетами.

3. Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. 
http://localhost:800/confirm/{id:uint}
Вызывается метод DebitingFunds из файла model.go. 
Данный метод 
- удаляет запись из таблицы с зарезервированными денежными средствами, 
- уменьшает баланс владельца счета с зарезервированными ранее денежными средствами,
- добавляет запись в таблицу для отчета бухгалтерии

4. Метод получения баланса пользователя. 
http://localhost:800/user/{id:uint64}
Возвращает JSON с балансом пользователя.
Вызывается метод BalanceUser из файла model.go. 
Данный метод 
- по преданному значению id пользователя находит его текущий баланс в таблице

**Проблемы:**
При проектировании данного микросервиса были некоторые сложности, часть из которых я решила, часть просто не успела.
1) Синтаксис golang (впервые писала на данном яп, поэтому долго писала),
2) БД (не могла подключить БД 3 дня, оказалось что были ошибки в запросе на подключение, это даже больше относится к пункту 1),
3) Потрачено много дней на изучение ошибок, выводимых в терминале (Java более подробно все расписывает).
