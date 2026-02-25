# Структура проекта Планировщик задач

├── Dockerfile
├── go.mod
├── go.sum
├── main.go                     # Точка входа: настройка env (порт, БД), инициализация БД и api, запуск сервера
├── pkg
│   ├── api
│   │   ├── addtask.go          # POST /api/task - добавление задачи (валидация даты, title, repeat)
│   │   ├── api.go              # Регистрация всех API-хендлеров и маршрутизация в taskHandler
│   │   ├── auth.go             # auth() для проверки JWT-токена (с проверкой хэша пароля)
│   │   ├── done.go             # POST /api/task/done - выполнение задачи (удаление или перенос даты)
│   │   ├── nextdate.go         # NextDate() для вычисления следующей даты по правилам (d/y/w/m)
│   │   ├── singin.go           # POST /api/signin - генерация JWT-токена при входе
│   │   └── tasks.go            # GET /api/tasks - список задач с поддержкой search
│   ├── db
│   │   ├── db.go               # Инициализация БД (Init()), схема таблицы scheduler
│   │   └── task.go             # Операции AddTask, GetTask, UpdateTask, DeleteTask, UpdateDate, Tasks()
│   └── server
│       └── server.go           # Запуск HTTP-сервера на указанном порту с файл-сервером web/
├── README.md
├── scheduler.db
├── tests...                    # файлы тестов из задания
└── web...                      # файлы web из задания

# Переменные окружения

TODO_PORT                       # порт сервера (по умолчанию: 7540)
TODO_DBFILE                     # файл БД (по умолчанию: scheduler.db)
TODO_PASSWORD                   # пароль для входа (если задан, то браузер запросит его)

# Запуск в Docker

docker run -d -p 7540:7540 \
  -e TODO_PASSWORD=ВАШПАРОЛЬ \  # если нужен пароль
  todo-scheduler
  
# Запуск локально

go run main.go                  # из папки проекта

# Запуск всех тестов

go test -v ./tests              # из папки проекта
