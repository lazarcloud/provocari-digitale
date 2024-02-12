#!/bin/bash

# Decode the base64 encoded C++ source file
echo "$CPP_SOURCE_BASE64" | base64 -d > source.cpp

# Compile the C++ source file
g++ -o output source.cpp

# Execute the compiled binary and capture the output
output=$(./output)

# Print the output
echo "Output:"
echo "$output"

# Sleep for 10 seconds before deleting the container
sleep 10
