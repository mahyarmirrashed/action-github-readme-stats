#!/bin/bash

# Initialize an empty array to hold the arguments
args=()

IFS=' ' read -r -a args <<<"$1"

# Pass the arguments array to your main script
exec app "${args[@]}"
