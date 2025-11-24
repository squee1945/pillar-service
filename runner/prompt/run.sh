#!/bin/bash

set -eo pipefail

TEMP_DIR="${TMPDIR:-/tmp}"
PROMPT_FILE=$(mktemp "${TEMP_DIR}/prompt-XXX")
SETTINGS_FILE="${HOME}/.gemini/settings.json"

gcscp --gcs-path="${PROMPT_GCS_PATH:?}" --local-path="${PROMPT_FILE}"

echo gcscp --gcs-path="${SETTINGS_GCS_PATH:?}" --local-path="${SETTINGS_FILE}"
gcscp --gcs-path="${SETTINGS_GCS_PATH:?}" --local-path="${SETTINGS_FILE}"

echo ""
echo "*** SETTINGS ***********************************************************"
cat  "${SETTINGS_FILE}" | sed -E 's/ghs_[^"]*"/ghs_<redacted>"/g'
echo "************************************************************************"
echo ""

echo ""
echo "*** TOOLS **************************************************************"
gemini -p "list the tools that you have available"
echo "************************************************************************"
echo ""

echo ""
echo "*** PROMPT *************************************************************"
cat  "${PROMPT_FILE}"
echo "************************************************************************"
echo ""

cd ${REPO:?}
cat "${PROMPT_FILE}" | gemini --yolo --debug --model=gemini-2.5-pro 2>&1
