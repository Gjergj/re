FROM golang:1.21.4 as builder

WORKDIR /opt/packagingcalculator

COPY ./ ./
COPY ./go.* .
COPY ./cmd/packagingcalculator/assets ./assets

RUN go build -o bin ./cmd/packagingcalculator/...

RUN cp -r /opt/packagingcalculator/cmd/packagingcalculator/migrations migrations

ENTRYPOINT [ "/opt/packagingcalculator/bin/packagingcalculator" ]