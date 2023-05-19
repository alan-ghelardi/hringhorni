#!/usr/bin/env bash
set -euo pipefail



while true; do
    curl http://localhost:8080/
    echo "--------------------------------------------------------------------------------"
    sleep 0.1
done
