#!/bin/bash

# Fetch tags
git fetch --tags
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)

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

# Set output for the next steps using environment file
echo "new_tag=$new_tag" >> $GITHUB_ENV

# Create the new tag
git tag $new_tag
git push --tags