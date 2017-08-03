# go-spotippos

## Run now with dply.co
### Click here (Run free for two hours):
[![Dply](https://dply.co/b.svg)](https://dply.co/b/BwXv6nmR) 

## How to test
### Expose your configs
```
export ADDR='104.131.130.30'
export PORT='80'

```
### Get properties in rectangle area
```
curl --request GET \
     --url "http://${ADDR}:${PORT}/properties?ax=0&ay=1000&bx=600&by=500"
```
### Get one propertie by ID
```
curl --request GET \
     --url "http://${ADDR}:${PORT}/properties/34"
```
### Register a new propertie
```
curl --request POST \
    --url "http://${ADDR}:${PORT}/properties" \
    --header 'content-type: application/json' \
    --data '{"title": "Imóvel código 1, com 3 quartos e 2 banheiros.","price": 643000,"description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.","x": 1257,"y": 929,"beds": 3,"baths": 2,"squareMeters": 61}'
```

## Build
```
git clone git@github.com:lichti/go-spotippos.git
cd go-spotippos
./utils/build.sh
```

## Run
```
./utils/run.sh <port>
docker ps | grep 'lichti/go-spotippos'
```

## Push
```
./utils/push.sh <DockerHubUser>
```

## dply.co script
```
cat ./utils/dplyco.sh
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
```