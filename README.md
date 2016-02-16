# Ape - Munki API and server

## API 
documented at ape-docs.groob.io

## Server
The Ape server is a single binary that includes the API and also supports a static file server to serve the munki repo 

## Release

OS And Linux binaries are available on the [release page](https://github.com/groob/ape/releases/)

## Docker

```bash
docker pull groob/ape:latest
docker run -d --name munki -v /path/to/repo:/data -p 80:80 groob/ape
```
