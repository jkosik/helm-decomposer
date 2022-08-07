## helm-decomposer
The tool templates the Helm package (.tgz or untarred folder) identifies all images in use and visualizes tree structure of the Chart and all dependencies (aliased dependencies are merged).

## Usage
- Download any Helm Chart. You will reference it later on.
- Run as `helm-decomposer -chart ./mychart.tar.gz
```
❯ ./helm-decomposer -h
Usage of ./helm-decomposer:
  -chart string
        Helm Chart to process. Submit tar.gz or folder name. (default "samples/nginx")
  -i    Inspect images used in the Helm Chart. (default "false")
  -o    Write output to helm-decomposer-output.md. (default "false")
```

## Compile
- `go mod init github.com/jkosik/helm-decomposer`
- `go mod tidy`
- `go build .`

## TODO
- Submitted Helm Chart must be healthy, i.e. Helm templating must end up without any warnings. Edge case hit when Helm chart uses dependency aliases combined with subchart parametrized on parent level only.
- including images into the visal chart hierarchy 


## Notes
- templating works on all levels
- aliased deps not working with Engine


