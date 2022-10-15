FROM  golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -v .
EXPOSE 5000
