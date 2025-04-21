FROM --platform=linux/amd64 ubuntu:latest

RUN apt-get update && \
    apt-get install -y wget tar

RUN wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.23.3.linux-amd64.tar.gz && \
    rm go1.23.3.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin    

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task_planner .

ENV TODO_PORT=3030

EXPOSE 3030

CMD ["/task_planner"]
