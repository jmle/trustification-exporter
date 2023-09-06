FROM ghcr.io/xkcd-2347/guac:v0.1.0-nightly.20230810 as guac

FROM golang:1.21 as build
WORKDIR /src
COPY . .
RUN go build -o /export .

FROM golang:1.21 as exporter
WORKDIR /export
COPY --from=build /export .
COPY --from=guac /opt/guac/guacone .

RUN mkdir ~/.aws && \
    cat <<EOF > ~/.aws/credentials
    [default]
    aws_access_key_id=access_key
    aws_secret_access_key=secret_key
EOF

ENTRYPOINT ["/export/export"]
