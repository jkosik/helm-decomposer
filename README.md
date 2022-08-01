## helm-decomposer
The tool takes Helm package (.tgz) as an input and visualizes hierarchy of dependencies (Helm subcharts and correspondig images) for further analysis.

## Usage
Currently just WIP!
- `helm pull bitnami/nginx`
- Update `chartPath` in the code accordingly.
- `go mod init github.com/jkosik/helm-decomposer`
- `go mod tidy`
- `go run main.go`
