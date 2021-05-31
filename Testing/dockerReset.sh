#! /bin/bash


#Close any active contianers
docker-compose down

# Prune all volumes - Should we keep this?
yes | docker volume prune

# Start up docker contianers
docker-compose up -d

echo "Initializing Database"
sleep 20
echo "Done"
sleep 5