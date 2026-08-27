FROM golang:1.21
ENV GOTOOLCHAIN=local
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/server .
EXPOSE 8080
CMD ["/app/bin/server"]
