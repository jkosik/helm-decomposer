## helm-decomposer
The tool templates the Helm package (.tgz or untarred folder) identifies all images in use and visualizes tree structure of the Chart and all dependencies.

## Usage
- `helm pull bitnami/nginx`
- `go mod init github.com/jkosik/helm-decomposer`
- `go mod tidy`
- `go run main.go`

## TODO
- Submitted Helm Chart must be healthy, i.e. Helm templating must end up without any warnings. Edge case hit when Helm chart uses dependency aliases combined with subchart parametrized on parent level only.
- flags for output files and help
- including images into the visal chart hierarchy 
