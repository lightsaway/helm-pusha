
## Building

Below command will produce archive with binaries both for osx and linux

````
  make all
````

## Installing

Just drop binary and plugin.yaml to $(HELM_HOME)/plugins/push


## Using
First add repo where you want to push/pull

````
helm repo add release http://my.release.com
````

Then use it :

````
helm push release my-awesome-chart-1.0.0.tgz
````
By default plugin assumes that there is /upload endpoint at the repo that accepts PUT requests

You can specify different endpoint by changing HELM_CHARTS_UPLOAD_ENDPOINT_PATH environment var

````
HELM_CHARTS_UPLOAD_ENDPOINT_PATH="/charts/upload/endpoint" helm push release my-awesome-chart-1.0.0.tgz 
````