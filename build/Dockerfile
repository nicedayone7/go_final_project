FROM golang:1.22.4

WORKDIR /app

ENV TODO_PORT="7540" TODO_DBFILE="scheduler.db" TODO_PASSWORD="password" TODO_SECRET="practicum" CGO_ENABLED="0" GOOS="linux" GOARCH="amd64" 

COPY . .

RUN go mod download

RUN go build -o /my_app_task

CMD ["/my_app_task"]