# gsheet-updater

```shell
export SPREADSHEET_ID="1WBAxWCxUQt9HXDWIbPXlQpKydYHvo2LUg5R4A-3d3LQ"
export FILE="hoursbytag.csv"
export TAB_ID="Sprint 25"

gsheet-updater
```

# Manual Release Building

```shell
git tag -a v1.0.0
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=v1.0.0"
```
