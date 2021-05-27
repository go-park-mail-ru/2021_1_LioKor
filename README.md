# LioKor Mail Backend

**Почтовый сервис с веб-интерфейсом**

### Команда:
* [Артем Королев](https://github.com/KoroLion)
* [Алтана Бадмацыренова](https://github.com/altanab)
* [Сергей Тяпкин](https://github.com/SergTyapkin)

### Менторы:
* [Джахонгир Тулфоров](https://github.com/bin-umar) (фронтенд)
* [Владимир Северов](https://github.com/hackallcode) (бэкенд)

### HowTo Run:
* go get liokor_mail/cmd/main
* go run liokor_mail/cmd/main

### Другие команды:
* Тесты: go test -coverpkg=./... -cover ./... -coverprofile test_cover
* Подробное покрытие (запускать после команды выше): go tool cover -func test_cover
* Покрытие в html: go tool cover -html=test_cover
* Автоформатирование: go fmt liokor_mail/...
* Сборка: go build liokor_mail/cmd/main

### Установка swagger
1. Скачайте последний релиз с https://github.com/swagger-api/swagger-ui/releases
2. Скопируйте все файлы, кроме index.html, из папки dist архива в папку swagger проекта

### Деплой
* [https://mail.liokor.ru](https://mail.liokor.ru)

### Swagger API Docs:
* [https://api.mail.liokor.ru/api/swagger/][https://api.mail.liokor.ru/api/swagger/]

### Репозиторий с фронтендом
* [2021_1_LioKor](https://github.com/frontend-park-mail-ru/2021_1_LioKor)
