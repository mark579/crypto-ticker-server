FROM golang:1.17-alpine as BUILD
WORKDIR /opt/crypto-ticker-server
COPY . .
RUN apk add git 
RUN go get -d -v ./...
RUN go build -o crypto-ticker-server

FROM alpine:3.13 as FINAL
COPY --from=BUILD /opt/crypto-ticker-server /bin/
EXPOSE 8080
CMD ["crypto-ticker-server"]