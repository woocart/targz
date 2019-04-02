# Dumb(but still safe) directory to tar.gz writer

Like tar but without annoying `tar:: file changed as we read it`. See https://bugzilla.redhat.com/show_bug.cgi?id=1058526#c9 for explanation.

This app aborts file compression if consistency check fails, but not if only `atime or permission changed`.

## Running tests

```make build lint test```

## Building Docker image

```make image```

## Making a release

```make release```
