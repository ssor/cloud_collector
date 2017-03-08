#!/bin/bash

result=`sudo lsof -nP -iTCP  | grep mongo `
echo "netstat:::"${result} 
