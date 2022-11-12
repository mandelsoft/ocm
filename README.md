# SAP R/3 connection demo plugin for the Open Component Model

**OCM Component**: `sap.com/r3trans/ocmplugin` \
**Resource Name**: `r3trans` \
**Delivery Repository**: `ghcr.io/mandelsoft/cnudie` 

## Description

This plugin simulates a connection of Open Component Model to
the SAP R/3 System.

It provides an access method to access R/3 transports in the transport system
and a repository uploader, which can be used to upload transport requests
described as resources in a component version directly into a transport
environment during he import of a component version into a local
repository landscape.

## Plugin Characteristics

```
Plugin Name:      r3trans.sap.com
Plugin Version:   0.0.1
Path:             /home/d021770/.ocm/plugins/r3trans.sap.com
Capabilities:     Access Methods, Repository Uploaders
Source:
  Component:       sap.com/r3trans/ocmplugin
  Version:         0.0.1-dev
  Resource:        r3trans
  Repository:      OCIRegistry/mandelsoft/cnudie
    Specification: {"baseUrl":"ghcr.io","componentNameMapping":"urlPath","subPath":"mandelsoft/cnudie","type":"OCIRegistry"}
Description: 
      This plugin provided support to access R/3 transport requests and to upload
      and to upload them to a transport environment again.
      The plugin uses the following configuration fields:
      - «systems» *map[string]<config>*
        The configuration used  fr a set of transport systems:
        - «path» *string* (default «/tmp/r3trans»)
          The base address to be used for the transport system.

Access Methods:
- Name; r3trans.sap.com
    demo access R/3 transport files
  Versions:
  - Version: v1
    The type specific specification fields are:
    - «transport» *string*
      name of transport request.
    - «transportSystem» *string*
      address of transport system
    - «path» *string*
      sub path in transport system.
    Command Line Options:
      - «accessTransport»: [*string*] name of R/3 transport request
      - «accessTransportPath»: [*string*] path in transport system
      - «accessTransportSystem»: [*string*] R/3 transport system
Repository Uploaders:
- Name: r3trans
  Upload R/3 transport requests to the transport system.
  It uses the following target specification fields:
  - «type</code *string* constant «r3trans/v1»
  - «transportSystem» *string*
    The address of the R/3 transport system.
  - «path» *string*
    The sub path used in the transport system.
  Registration Contraints:
  - Artefact Type: r3trans.sap.com/transportRequest
    Media Type   :
```


## Playground

This project is a demo, only. It can be build as part of the ocm project.
Clone the `r3trans/master` branch into you ocm project under folder `local/r3trans`.

Now you can use the Makefile in this folder to build and push he plugin.
Set the variable `OCMREPO` to configure your own OCM respository.


