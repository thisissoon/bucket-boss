# bucket-boss

CLI tool for managing bucket contents, currently only support AWS S3 but designed to support different providers like GCS


## Development

 - Go 1.11+
 - Dependencies managed with `go mod`

### Setup

These steps will describe how to setup this project for active development. Adjust paths to your desire.

1. Clone the repository: `git clone github.com/thisissoon/bucket-boss bucket-boss`
2. Build: `make build`
3. 🍻

### Dependencies

Dependencies are managed using `go mod` (introduced in 1.11), their versions
are tracked in `go.mod`.

To add a dependency:
```
go get url/to/origin
```

### Configuration

Configuration can be provided through a toml file, these are loaded
in order from:

- `/etc/bucket-boss/bucket-boss.toml`
- `$HOME/.config/bucket-boss.toml`

Alternatively a config file path can be provided through the
-c/--config CLI flag.

#### Example bucket-boss.toml
```toml
[log]
console = true
level = "debug"  # [debug|info|error]

[aws]
enabled = true
bucketName = "bucketname" # name of the bucket you want to manage
region = "region" # bucket region
accessKey = "AWS_ACCESS_KEY" # AWS access key usually set as an environment variable
secretKey = "AWS_SECRET_ACCESS_KEY" # AWS secret key usually set as an environment variable
```
