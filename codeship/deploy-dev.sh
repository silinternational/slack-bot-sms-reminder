#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# Build binaries
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
$DIR/build.sh

# Export env vars
export SLACK_APP_ID="${SLACK_APP_ID}"
export SLACK_CLIENT_ID="${SLACK_CLIENT_ID}"
export SLACK_CLIENT_SECRET="${SLACK_CLIENT_SECRET}"
export SLACK_SIGNING_SECRET="${SLACK_SIGNING_SECRET}"
export SLACK_VERIFICATION_TOKEN="${SLACK_VERIFICATION_TOKEN}"

serverless deploy -v --stage dev