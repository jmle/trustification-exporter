trustification-exporter
--
### Building with Docker
```bash
docker buildx build . -t localhost/trustification-exporter
```

### Running with Docker (example)
```bash
docker run --network host          \
  -e KAFKA_HOSTNAME=localhost      \
  -e KAFKA_PORT=9092               \
  -e QUEUES=sbom-stored,vex-stored \
  -e S3_HOSTNAME=localhost         \
  -e S3_PORT=9000                  \
  -e MESSAGE_PROVIDER=kafka localhost/trustification-exporter
```

### Configuration
Configuration is done through the following environment variables, shown here with sample local values:
```
KAFKA_HOSTNAME=localhost;
KAFKA_PORT=9092;
MESSAGE_PROVIDER=kafka/sqs;
S3_HOSTNAME=localhost;
S3_PORT=9000;
SQS_HOSTNAME=localhost;
SQS_PORT=4566;
QUEUES=sbom-stored,vex-stored
```
The `S3_*` variables are used for both MinIO and S3.

### Running locally
####  Localstack: S3 + SQS
**NOTE**: At the moment, localstack only works with docker.
1. Run docker-compose
```bash
cd docker-compose
docker-compose -f compose-aws.yaml up
```
2. Upload a document to a bucket. This can be done by either installing the aws CLI in your local machine, or by `docker exec`'ing into the running localstack container.
```bash
aws --endpoint-url=http://localhost:4566 s3api put-object --bucket <bucket name> --key <name of file in bucket> --body <file to be uploaded>
```
3. Uploading a document to S3 will trigger a notification onto the SQS queue.

#### MinIO + Kafka
1. Run docker-compose
```bash
cd docker-compose
docker-compose -f compose-local.yaml up
```
2. Log into minio (http://localhost:9001) with credentials `admin/password` and upload a document into a bucket.
3. An event will be triggered and sent to the Kafka topic.


### SSL in MinIO/S3
SSL is deactivated for MinIO, so only HTTP will work. In order to activate SSL:
1. Generate or copy certs in minio container
    1. Can use the utility from minio
2. Add cert to local trust chain (a bit cumbersome?)
   1. Put public.crt in /etc/pki/ca-trust/source/anchors
   2. Run update-ca-trust
3. Create access key in minio
   1. Put in ~/.aws/credentials

### AWS credentials file
Both docker-compose files are set up to use the following credentials:
```
[default]
aws_access_key_id=key_id
aws_secret_access_key=access_key
```
You can set these up in `~/.aws/credentials`.
