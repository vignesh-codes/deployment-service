FROM golang:1.22-alpine3.19

RUN apk update
RUN apk add curl

COPY . .

ENV GOPATH=${pwd}
# RUN go mod tidy
RUN cd ./ && unset GOPATH && go get

CMD go run *.go
