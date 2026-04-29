# build
FROM golang:1.26.1-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o hireradar ./cmd/bot

# final
FROM alpine:latest

RUN addgroup -g 1000 -S appuser && adduser -u 1000 -S appuser -G appuser
COPY --from=build /app/hireradar /hireradar
EXPOSE 8080
USER appuser
CMD ["/hireradar"]



