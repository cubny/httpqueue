#!/bin/bash

set -e
echo "Load testing the API server using -c 100 -n 10000"
ab -c 100 -n 10000 -p createitems.json -T application/json http://localhost:8080/timers
