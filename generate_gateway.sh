#!/bin/bash
set -e

# This script generates gRPC-Gateway stubs from Cloud Tasks proto files
# We use standalone=true mode to only generate the gateway HTTP handlers

# Install required protoc plugins if not present
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Remove old generated files
rm -rf gateway/cloudtasksgw

# Create output directory
mkdir -p gateway/cloudtasksgw

# Set up include paths
PROTO_PATH="temp_googleapis"

# Generate only the gRPC-Gateway stubs (standalone mode)
# Use the newer cloud.google.com/go/cloudtasks packages which implement protoreflect.ProtoMessage
protoc \
    -I "${PROTO_PATH}" \
    --grpc-gateway_out=gateway/cloudtasksgw \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=standalone=true \
    --grpc-gateway_opt=Mgoogle/cloud/tasks/v2/cloudtasks.proto=cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb \
    --grpc-gateway_opt=Mgoogle/cloud/tasks/v2/queue.proto=cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb \
    --grpc-gateway_opt=Mgoogle/cloud/tasks/v2/task.proto=cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb \
    --grpc-gateway_opt=Mgoogle/cloud/tasks/v2/target.proto=cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb \
    --grpc-gateway_opt=Mgoogle/iam/v1/iam_policy.proto=cloud.google.com/go/iam/apiv1/iampb \
    --grpc-gateway_opt=Mgoogle/iam/v1/policy.proto=cloud.google.com/go/iam/apiv1/iampb \
    --grpc-gateway_opt=Mgoogle/iam/v1/options.proto=cloud.google.com/go/iam/apiv1/iampb \
    google/cloud/tasks/v2/cloudtasks.proto

echo "Gateway stubs generated successfully in gateway/cloudtasksgw/"
