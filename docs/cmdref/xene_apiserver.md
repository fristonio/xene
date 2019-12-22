## xene apiserver

Run xene apiserver.

### Synopsis

Run xene apiserver which can then be used to communicate to user facing interface of xene.

```
xene apiserver [flags]
```

### Options

```
  -l, --cert-file string     Certificate to use for the API Server when running under HTTPS scheme.
  -c, --config string        Config file for API server.
  -h, --help                 help for apiserver
  -b, --host string          Host to bind the api server to. (default "0.0.0.0")
  -k, --key-file string      Key to use when using HTTPS scheme for the server.
  -p, --port uint32          Port to bind the xene api server on. (default 6060)
  -s, --scheme string        Scheme to use for the api server. (default "http")
  -u, --unix-socket string   Default path for the unix domain socket, when using unix scheme (default "/var/run/xene/xene.sock")
```

### SEE ALSO

* [xene](xene.md)	 - xene is an open source workflow builder and executor tool.

###### Auto generated by spf13/cobra on 22-Dec-2019