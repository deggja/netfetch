#!/bin/bash

# Fetch tags
git fetch --tags
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)

# Initialize variables
major=0
minor=0
patch=0

if [ ! -z "$latest_tag" ]; then
  # Parse the current version
  IFS='.' read -r major minor patch <<< "${latest_tag}"
fi

# Analyze commit messages since the last tag for versioning
for commit in $(git rev-list $latest_tag..HEAD); do
    message=$(git log --format=%B -n 1 $commit)
    
    if [[ $message == *"#major"* ]]; then
      let major+=1
      minor=0
      patch=0
      break
    elif [[ $message == *"#minor"* ]]; then
      let minor+=1
      patch=0
    elif [[ $message == *"#patch"* ]]; then
      let patch+=1
    fi
done

new_tag="${major}.${minor}.${patch}"

# Set output for the next steps using environment file
echo "new_tag=$new_tag" >> $GITHUB_ENV

# Create the new tag
git tag $new_tag
git push --tags