#!/bin/bash

# This is a test script that contains some errors

name = "John Doe"

if [ $name == "John Doe" ]; then
    echo "Hello, $name"
fi

for i in {1..10}; do
    echo "Number $i"
done

ech "This is a test"

while read line; do
    echo $line
done <"nonexistentfile.txt"

function greet {
    echo "Hello, $1"
}

greet "World"
