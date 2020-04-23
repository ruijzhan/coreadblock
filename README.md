# coreadblock
A CoreDNS plugin to block ads

Issue of this branch:

Out and Setup are exported in export_test.go file, but this does not work. To be resolved later. 

```shell script
$ go test ./test/ -v
# github.com/ruijzhan/coreadblock/test [github.com/ruijzhan/coreadblock/test.test]
test\coreadblock_test.go:19:2: undefined: coreadblock.Out
test\setup_test.go:12:12: undefined: coreadblock.Setup
FAIL    github.com/ruijzhan/coreadblock/test [build failed]
FAIL
```