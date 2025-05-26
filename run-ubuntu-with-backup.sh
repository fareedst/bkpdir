#!/usr/bin/env bash
##
# run-ubuntu-with-backup.sh
#
# Starts an Ubuntu 20.04 Docker container named "my-ubuntu-container"
# and mounts the host's `bin/bkpdir-ubuntu20.04` into the container's
# `/usr/local/bin` (read-only), making it available on the PATH.
#
# Set RUWB_KEEP_CONTAINER=1 to prevent container removal on exit
#
# Compatible with both macOS and Ubuntu
##

# Exit on error, unset variables, and pipe failure
set -euo pipefail

# bkpdir version
: "${RUWB_BKPDIR_VERSION:=20.04}"

# Ubuntu version
: "${RUWB_UBUNTU_VERSION:=20.04}"

# Local path to binary
: "${RUWB_HOST_BINARY_NAME:=bin/bkpdir-ubuntu${RUWB_BKPDIR_VERSION}}"

# Path to binary in container, on the PATH
: "${RUWB_CONTAINER_BINARY_NAME:=/usr/local/bin/bkpdir}"

# Control container removal
: "${RUWB_KEEP_CONTAINER:=0}"

# Container name
: "${RUWB_CONTAINER_NAME:=my-ubuntu-container}"

# Resolve the absolute path to the host binary
# Use readlink -f on Linux, fallback to simpler method on macOS
if command -v readlink >/dev/null 2>&1 && readlink -f / >/dev/null 2>&1; then
    # Linux
    HOST_BINARY="$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")/${RUWB_HOST_BINARY_NAME}"
else
    # macOS
    HOST_BINARY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/${RUWB_HOST_BINARY_NAME}"
fi

# Fix volume mount point
VOLUME_MOUNT_POINT="${HOST_BINARY}:${RUWB_CONTAINER_BINARY_NAME}:ro"

# Build docker run options
DOCKER_OPTS=(-it --name "${RUWB_CONTAINER_NAME}" -v "${VOLUME_MOUNT_POINT}")
if [ "${RUWB_KEEP_CONTAINER}" = "0" ]; then
    DOCKER_OPTS+=(--rm)
fi

# Print the command to run the container
# Use tput for colors if available, fallback to plain text
if command -v tput >/dev/null 2>&1; then
    PREFIX="$(tput setaf 14)$(tput setab 5)++$(tput sgr0) "
else
    PREFIX="++ "
fi
echo -e >&2 "${PREFIX}docker run ${DOCKER_OPTS[*]} ubuntu:${RUWB_UBUNTU_VERSION}"

docker run "${DOCKER_OPTS[@]}" ubuntu:${RUWB_UBUNTU_VERSION}
