#!/bin/bash
cd /home/aro/Desktop/library
go build -o /tmp/go-micro-gen .
mkdir -p /tmp/tester
cd /tmp/tester

for db in postgres mongo none; do
  for broker in kafka rabbitmq nats none; do
    for transport in http grpc both; do
      echo "Testing db=$db broker=$broker transport=$transport"
      rm -rf testsvc
      /tmp/go-micro-gen generate --name testsvc --module x/y --db $db --broker $broker --transport $transport --arch clean --docker=false --cloud none --ci none --redis=false --k8s=false --helm=false --output testsvc --yes > /dev/null
      cd testsvc
      go mod tidy > /dev/null 2>&1
      go build ./... > /dev/null
      if [ $? -ne 0 ]; then
        echo "FAIL: db=$db broker=$broker transport=$transport"
        go build ./...
        exit 1
      fi
      cd ..
    done
  done
done
echo "ALL PASSED!"
