FROM golang:1.19 AS build

WORKDIR /build
COPY . /build

#RUN go mod download
RUN go build -o bin/fn-test main.go

#FROM gcr.io/distroless/static-debian11
FROM gcr.io/distroless/static@sha256:11364b4198394878b7689ad61c5ea2aae2cd2ed127c09fc7b68ca8ed63219030
COPY --from=build /build/bin/fn-test /usr/local/bin/fn-test

ENTRYPOINT ["/usr/local/bin/fn-test"]