#!/bin/bash

# Fetch tags
git fetch --tags

# Get the latest tag
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

# If there are no tags yet, start with 0.0.1
if [ -z "$latest_tag" ]; then
  new_tag="0.0.1"
else
  # Increment the version
  major=$(echo $latest_tag | cut -d. -f1)
  minor=$(echo $latest_tag | cut -d. -f2)
  patch=$(echo $latest_tag | cut -d. -f3)
  if [ $patch -lt 99 ]; then
    let patch+=1
  else
    let minor+=1
    patch=0
  fi
  new_tag="${major}.${minor}.${patch}"
fi

# Set output for the next steps
echo "Setting new tag to $new_tag"
echo "::set-output name=new_tag::$new_tag"

# Create the new tag
git tag $new_tag