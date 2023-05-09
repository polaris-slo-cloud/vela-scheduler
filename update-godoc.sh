#!/bin/bash
# set -x
# set -o errexit

# The path, where the static HTML files with the godoc should be stored.
DEST_PATH="./gh-pages/godoc"

if [ ! -d "$DEST_PATH" ]; then
    echo "The destination directory '$DEST_PATH' does not exist"
    exit 1
fi

cd go
DEST_PATH="../$DEST_PATH"

if [ "$(git status --porcelain)" != "" ]; then
    echo "This script must be run on a clean working tree, because it will temporarily modify it."
    exit 1
fi

echo "Rewriting go.work and creating symlinks to temporarily move all modules to the 'src' subfolder"

# Rewrite go.work and temporarily moving all modules to the src subfolder
sed -e "s/\.\//\.\/src\//" ./go.work > ./go.work

mkdir -p src
mv ./cluster-agent ./src/cluster-agent
mv ./context-awareness ./src/context-awareness
mv ./framework ./src/framework
mv ./k8s-connector ./src/k8s-connector
mv ./scheduler ./src/scheduler

# Run godoc in the background
godoc -index -goroot=. -http=:8080 &
godocPid=$!
sleep 10s

# Delete existing files and create temp directory for output
rm -rf $DEST_PATH/*
mkdir "$DEST_PATH/tmp"

# Download all godoc html pages into static files (based on https://github.com/golang/go/issues/2381#issuecomment-66059484)
wget -r -np -N -E -p -k -e robots=off --directory-prefix="$DEST_PATH/tmp" http://localhost:8080/pkg/
# -r  : download recursive
# -np : don't ascend to the parent directory
# -N  : don't retrieve files unless newer than local
# -E  : add extension .html to html files (if they don't have)
# -p  : download all necessary files for each page (css, js, images)
# -k  : convert links to relative
# -e robots=off: ignore robots.txt file, which would otherwise needed to be modified
# --directory-prefix: destination directory for the downloaded files

# Move the files from the temp directory to the DEST_PATH
mv $DEST_PATH/tmp/localhost:8080/* $DEST_PATH
rm -rf "$DEST_PATH/tmp"

kill -SIGTERM $godocPid
echo "Resetting repository to latest commit"
git reset --hard
rm -rf ./src
