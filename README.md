# A Demo Using the r3trans OCM plugin.

This demo project builds a demo component
refering to a faked transport request.

The component version just contains a single resource
called `transport`, which describes the transport
request.

This component can then be transferred (by value) into
another environment (the example just uses another CTF).
Hereby, an automatic import into a faked transport system
can be done.

The access method of the transport resource will be adapted
accordingly, and the transport content will be visible in the
target transport environment.

## Steps

All operations are preconfigured as make targets.


- Setup your environment to use the plugin:

  Set your OCM repository according to your needs. The default
  delivery repository is `ghcr.io//mandelsoft/cnudie`.

  ```bash
  $ OCMREPO=ghcr.io//mandelsoft/cnudie
  $ install plugin -f $OCMREPO//sap.com/r3trans/ocmplugin 
  ```

  With make the command would look like:

  ```bash
  $ make setup
  ```

  Afterwards the plugin capabilities can be examined with

  ```bash
  $ ocm describe plugin r3trans
  ```

  After the plugin has been installed, a new access method
  `r3trans.sap.com/v1` will be available.

  ```bash
  ocm help accessmethods
  ```
  
  will now show the new method specification and the new 
  command line options usable to compose an appropriate
  access specification.


- Create your transport and provide a component version

  ```bash
  $ make ctf
  ```
  
  It creates a demo transport request BINK08154711 in
  `/tmp/r3trans/source`.

  ```bash
  $ tree /tmp/r3trans
  ```

- Now the content of the transport request resource can be downloaded
  from the local source

  ```bash
  $ ocm download resource gen/ctf -O local/transport.local transport
  $ cat local/transport.local
  ```

  or with make

  ```bash
  $ make local-download
  ```

- Transport the component into a foreign environment

  ```bash
  $ ocm transfer ctf -V gen/ctf local/target
  ```
  
  or with make

  ```bash
  $ make local-transport
  ```

  The content is now contained as local blob in the target CTF


  ```bash
  $ ocm get resources local/target -o yaml
  ```

- Download the locally stored transprt request content

  It can be downloaded as usual:

  ```bash
  $ ocm download resource local/target -O local/transport.target transport
  $ cat local/transport.target
  ```

  or with make

  ```bash
  $ make target-download
  ```

- Transport with automatic upload into target transport environment

  Alternatively the transport can be requested to upload all transport requests
  into a a target transport environment.

  ```bash
  $ transfer ctf -V --uploader plugin/r3trans.sap.com=@importtarget.yaml gen/ctf local/target
  ```

  The file `importtarget.yaml`described the transport environment to use for the uplpoad.
  
  ```
    type: r3trans/v1
    transportSystem: target
    path: TST
  ```

  or with make

  ```bash
  $ make local-import
  ```

  As a result, the transport request will be stored in the target transport environmemt:

  ```bash
  $ tree /tmp/r3trans/target
  ```

  The access specification used in the local OCM repository is adapted accordingly:

  ```bash
  $ ocm get resources local/target -o yaml
  ```

- Download the transport request from the local transport environment

  ```bash
  $ ocm download resource local/target -O local/transport.target transport
  ```

  or with make

  ```bash
  $ make target-download
  ```

  The modified access specification of the component version is now
  used to access the content directly in the target system.
