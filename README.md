# Opap

Made to generate krakend configuration file with swagger.json.
> This project does not give you a complete krakend configuration.
> It generates only the endpoint scopes in the swagger-api that you need to define.


### Usage

* download repostory
* copy `example.config.json` to `config.json`
* update configuration file
* `make run` and create krakend configurations files in `config.json` base path

### Configuration File
**Base Scope**
```sh
base_dir  ->  The path where the krakend files were created
file_prefix  ->  prefix of generated files
prefix_separator  ->  separator of created files
```
**Service Scope**

```sh
name  ->  service name
swagger_url  ->  swagger json url
internal_address  ->  backend host
output_encoding  ->  endpoint output encoding
backend_encoding  ->  backend encoding
excluded_endpoints -> endpoints to be excluded from swagger
```
**Cross Mapping**

A list of routing maps with key value mapping.

Generated separately for each service, use <service_name> for your service name

```sh
from  ->  krakend input url
to  ->  backend url
internal_address  ->  backend address
headers_to_pass  ->  allowed headers
method  ->  request method
```
###### It will continue to be developed in line with the need.