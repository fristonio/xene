## xene apiserver

Run xene apiserver.

### Synopsis

Run xene apiserver which can then be used to communicate to user facing interface of xene.

```
xene apiserver [flags]
```

### Options

```
  -l, --cert-file string           Certificate to use for the API Server when running under HTTPS scheme.
  -n, --disable-auth               If the authentication should be disabled for the API server.
  -h, --help                       help for apiserver
  -b, --host string                Host to bind the api server to. (default "0.0.0.0")
  -j, --jwt-secret string          JWT secret for authentication purposes, make sure it is secure and non bruteforcable.
  -k, --key-file string            Key to use when using HTTPS scheme for the server.
  -p, --port uint32                Port to bind the xene api server on. (default 6060)
  -s, --scheme string              Scheme to use for the api server. (default "http")
  -z, --standalone                 Run xene apiserver in standalone mode.
  -d, --storage-directory string   Storage directory to use for xene apiserver. (default "/var/run/xene/store")
  -e, --storage-engine string      Storage engine to use for the API server (default "badger")
  -u, --unix-socket string         Default path for the unix domain socket, when using unix scheme (default "/var/run/xene/xene.sock")
  -v, --verbose-logs               Print verbose APIServer request logs.
```

### SEE ALSO

* [xene](xene.md)	 - xene is an open source workflow builder and executor tool.

###### Auto generated by spf13/cobra on 25-May-2020
