FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o main cmd/go_final_project/main.go

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=secure_password

CMD ["./main"]