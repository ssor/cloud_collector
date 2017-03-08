#!/bin/bash

result=`sudo netstat -apn | grep ESTABLISHED | grep 27017`
echo "netstat:::"${result}
