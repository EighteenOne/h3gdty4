# Application design

В коде представлен прототип сервиса бронирования номеров в отелях,
в котором реализована возможность забронировать свободный номер в отеле.

Сервис будет развиваться, например:

- появится отправка письма-подтверждения о бронировании
- появятся скидки, промокоды, программы лояльности
- появится возможность бронирования нескольких номеров

## Задание

Провести рефакторинг структуры и кода приложения, исправить существующие
проблемы в логике. Персистентное хранение реализовывать не требуется,
все данные храним в памяти сервиса.

В результате выполнения задания ожидается структурированный код сервиса,
с корректно работающей логикой сценариев бронирования номеров в отелях.

Чеклист:

- код реорганизован и выделены слои
- выделены абстракций и интерфейсы
- техническе и логические ошибки исправлены

Ограничения:

- ожидаем реализацию, которая управляет состоянием в памяти приложения,
 но которую легко заменить на внешнее хранилище
- если у тебя есть опыт с Go: для решения надо использовать только
 стандартную библиотеку Go + роутер (например chi)
- если у тебя нет опыта с Go: можно реализовать решение на своем
 любимом стеке технологий

## Что будет на встрече

На встрече ожидаем что ты продемонстрируешь экран и презентуешь свое решение:
расскажешь какими проблемами обладает исходный код и как они решены в твоем варианте.
Мы будем задавать вопросы о том почему было решено разделить ответственность между
компонентами тем или иным образом, какими принципами ты при этом руководствуешься.
Спросим что будет если продакт решит добавить какую-то новую фичу — как она ляжет
на предложенную тобой структуру. Также можем поговорить и о более технических вещах:
о значениях и указателях, многопоточности, интерфейсах, каналах.

## Например

```sh
go run main.go
```

```sh
curl --location --request POST 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data-raw '{
    "hotel_id": "reddison",
    "room_id": "lux",
    "email": "guest@mail.ru",
    "from": "2024-01-02T00:00:00Z",
    "to": "2024-01-04T00:00:00Z"
}'
```



var localTestAvailability = map[string]domain.RoomAvailability{
	"reddison_lux_2024-01-01": {"reddison", "lux", date(2024, 1, 1), 1},
	"reddison_lux_2024-01-02": {"reddison", "lux", date(2024, 1, 2), 1},
	"reddison_lux_2024-01-03": {"reddison", "lux", date(2024, 1, 3), 2},
	"reddison_lux_2024-01-04": {"reddison", "lux", date(2024, 1, 4), 1},
	"reddison_lux_2024-01-05": {"reddison", "lux", date(2024, 1, 5), 0},
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}