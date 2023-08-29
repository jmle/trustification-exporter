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

