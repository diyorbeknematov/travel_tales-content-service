FROM golang:1.22.2 AS builder

WORKDIR /travel-content

COPY . .
RUN go mod download

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../travel_app .

FROM alpine:latest

WORKDIR /travel-content 

COPY --from=builder /travel-content/travel_app .
COPY --from=builder /travel-content/logs/app.log ./logs/
COPY --from=builder /travel-content/.env .

EXPOSE 50051

CMD [ "./travel_app" ]