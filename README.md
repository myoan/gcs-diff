# gcs-diff

Simple diff viewer for Google Cloud Storage

## install

Install from https://github.com/myoan/gcs-diff/releases

## usage

```
❯ gcs-diff -help
Usage of gcs-diff:
  -b string
        backetname
  -conc int
        upload cuncurrency (default 4)
  -cred string
        credential path
  -dst string
        dst path to compare
  -src string
        src path to compare
```

example
```
❯ gcs-diff -b example-storage -cred credential.json -src assets1 -dst assets2
---------------------------------------------
bucket:      example-storage
credential:  credenital.json
source:      assets1
destination: assets2
---------------------------------------------

+ gs:/example-storage/assets1/only2.txt
- gs:/example-storage/assets2/only1.txt
~ gs:/example-storage/assets1/samename.txt
```

- `+`: means does not exist in src dir, but exists dst dir
- `-`: means exists in src dir, but does not exist dst dir
- `~`: means exist both dir, but CRC32C is changed
- `(not output)`: means exist both dir, and same CEC32C