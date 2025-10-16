#!/bin/bash

set -eo pipefail

echo "Cloning ${OWNER:?}/${REPO:?}"
git clone \
  --depth 1 \
  --branch "${DEFAULT_BRANCH:?}" \
  --single-branch \
  "https://x-access-token:${GITHUB_TOKEN:?}@github.com/${OWNER:?}/${REPO:?}.git"

cd ${REPO:?}

echo "Creating branch ${DEV_BRANCH:?}"
git switch -c "${DEV_BRANCH:?}"
