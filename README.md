# MirageWyrm

MirageWyrm is a command-line tool for listing and fetching a random count of text files from  S3 buckets.

## Features

- List objects in S3 buckets recursively
- Filter listings by file extension
- Randomly fetch files that don't exist locally
- Configurable via command line flags or config file
- JSON or text logging with adjustable verbosity

## Installation

To install MirageWyrm, you need Go 1.21 or later. Use the following command:

```bash
go install github.com/gkwa/miragewyrm@latest
```

## Configuration

MirageWyrm can be configured through:

1. Command line flags
2. Environment variables
3. Configuration file (.miragewyrm.yaml)

The default configuration file location is `$HOME/.miragewyrm.yaml`. You can specify a different location using the `--config` flag.

### Configuration File Example

```yaml
bucket: my-s3-bucket
verbose: 1
log-format: json
```

## Usage

Here are some common usage examples:

```bash
# List files in default bucket
./miragewyrm list

# Use verbose logging
./miragewyrm list -v

# List files in a specific bucket
./miragewyrm list --bucket mybucket

# Fetch 10 random files from a specific bucket
./miragewyrm fetch --bucket mybucket random --count=10 -v -v

# Use JSON format logging
./miragewyrm --log-format json list
```

## Usage details

```log
root@ip-172-31-24-137:~# miragewyrm fetch random --count=3 -v -v
3:43AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:60 > Fetching manifest from S3... v=0
tree data
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:65 > Found objects in bucket total=71896 v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:68 > Filtered for missing files missing=71896 total=71896 v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:70 > Selecting random files available=71896 requested=3 v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:72 > Selected files for download count=3 v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:143 > Downloading file key=data/cache/epub/18055/pg18055.txt to=data/cache/epub/18055/pg18055.txt v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:169 > Downloaded file key=data/cache/epub/18055/pg18055.txt size_bytes=659297 to=data/cache/epub/18055/pg18055.txt v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:143 > Downloading file key=data/cache/epub/33888/pg33888.txt to=data/cache/epub/33888/pg33888.txt v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:169 > Downloaded file key=data/cache/epub/33888/pg33888.txt size_bytes=134174 to=data/cache/epub/33888/pg33888.txt v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:143 > Downloading file key=data/cache/epub/117/pg117.txt to=data/cache/epub/117/pg117.txt v=0
3:44AM INF ../home/runner/work/miragewyrm/miragewyrm/fetch/fetcher.go:169 > Downloaded file key=data/cache/epub/117/pg117.txt size_bytes=20734 to=data/cache/epub/117/pg117.txt v=0
root@ip-172-31-24-137:~# tree data
data
`-- cache
    `-- epub
        |-- 117
        |   `-- pg117.txt
        |-- 18055
        |   `-- pg18055.txt
        `-- 33888
            `-- pg33888.txt

6 directories, 3 files
root@ip-172-31-24-137:~#
```


### Global Flags

- `--config`: Config file path (default: $HOME/.miragewyrm.yaml)
- `--bucket`: S3 bucket name (default: streambox-helpfulferret)
- `-v, --verbose`: Increase verbosity (can be used multiple times)
- `--log-format`: Log format (text or json)

### Commands

#### `list`

Lists objects in an S3 bucket recursively, filtering by file extension.

```bash
./miragewyrm list [flags]
```

#### `fetch random`

Randomly selects and downloads files from S3 that don't exist locally.

```bash
./miragewyrm fetch random [flags]

Flags:
  -n, --count int        Number of random files to fetch (default 1)
  -o, --outdir string    Output directory for downloaded files (default ".")
```

#### `version`

Displays version information for the MirageWyrm binary.

```bash
./miragewyrm version
```

## Development

### Prerequisites

- Go 1.21 or later
- AWS credentials configured
- Access to an S3 bucket

### Building

```bash
make build
```

### Testing

```bash
go test ./...
```
