#!/bin/bash

curl \
  -i \
  -X 'POST' \
  -H "Content-Type: application/json"  \
  --data @./docs/provinces.json  \
  http://127.0.0.1:8000/province/populate

curl \
  -i \
  -X 'POST' \
  -H "Content-Type: application/json"  \
  --data @./docs/properties.json  \
  http://127.0.0.1:8000/properties/populate


