#!/bin/bash 

echo "Docker pull"
docker pull lichti/go-spotippos
echo "Docker run"
docker run -d -p 80:8000 lichti/go-spotippos
echo "Waiting"
sleep 5s
echo "Getting properties"
curl -Ls https://github.com/lichti/go-spotippos/raw/master/docs/properties.json -o /tmp/properties.json
echo "Getting provinces"
curl -Ls https://raw.githubusercontent.com/lichti/go-spotippos/master/docs/provinces.json -o /tmp/provinces.json
echo "Populating provinces"
curl -i -X 'POST' -H "Content-Type: application/json" --data @/tmp/provinces.json http://127.0.0.1/provinces/populate
echo "Populating properties"
curl -i -X 'POST' -H "Content-Type: application/json"  --data @/tmp/properties.json http://127.0.0.1/properties/populate