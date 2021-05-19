#! /bin/bash

echo "Setting Env vars"
#Set up env variables
export MINIO_ENDPOINT=0.0.0.0:9001
export MINIO_ACCESS_ID=iisharankov
export MINIO_SECRET_KEY=iisharankov
export MINIO_SSL=false

export S3_ENDPOINT=0.0.0.0:9002
export S3_ACCESS_ID=iisharankov
export S3_SECRET_KEY=iisharankov
export S3_SSL=false

export DATABASE_NAME=mydb
export MYSQL_IP=0.0.0.0
export MYSQL_USER=iisharankov
export MYSQL_PASSWORD=iisharankov
