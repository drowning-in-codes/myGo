# run go application in docker
FROM golang:1.20.13-bookworm
WORKDIR /app
COPY . .
ENV GOPROXY="https://goproxy.cn"
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
