# HTTP PROXY

Поддерживает только `http` 

Запуск:

* `git clone https://github.com/Kudesnjk/http_proxy.git`
* `cd http_proxy`
* `docker-compose build`
* `docker-compose up`

Использование:

Прокси - `localhost:8080`

Веб-интерфейс - `localhost:8000`

Использование веб-интерфейса:

`/requests` – список запросов

`/requests/id` – вывод 1 запроса

`/repeat/id` – повторная отправка запроса

`/scan/id` – сканирование запроса на наличие XSS
