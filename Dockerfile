FROM  golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy && go build -v .
CMD ["./zettel-bot"]

