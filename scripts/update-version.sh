#!/bin/bash

# Get the version
tag=`cat pkg/builder/version`

# Get the latest tag from git
latest_tag=`git describe --tags --abbrev=0`

# Check if the version is the same as the latest tag
if [ "$tag" == "$latest_tag" ]; then
    echo "Version in version file is the same as the latest tag"
else
    git tag -a $tag -m "Release $tag" && \
    git push -f --tag && \
    echo "Tag $tag has been pushed to the repository"
fi
