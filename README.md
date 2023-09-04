### Compiling
```bash
go build .
```

### Running


### Localstack
This project uses localstack to test locally against SQS and S3.

1. Run docker-compose
```bash
cd docker-compose
docker-compose up
```
2. Upload a document to a bucket. This can be done by either installing the aws CLI in your local machine, or by `docker exec`ing into the running localstack container.
```bash
aws --endpoint-url=http://localhost:4566 s3api put-object --bucket <bucket name> --key <name of file in bucket> --body <file to be uploaded>
```
3. Uploading a document to S3 will trigger a notification

### Podman
#### Run
```
podman run
   -e KAFKA_TOPIC=sbom-stored
   -e KAKFA_HOSTNAME=localhost
   -e MINIO_HOSTNAME=localhost
   -e AWS_ACCESS_KEY_ID=<access key id>
   -e AWS_SECRET_ACCESS_KEY=<secret access key>
   --network host -it localhost/trust-exporter
```

### SSL in minio
1. Generate or copy certs in minio container
    1. Can use the utility from minio
2. Add cert to local trust chain (a bit cumbersome?)
   1. Put public.crt in /etc/pki/ca-trust/source/anchors
   2. Run update-ca-trust
3. Create access key in minio
   1. Put in ~/.aws/credentials

### AWS credentials file
```
[default]
aws_access_key_id = key id
aws_secret_access_key = access key
```

