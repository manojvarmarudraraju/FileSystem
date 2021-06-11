FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/demo

FROM scratch
COPY --from=build /bin/demo /bin/demo
ENTRYPOINT ["/bin/demo","--port=:1999"]