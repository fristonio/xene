#!/bin/bash

set -euo pipefail

APIDOCS_DIR=./pkg/apiserver/docs/

TMP_DIR=`mktemp -d`

trap 'rm -rf $TMP_DIR' EXIT INT TERM

swag init --dir pkg/apiserver/ --output $TMP_DIR --generatedTime=false

if ! $(diff ${APIDOCS_DIR}/swagger.json ${TMP_DIR}/swagger.json); then
  echo "Detected a difference in xene API docs"
  echo "diff: `diff ${APIDOCS_DIR}/swagger.json ${TMP_DIR}/swagger.json`"
  echo "Please rerun 'cd pkg/apiserver/ && swag init' and commit your changes"
  exit 1
fi

echo "[*] API documentation is up to date."
