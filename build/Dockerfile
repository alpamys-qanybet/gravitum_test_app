FROM golang:1.23-alpine AS golangbuild

WORKDIR /app

COPY ../ ./

RUN go mod tidy && go mod vendor
RUN go build ./cmd/gravitum-test-app

FROM alpine:3.14
RUN apk add --no-cache bash curl net-tools tzdata 
ENV TZ="Asia/Aqtau"

WORKDIR /app

COPY --from=golangbuild /app/gravitum-test-app ./gravitum-test-app
COPY --from=golangbuild /app/config/config.yml ./config/config.yml

CMD ["./gravitum-test-app"]