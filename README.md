## Installation
```shell
go install github.com/xh-dev-go/pathTo@latest
```

## Usage
Convert from stdin to windows format and copy result to clipboard
```shell
# convert ~/abc to c:\users\{user}\abc
# and the result copy to clipboard
echo ~/abc | pathTo -win -from-stdin -outclipboard
```

Convert from clipboard to unix format
```shell
# read clipboard to get the path
# convert path as c:\users\{user}\abc to /c/users/abc
echo c:\\users\\user\\abc | pathTo -unix -from-clipboard
```
