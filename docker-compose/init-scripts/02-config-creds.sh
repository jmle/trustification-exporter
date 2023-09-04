#!/bin/bash

cat <<EOF >> ~/.aws/credentials
[default]
aws_access_key_id=access_key
aws_secret_access_key=secret_key
EOF