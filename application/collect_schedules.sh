#!/bin/bash
set -e

for j in {1..10}
do

for i in {1..10}
do
  echo "Iteration $j $i of [10, 10]"
  go run cmd/collectschedules/main.go -env prod -airport '*' -start 2026-04-01 -end 2026-05-01
done

done