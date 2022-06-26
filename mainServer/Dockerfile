FROM golang:1.18
WORKDIR /app
COPY . /app
RUN go get -v -d ./...
RUN go build
EXPOSE 8080
CMD [ "go", "run", "main.go", "-dockerconfig"]