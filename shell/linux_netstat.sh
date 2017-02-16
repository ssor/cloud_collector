#!/bin/bash

sudo netstat -apn | grep ESTABLISHED | grep 27017
