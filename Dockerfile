FROM --platform=linux/amd64 ubuntu:latest

ENV TODO_PORT=8080

RUN apt-get update && \
    apt-get install -y golang-go && \
    apt-get clean

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o TaskPlanner .

EXPOSE ${TODO_PORT}

CMD ["./TaskPlanner"]
