## helm-decomposer
The tool templates the Helm package (.tgz or untarred folder) identifies all images in use and visualizes tree structure of the Chart.

## Demo
![](readme/readme.mp4)

## Build the binary
```
git clone git@github.com:jkosik/helm-decomposer.git
cd helm-decomposer
go build .
```

## Usage
1. Download any Helm Chart. You will reference it later on.
2. Run as `./helm-decomposer -chart mychart.tgz -o`
```
❯ ./helm-decomposer -h
Usage of ./helm-decomposer:
  -chart string     
        Helm Chart to process. Submit .tgz or folder name. (default "sample-helm-charts/nginx")
  -ij   Write image list to images.json. (default "false")
  -iy   Write image list to images.yaml. (default "false")
  -o    Write Helm Chart tree to helm-decomposer-output.md. (default "false")
```

## Issues
- Processed Helm Chart must have all variables properly set to allow helm-decomposer properly template the whole Helm Chart and Subcharts.
- Edge case appears when Helm chart uses dependency aliases combined with subchart parametrized on parent level only.

## License
MIT
