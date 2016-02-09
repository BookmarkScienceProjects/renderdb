#!/usr/bin/env bash
BASE64FILE=$(mktemp)
cat $1 | jq --compact-output -r '.[].geometryData' > $BASE64FILE
echo $BASE64FILE

i=0
while read line;
do
	i=$((i+1))
	echo $line | base64 --decode > $2/$i.obj
done < $BASE64FILE
rm -f $BASE64FILE

