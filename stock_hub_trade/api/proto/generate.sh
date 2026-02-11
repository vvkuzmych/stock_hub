#!/bin/bash

# Generate Go code from proto files
# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔧 Generating Go code from proto files...${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Install it:"
    echo "   brew install protobuf"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "❌ protoc-gen-go not found. Install it:"
    echo "   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ protoc-gen-go-grpc not found. Install it:"
    echo "   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Generate Go code
protoc --go_out=../.. --go_opt=paths=source_relative \
       --go-grpc_out=../.. --go-grpc_opt=paths=source_relative \
       api/proto/*.proto

echo -e "${GREEN}✅ Proto files generated successfully!${NC}"
echo ""
echo "Generated files:"
echo "  - api/proto/model.pb.go"
echo "  - api/proto/stock_hub.pb.go"
echo "  - api/proto/stock_hub_grpc.pb.go"
