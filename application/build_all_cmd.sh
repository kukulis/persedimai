#!/bin/bash

go build -o bin/webapp ./cmd/webapp
go build -o bin/seeder ./cmd/seeder
go build -o bin/createclusters ./cmd/createclusters
go build -o bin/createdb ./cmd/createdb

