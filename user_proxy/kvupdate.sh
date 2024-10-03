#!/bin/bash
while read line
do
echo "line: "$line > /app/logs.txt
done