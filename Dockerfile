FROM ghcr.io/xkcd-2347/guac:v0.1.0-nightly.20230810 as guac

FROM golang:1.21 as build
WORKDIR /src
COPY . .
RUN go build -o /export .

FROM golang:1.21 as exporter
WORKDIR /export
COPY --from=build /export .
COPY --from=guac /opt/guac/guacone .
# TODO: Pass MinIO credentials as Dockerfile args
#COPY credentials /root/.aws/

ENTRYPOINT ["/export/export"]
