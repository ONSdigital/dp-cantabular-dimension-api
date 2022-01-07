#!/bin/bash -eux

pushd dp-cantabular-dimension-api
  make build
  cp build/dp-cantabular-dimension-api Dockerfile.concourse ../build
popd
