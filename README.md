# build-struct

This is a simple command line tool which converts a given configuration file into a Go struct. It accepts a configuration file in JSON or YAML format. The tool looks at the content within the configuration file, not necessarily the file extension. I designed it this way as sometimes I use `config` instead of `config.yaml` so I wanted it to be flexible.

## Use Case

For my personal use case, I have structs which have a lot of fields and I have to manually create them. This tool helps me to automate the process of creating structs.

## Installation

```bash
go install github.com/mitchs-dev/build-struct@latest
```

## Usage

```bash
build-struct <struct-name> <path-to-config-file>
```

## Example

I have a configuration file (`config.yaml`) with the following content:

```yaml
apiVersion: v1
metadata:
  name: my-app
  namespace: default
settings:
  logging: 
    debug: false
    format: json
  database:
    host: localhost
    port: 5432
    username: root
    password: password
  data:
  - path: /var/data
    threshold: 0.80
    enable: true
    size: 100Gi
  - path: /var/log
    threshold: 0.90
    enable: true
    size: 10Gi
```

I can run the following command:

```bash

build-struct Config config.yaml
```

This will generate the following Go struct:

```go
type Config struct {
        ApiVersion string `yaml:"apiVersion"`
        Metadata struct {
                Name string `yaml:"name"`
                Namespace string `yaml:"namespace"`
        } `yaml:"metadata"`
        Settings struct {
                Logging struct {
                        Debug bool `yaml:"debug"`
                        Format string `yaml:"format"`
                } `yaml:"logging"`
                Database struct {
                        Password string `yaml:"password"`
                        Host string `yaml:"host"`
                        Port int `yaml:"port"`
                        Username string `yaml:"username"`
                } `yaml:"database"`
                Data []struct {
                        Threshold float64 `yaml:"threshold"`
                        Enable bool `yaml:"enable"`
                        Size string `yaml:"size"`
                        Path string `yaml:"path"`
                } `yaml:"data"`
        } `yaml:"settings"`
}
```
