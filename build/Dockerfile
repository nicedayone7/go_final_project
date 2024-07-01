FROM golang:1.22

WORKDIR /app

ENV TODO_PORT="7540" TODO_DBDILE="scheduler.db" TODO_PASSWORD="password" TODO_SECRET="practicum"

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app_task

EXPOSE 7540

CMD ["/my_app_task"]