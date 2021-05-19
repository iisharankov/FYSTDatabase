#! /bin/bash

# Init Env Variables
. initEnvVars.sh


# Restart Docker
. dockerReset.sh


NUMBEROFFILES=1000


for fileID in $(seq 1 $NUMBEROFFILES); do

# set variables to use during upload
echo $fileID / $NUMBEROFFILES
filename="test"
instrument=1
hash="notARealHash"
theDate=$(date +"%Y-%m-%dT%T.%NZ")
size=1024
url=$(realpath $0)

# Request to add a new file to the database
curl 'localhost:8700/files' -d@- <<___
{
  "name": "$filename",
  "instrument": $instrument,
  "md5sum": "$hash",
  "date_created": "$theDate",
  "size": "$size",
  "url": "$URL"
}
___

# Add row in Record table to say object was uploaded to server (It wasn't, but if it was this would be called)
curl "localhost:8700/logs/$fileID/fyst" -d@- <<___
{  :"NAME"  }
___


# Add a row to the Copies table saying the object was uplaoded to the server (It wasn't, but just to fill it)
curl "localhost:8700/files/$fileID/copies" -d@- <<___
{
  "file_id": $fileID,
  "location_id": 1,
  "url": "$(realpath $0)"
}
___


# Add record to Germany so server doesn't try to request object that does not exist (never uploaded by minio
curl "localhost:8700/logs/$fileID/germany" -d@- <<___
{  :"NAME"  }
___

# Add record to Toronto so server doesn't try to request object that does not exist (never uploaded by minio 
curl "localhost:8700/logs/$fileID/toronto" -d@- <<___
{  :"NAME"  }
___

done
#1 1 date
#1 3 date
#1 5 date
