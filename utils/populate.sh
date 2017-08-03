#!/bin/bash

curl -Ls https://github.com/lichti/go-spotippos/raw/master/docs/properties.json -o /tmp/properties.json

curl -Ls https://raw.githubusercontent.com/lichti/go-spotippos/master/docs/provinces.json -o /tmp/provinces.json

curl -i -X 'POST' -H "Content-Type: application/json" --data @/tmp//provinces.json http://127.0.0.1:${1:-8000}/province/populate

curl -i -X 'POST' -H "Content-Type: application/json"  --data @/tmp//properties.json http://127.0.0.1:${1:-8000}/properties/populate


