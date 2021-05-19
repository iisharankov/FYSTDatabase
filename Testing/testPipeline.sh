#! /bin/bash
# Script to automatically create huge sums of files and pass them throught he FYST Data Migration Pipeline

# Init Env Variables
. initEnvVars.sh


# Restart Docker
#. dockerReset.sh

#Constants
NUMBEROFFILES=50
SIZE=1 # In MB

echo "Creating $NUMBEROFFILES files of size $SIZE MB and sending to FYST database via client."

# Remove any old files created before
#rm file*.txt

fallocate -l ${SIZE}M file0.txt

for element in $(seq 1 $NUMBEROFFILES)
do
echo "--"

# Create a random file
filename=file${element}
#fallocate -l ${SIZE}M $filename.txt

echo "changing file$((element-1)).txt to $filename.txt"
mv file$((element-1)).txt $filename.txt
full_path=$(readlink -f $filename.txt) #  Gives path to file given
#full_path=$(realpath $0) # Gives bash file!


# Set up flags for upload command
hash=`/bin/echo $filename | /usr/bin/md5sum | /bin/cut -f1 -d" "`
cur_time=$(date +"%Y-%m-%dT%T.%NZ")
filesize=$(stat --printf="%s" "$filename.txt")

# upload to FNYST database
client upload -name=$filename -instrument=1 -md5=$hash -date=$cur_time -size=$filesize -url=$full_path

done


