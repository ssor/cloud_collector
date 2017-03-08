#!/bin/bash

result=""

# for host in "172.16.1.11" "172.16.1.41" "172.16.1.42"
for host in "127.0.0.1"
do 
    # echo $host
    conn=`mongostat --host=${host} --rowcount=1 --noheaders  | awk '{print $16}'`
    # echo "conn -> "${conn} 
    result=${result}${host}"->"${conn}"|"
done 
echo "mongostat:::"${result}
