FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/demo

FROM alpine:latest
RUN mkdir /vol
COPY --from=build /bin/demo /bin/demo
ENTRYPOINT ["/bin/demo","--path =/vol","--url=http://filestats:1999/files","--help"]