#!/bin/bash
## This script is used to update the version of the project
## - If you provide a symantec versioning number as an argument (I.e v1.2.3)
##   it will update the version file with the provided version
## - If you don't provide a version, it will check if the version file has a version
##   (without updating the version file) and tag the repository with the version in the file

# Check if the version is set
versionToSet=$1
if [ "$versionToSet" != "" ]; then
    # Ensure the version follows the semantic versioning
    if [[ ! $versionToSet =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then

        # It's acceptable to have the version without the v prefix
        # But it will be appended to follow the same pattern
        if  [[ $versionToSet =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        versionToSet="v$versionToSet"
        else
            echo "Version should follow the semantic versioning pattern (I.e v1.2.3)"
            exit 1
        fi

    fi

    echo "Setting version to $versionToSet"
    echo $versionToSet > pkg/builder/version
fi

# Get the version
tag=`cat pkg/builder/version`

# Ensure the version file has a version and that the version is valid
if [ "$tag" == "" ]; then
    echo "Version file does not have a version"
    exit 1
fi
if [[ ! $tag =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Version in version file is not valid"
    exit 1
fi

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
