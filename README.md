## Build

```sh
rsrc -manifest sav.manifest -o rsrc.syso
go build -ldflags='-s -w -H windowsgui' -trimpath
```
