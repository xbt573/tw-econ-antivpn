FROM golang:alpine as build

WORKDIR /app

COPY . .
RUN go build

FROM alpine:latest

COPY --from=build /app/tw-econ-antivpn /
ENTRYPOINT ["/tw-econ-antivpn"]