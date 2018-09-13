# LiveRebuild: powerful and flexible livereload server.

## Config file

Example Config:

```toml
verbose = false

[static]
address = ":4000"
root = "./build/static"
paths = [ "./build/static/*" ]
fallback = "index.html"

[lr]
address = ":37529"
paths = [ "./build/static/*" ]

[build]
paths = [ "src/*.elm", "src/*.tmpl"]
cmd = [ "make" ]

[daemon]
paths = [ "./build/server/*" ]
cmd = [ "sleep", "1d" ]
```
