## ocm install plugins &mdash; Get Plugins

### Synopsis

```
ocm install plugins [<options>] <component version> <plugin name>
```

### Options

```
  -d, --describe   describe plugin, only
  -f, --force      overwrite existing plugin
  -h, --help       help for plugins
```

### Description


Download and install a plugin provided by an OCM component version.

If the <code>--repo</code> option is specified, the given names are interpreted
relative to the specified repository using the syntax

<center>
    <pre>&lt;component>[:&lt;version>]</pre>
</center>

If no <code>--repo</code> option is specified the given names are interpreted 
as located OCM component version references:

<center>
    <pre>[&lt;repo type>::]&lt;host>[:&lt;port>][/&lt;base path>]//&lt;component>[:&lt;version>]</pre>
</center>

Additionally there is a variant to denote common transport archives
and general repository specifications

<center>
    <pre>[&lt;repo type>::]&lt;filepath>|&lt;spec json>[//&lt;component>[:&lt;version>]]</pre>
</center>

The <code>--repo</code> option takes an OCM repository specification:

<center>
    <pre>[&lt;repo type>::]&lt;configured name>|&lt;file path>|&lt;spec json></pre>
</center>

For the *Common Transport Format* the types <code>directory</code>,
<code>tar</code> or <code>tgz</code> is possible.

Using the JSON variant any repository type supported by the 
linked library can be used:

Dedicated OCM repository types:
- `ComponentArchive`

OCI Repository types (using standard component repository to OCI mapping):
- `ArtefactSet`
- `CommonTransportFormat`
- `DockerDaemon`
- `Empty`
- `OCIRegistry`
- `oci`
- `ociRegistry`


### Examples

```
$ ocm install plugin ghcr.io/github.com/mandelsoft/cnudie//github.com/mandelsoft/ocmplugin:0.1.0-dev
```

### SEE ALSO

##### Parents

* [ocm install](ocm_install.md)	 &mdash; Install elements.
* [ocm](ocm.md)	 &mdash; Open Component Model command line client

