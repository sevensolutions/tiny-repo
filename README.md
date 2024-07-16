# TinyRepo

TinyRepo is a very basic repository server for storing versioned artifacts.
Artifacts are simply treated as blobs, so you can version anything. I'am versioning ZIP files.
The storage backend can be a local directory or an S3 compatible service like MinIO.

## Artifact Versions

Every artifact has a artifact name and is part of a namespace.

Artifacts are versioned using [Semantic Versioning(Semver)](https://semver.org/).

## Configuration

Configuration is done using environment variables.

TODO

## Authentication

TODO

## HTTP API

### Push an Artifact (Upload)
```
PUT http://localhost:8080/:namespace/:name/:version[?keep=3]
```

This endpoint is used to push a new artifact.

You may optionally specify a `keep`-parameter to automatically delete old version. The default is 0, which means all versions will be kept.
A value of 1 will keep just 1 version, including the one currently beeing pushed.

The request must be sent as multipart-form-data containing a single field called `file`.

Example cURL call:

```bash
curl -X PUT -H "Authorization: Bearer {token}" -F file='@filename' http://localhost:8080/foo/bar/1.0.0
```

### Pull an Artifact (Download)

```
GET http://localhost:8080/:namespace/:name/:version|latest[?name={targetFileName}]
```

This endpoint is used to download an artifact of a specific version.
You may use the `latest` keyword to download the latest version.

An optional target file name may be supplied using the `name`parameter to specify the name, which will be used as the filename of the attachment.
Otherwise the file is just called *blob*.

### List all Versions

```
GET http://localhost:8080/:namespace/:name
```

This endpoint returns a JSON, containing all available versions of the artifact.

Here is an example:

```json
{
  "count": 2,
  "latest": "1.3.17",
  "versions": [
    "1.3.17",
    "1.3.16"
  ]
}
```

### Deleting a Version

This is not supported at the moment.

## Disclaimer

This is one of my first "bigger" projects, written in GO, so I'am pretty sure a lot of things are wrong 🙈
