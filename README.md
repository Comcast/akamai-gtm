# akamai-gtm

A Golang-based CLI for Akamai [GTM](https://developer.akamai.com/api/luna/config-gtm/overview.html).

## Installation

Download the desired [release](https://github.com/Comcast/akamai-gtm/releases) version for your operating system. Untar and install the `akamai-gtm` executable
to your `$PATH`.

### Compiling from Golang source

Alternatively, if you choose to compile from Golang source code:

* install Golang
* set up your `$GOPATH`
* clone `comcast/akamai-gtm` to `$GOPATH/src/github.com/comcast/akamai-gtm`
* cd `$GOPATH/src/github.com/comcast/akamai-gtm && make`

## Usage

```
akamai-gtm --help

NAME:
   akamai-gtm - A CLI to Akamai GTM configuration

USAGE:
   akamai-gtm [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
    domains                     domains
    domain                      domain <domain.akadns.net>
    domain-create               domain-create --type <domainType> <domain.akadns.net>
    domain-update               domain-update --json <DomainJSONFile>
    data-centers                data-centers <domain.akadns.net>
    data-centers-delete         data-centers-celete --id <dataCenterId> --id <dataCenterId> <domain.akadns.net>
    data-centers-delete-all     data-centers-delete-all <domain.akadns.net>
    data-center                 data-center --id <dataCenterId> <domain.akadns.net>
    data-center-create          data-center-create --json <DataCenterJSONFile> <domain.akadns.net>
    data-center-update          data-center-update --json <DataCenterJSONFile> <domain.akadns.net>
    data-center-delete          data-center-delete --id <dataCenterId> <domain.akadns.net>
    properties                  properties
    properties-delete           properties-delete --names <PropertyName>,<PropertyName> <domain.akadns.net>
    properties-delete-all       properties-delete-all <domain.akadns.net>
    property                    property --name <PropertyName> <domain.akadns.net>
    property-create             property-create --json <PropertyJSONFile> <domain.akadns.net>
    property-update             property-update --json <PropertyJSONFile> <domain.akadns.net>
    property-delete             property-delete --name <PropertyName> <domain.akadns.net>
    traffic-targets             traffic-targets --name <PropertyName> <domain.akadns.net>
    liveness-tests              liveness-tests --name <PropertyName> <domain.akadns.net>
    status                      status <domain.akadns.net>

GLOBAL OPTIONS:
   --host value                         Luna API Hostname [$AKAMAI_EDGEGRID_HOST]
   --client_token value, --ct value     Luna API Client Token [$AKAMAI_EDGEGRID_CLIENT_TOKEN]
   --access_token value, --at value     Luna API Access Token [$AKAMAI_EDGEGRID_ACCESS_TOKEN]
   --client_secret value, -s value      Luna API Client Secret [$AKAMAI_EDGEGRID_CLIENT_SECRET]
   --help, -h                           show help
   --version, -v                        print the version
```
