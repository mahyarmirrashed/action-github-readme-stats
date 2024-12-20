#!/bin/bash

# Initialize an empty array to hold the arguments
args=()

# Use a while loop to split INPUT_INCLUDES by spaces and handle each token
while IFS= read -r -d ' ' token; do
  args+=("$token")
done <<<"$INPUT_INCLUDES "

# Pass the arguments array to your main script
exec app "${args[@]}"
