#!/bin/bash

set -eo pipefail

TEMP_DIR="${TMPDIR:-/tmp}"
PROMPT_FILE=$(mktemp "${TEMP_DIR}/prompt-XXX")
SETTINGS_FILE="${HOME}/.gemini/settings.json"

gcscp --gcs-path="${PROMPT_GCS_PATH:?}" --local-path="${PROMPT_FILE}"

echo gcscp --gcs-path="${SETTINGS_GCS_PATH:?}" --local-path="${SETTINGS_FILE}"
gcscp --gcs-path="${SETTINGS_GCS_PATH:?}" --local-path="${SETTINGS_FILE}"

echo ""
echo "*** PROMPT *************************************************************"
cat  "${PROMPT_FILE}"
echo "************************************************************************"
echo ""

cd ${REPO:?}
cat "${PROMPT_FILE}" | gemini --yolo --debug 2>&1
