#!/bin/sh -x


. ./auth.env

packr build
result=$?
echo result = $result

if [ "$result" != "0" ];then
    echo "go build failed"
    exit 1
fi

./server


