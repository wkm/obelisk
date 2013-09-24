#!/bin/sh
# nasty little shell script

# get the release version, trimming whitespace
version=$(cat version.txt | xargs)
echo "attempting to release $version"

# test if the version was released
if git show-ref --verify --quiet refs/tags/$version
then
	echo "ERR tag $version already created"
	exit 1
else
	echo "OK  tag does not exist"
fi

# let's cut this bad boy
if ! git diff --exit-code
then
	echo "ERR there are uncommitted changes"
	exit 1
else
	echo "OK  all changes committed"
fi

# make the tag
if ! git tag $version
then
	echo "ERR could not create tags"
	exit 1
else
	echo "OK  tagged"
fi

# and push all changes
if ! git push --tags
then
	echo "ERR could not push tags"
	exit 1
else
	echo "OK  pushed."
fi