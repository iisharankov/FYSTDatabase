#!/bin/bash

echo "Logging into database"
password=iisharankov

/usr/bin/expect<<EOF
set timeout 30
spawn /bin/bash
send "mysql -h 127.0.0.1 -P 3306 -u iisharankov -p  < ../Database/profiling.sql \r"
expect "password:"
send "$password\r"
expect "> "

# Read DB
#send "use mydb;\r"
#expect "$ "
#send "select * from Rules;\r"
#expect "$ "
send "quit \r"
EOF