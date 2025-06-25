FROM golang:1.24-alpine
ENV TZ=Asia/Bangkok
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .env
RUN go build -o main ./cmd
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
