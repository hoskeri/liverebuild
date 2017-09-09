# LiveRebuild: powerful and flexible livereload server.

## Config file

Example Config:

```ini

# watch files and build.
buildCommandRoot = ./src
buildFiles = *.elm api/*.py # paths are relative to buildCommandRoot
buildCommand = "make"

# live reload patterns
watchServeRoot = ./dist/static # paths are relative to watchServeRoot
watchServeFiles = *.js *.css
watchFallback = index.html # render on 404

# customize ports and logging
listenLivereload = :37529
listenStatic = :4000
verbose = true

```
