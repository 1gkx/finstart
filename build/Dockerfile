FROM golang:1.17-alpine as build
WORKDIR /opt/build
COPY . .
RUN go mod download
RUN go build -o app main.go 

FROM scratch as release
COPY --from=build /opt/build/app /finstar
ENTRYPOINT ["/finstar"]