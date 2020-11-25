FROM golang:1.15.5 as base
EXPOSE 8080

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 go build -o app main.go

FROM scratch
COPY --from=base /app .

ENTRYPOINT ["./app"]