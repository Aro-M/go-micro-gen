#!/bin/bash
set -e

echo "=== 🚀 Testing go-micro-gen INIT Mode Injection ==="

# Build the latest binary
go build -o go-micro-gen .
CLI_PATH=$(realpath ./go-micro-gen)

echo ""
echo "--- 1. Testing 'vertical' architecture + mongo + rabbitmq ---"
rm -rf .test-init-vertical
mkdir .test-init-vertical
cd .test-init-vertical
go mod init github.com/testorg/my-existing-app
$CLI_PATH init --arch vertical --db mongo --broker rabbitmq --transport both --redis=false --docker=false --k8s=false --helm=false --cloud=none --ci=none --serverless=false --seeding=false --graphql=false --jwt=false --yes
go mod tidy
go build ./cmd/main.go
echo "✅ Vertical architecture injected successfully!"
cd ..

echo ""
echo "--- 2. Testing 'clean' architecture + postgres + kafka + graphql ---"
rm -rf .test-init-clean
mkdir .test-init-clean
cd .test-init-clean
go mod init github.com/testorg/clean-app
$CLI_PATH init --arch clean --db postgres --broker kafka --graphql=true --seeding=true --redis=false --docker=false --k8s=false --helm=false --cloud=none --ci=none --serverless=false --jwt=false --transport=http --yes
go mod tidy
# Need gofakeit for the seeder
go mod download github.com/brianvoe/gofakeit/v6
go build ./cmd/main.go
go build ./cmd/seed/main.go
echo "✅ Clean architecture + GraphQL + Seeding injected successfully!"
cd ..

echo ""
echo "--- 3. Testing 'ddd' architecture + redis + AWS serverless ---"
rm -rf .test-init-ddd
mkdir .test-init-ddd
cd .test-init-ddd
go mod init github.com/testorg/ddd-app
$CLI_PATH init --arch ddd --db none --broker none --redis=true --cloud aws --serverless=true --docker=false --k8s=false --helm=false --ci=none --seeding=false --graphql=false --jwt=false --transport=http --yes
go mod tidy
go build ./cmd/lambda/main.go
echo "✅ DDD AWS Serverless architecture injected successfully!"
cd ..

echo ""
echo "🎉 All Init Mode tests passed successfully without corrupting existing modules!"
