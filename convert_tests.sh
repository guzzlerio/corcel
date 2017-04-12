#! /bin/bash
for i in `find ./ -type f -name "*_test.go"`
do
    ./convert_test.sh $i  
done
