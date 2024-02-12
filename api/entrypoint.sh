#!/bin/bash

# Decode the base64 C++ source code and write it to a file
echo "$CPP_SOURCE_BASE64" | base64 -d > source.cpp

# Compile the C++ source code
if g++ source.cpp -o output; then
    # Run the compiled program and redirect stdout and stderr to different files
    (/usr/bin/time -f "\nElapsed Time: %e\nMax Memory: %M" ./output) > program_output.txt 2> time_output.txt

    # Extract the maximum memory usage and time taken
    TIME_OUTPUT=$(cat time_output.txt)
    MAX_MEMORY=$(echo "$TIME_OUTPUT" | grep "Max Memory" | cut -d: -f2 | tr -d ' ')
    TIME_TAKEN_SECONDS=$(echo "$TIME_OUTPUT" | grep "Elapsed Time" | cut -d: -f2 | tr -d ' ')

    # Read the output of the program from the file
    PROGRAM_OUTPUT=$(cat program_output.txt)

    # Echo the maximum memory usage and time taken
    echo "Maximum memory usage: $MAX_MEMORY KB"
    echo "Time taken: $TIME_TAKEN_SECONDS seconds"

    # Echo the output of the program
    echo "Program output: $PROGRAM_OUTPUT"

    # Return stdout
    echo "$PROGRAM_OUTPUT"
else
    # Return stderr
    error=true
    echo "Compilation failed." >&2
    echo "$error"
fi

sleep 25
