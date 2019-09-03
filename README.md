# Helm-client

Helm-client provides RESTful API that enables manage Helm repositories and deploy charts.

Click [here](https://github.com/helm/helm) to learn more about Helm.

## Prerequisites
- [Go 1.11+ installation or later](https://github.com/golang/go)
- [Beego version 1.10.0 or later](https://github.com/astaxie/beego)
- [Helm and Tiller installed](https://github.com/helm/helm/blob/master/docs/quickstart.md)

## Getting Started

#### Start Helm Tiller Server
Detail see [official docs](https://helm.sh/docs/helm/#helm)

#### Satrt Helm-client Api server in Local
```
bee run -downdoc=true -gendoc=true
```
if you don't want to use beego, you can 
```
go run main.go
``` 

#### [optional] Start Helm-client Api server with Docker 
```
docker build -t zhihui/helm-client:latest .

docker run --rm -it \
  -p 8080:8080 \
  --name helm-client \
  -v $(pwd)/.helm:/helm \
  -e TILLER_HOST=${tillerHost} \
  -e HELM_HOME=/helm \
  -e SKIP_REFRESH=true \
  zhihui/helm-client:latest
```
#### Api Command Examples

List all releases
```
curl http://localhost:8080/v1/releases

Results:
{
  "Next": "",
  "Releases": [
    {
      "Name": "wobbling-prawn",
      "Revision": 1,
      "Updated": "Mon Apr 22 16:15:12 2019",
      "Status": "DEPLOYED",
      "Chart": "fabric-ca-0.1.0",
      "AppVersion": "1.0",
      "Namespace": "default"
    }
  ]
}
```

Install Helm Charts 
```
curl -s -X POST http://localhost:8080/v1/releases \
    -H "content-type: application/json" \
    -d '{
        "chartName": "stable/mysql",
        "namespace": "helm-client"
    }'
    
Results:
{
  "chartIcon": "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
  "chartName": "mysql",
  "chartVersion": "0.15.0",
  "name": "nuanced-goat",
  "namespace": "helm-client",
  "status": "DEPLOYED",
  "updated": "Mon Apr 22 16:53:45 2019"
}
```

## Api Docs
[http://127.0.0.1:8080/swagger/](http://localhost:8080/swagger/)
