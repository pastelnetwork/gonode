#!/bin/bash
# if we are testing a PR, merge it with the latest master branch before testing
# this ensures that all tests pass with the latest changes in master.

set -o nounset
set -o errexit
set -eu -o pipefail

err=0

(set -x && git pull --ff-only origin "refs/heads/master") || err=$?

if [ "$err" -ne "0" ]; then
    echo
    echo -e "\033[0;31mERROR: Failed to merge your branch with the latest master."
    echo -e "Please manually merge master into your branch, and push the changes to GitHub.\033[0m"
    exit $err
fi
