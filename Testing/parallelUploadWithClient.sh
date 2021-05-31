#! /bin/bash
# Script to automatically create huge sums of files and pass them throught he FYST Data Migration Pipeline

. initEnvVars.sh  # Init Env Variables
#. dockerReset.sh  # Restart Docker


#Constants
NUMBEROFFILES="${NUMBEROFFILES:-5}"  # Default 10
PARALLELLOOPS="${PARALLELLOOPS:-10}"  # Default 10
SIZE="${FILESIZE:-1}" # In MB

echo "Starting $PARALLELLOOPS loops which will each create $NUMBEROFFILES files and send them to the FYST database via client"


# Filename must be unique so generate random string to append to filename
rand_seed=$[($RANDOM%1000)+1] 

uploadFileWithClient(){
    # random sleep (0-1s) to randomize upload stream timings
    sleep .$[( $RANDOM % 10) + 1 ]s
    
    # Create a random file
    filename=${rand_seed}_${1}_file$2.txt  
    fallocate -l ${SIZE}M $filename
    head -c 1M </dev/urandom >$filename

    full_path=$(readlink -f $filename) #  Gives path to file given
    #full_path=$(realpath $0) # Gives bash file!


    # Set up flags for upload command
    # hash=`/bin/echo $filename | /usr/bin/md5sum | /bin/cut -f1 -d" "`  # Hash of string
    hash=`md5sum ${full_path} | awk '{ print $1 }'`  # Hash of file
    cur_time=$(date +"%Y-%m-%dT%T.%NZ")
    filesize=$(stat --printf="%s" "$filename")

    # upload to FNYST database
    client upload -name=$filename -instrument=1 -md5=$hash -date=$cur_time -size=$filesize -url=$full_path

    rm $filename  # delete file
}

loopUpload(){
    for element in $(seq 1 $NUMBEROFFILES); do
        uploadFileWithClient $1 $element
    done
}

for i in $(seq 1 $PARALLELLOOPS); do
    loopUpload "$i" &
done
wait
echo "Done"



