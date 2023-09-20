FROM golang:1.21.1-alpine3.18 as build

RUN apk add --no-cache --update gcc g++

WORKDIR /code
COPY go.mod .
COPY go.sum .

RUN go mod download

ADD . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w"  -o app .

FROM alpine:3.18.3

COPY --from=build /code/app /app/nanourl
COPY --from=build /code/views /app/views
COPY --from=build /code/static /app/static

VOLUME /app/data
ENV NANOURL_PATH=/app/data
ENV NANOURL_REAL_IP_HEADER=CF-Connecting-IP

WORKDIR /app
CMD ["/app/nanourl"]
