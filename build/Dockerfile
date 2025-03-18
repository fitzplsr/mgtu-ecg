FROM golang:1.24.0-alpine as builder

COPY . /github.com/fitzplsr/mgtu-ecg/
WORKDIR /github.com/fitzplsr/mgtu-ecg/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go

FROM scratch AS runner

WORKDIR /docker-main/

COPY --from=builder /github.com/fitzplsr/mgtu-ecg/.bin .
COPY --from=builder /github.com/fitzplsr/mgtu-ecg/config config/

EXPOSE 4000
ENTRYPOINT ["./.bin"]