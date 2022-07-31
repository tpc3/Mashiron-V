FROM golang:alpine AS build
ADD . /go/src/Mashiron-V/
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 0
WORKDIR /go/src/Mashiron-V
RUN go build .

FROM alpine
COPY --from=build /go/src/Mashiron-V/Mashiron-V /go/src/Mashiron-V/config.yaml /Mashiron-V/
RUN apk add --no-cache ca-certificates
WORKDIR /Mashiron-V
ENTRYPOINT [ "/Mashiron-V/Mashiron-V" ]