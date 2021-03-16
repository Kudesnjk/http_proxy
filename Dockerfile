FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o web ./app/web_interface/main.go
RUN go build -o proxy ./app/proxy/main.go
RUN chmod +x run_proxy.sh

EXPOSE 8080
EXPOSE 8000

CMD "./run_proxy.sh"