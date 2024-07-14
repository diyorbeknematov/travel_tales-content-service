#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/generated

for x in $(find ${CURRENT_DIR}/travel_tales-grpc-proto/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/travel_tales-grpc-proto/ -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done