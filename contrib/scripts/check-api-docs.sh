#!/bin/bash

set -euo pipefail

APIDOCS_DIR=./pkg/apiserver/docs/

TMP_DIR=`mktemp -d`

trap 'rm -rf $TMP_DIR' EXIT INT TERM

swag init --dir pkg/apiserver/ -o $TMP_DIR --generatedTime=false

if ! $(diff -r ${APIDOCS_DIR} ${TMP_DIR}); then
  echo "Detected a difference in xene API docs"
  echo "diff: `diff -r ${APIDOCS_DIR} ${TMP_DIR}`"
  echo "Please rerun 'cd pkg/apiserver/ && swag init' and commit your changes"
  exit 1
fi

echo "[*] API documentation is up to date."
