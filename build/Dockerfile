FROM golang:1.17-alpine as build
WORKDIR /opt/build
COPY . .
RUN go build -o app main.go && chmod +x app

FROM alpine:latest as release
COPY --from=build /opt/build/app /finstar
CMD [ "/finstar" ]