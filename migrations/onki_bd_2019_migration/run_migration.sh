#!/bin/bash

echo "Running Onki BD 2019 Migration..."
echo "This will migrate data from MySQL onki_bd_2019 table to PostgreSQL old_registries table"
echo ""

cd "$(dirname "$0")"
go run main.go
