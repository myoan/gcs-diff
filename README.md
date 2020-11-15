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
gcs-diff -b yoan-asset-stoage -cred ../gcs-client/lp-test-263613-a14c3443df95.json -src assets1 -dst assets2
---------------------------------------------
bucket:      yoan-asset-stoage
credential:  ../gcs-client/lp-test-263613-a14c3443df95.json
source:      assets1
destination: assets2
---------------------------------------------

storage: object doesn't exist: gs:/yoan-asset-stoage/assets1/only2.txt
storage: object doesn't exist: gs:/yoan-asset-stoage/assets2/only1.txt
storage: Object doesn't match: gs:/yoan-asset-stoage/assets1/samename.txt
```