FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV PORT=":2565"
ENV DATABASE_URL="postgres://zviytute:cuD_PKAR8GQBTGKKLf2A1pswdtPdi1hA@tiny.db.elephantsql.com/zviytute?sslmode=disable"
ENV AUTH_TOKEN="November 10, 2009"

RUN go test --tags=unit -v ./...
RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-app /app/go-app

CMD ["/app/go-app"]
