# gotree
Create file structure trees in terminals

```
.
├── cmd/
│   ├── list.go
│   ├── root.go
│   └── tree.go
├── empty/
├── files/
│   ├── getFiles.go
│   └── getFiles_test.go
├── tree/
│   └── print.go
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── gotree
├── main.go
└── tree_output.txt
```


# Usage

```
// list directore
gotree list
// show hidden files and directories
gotree list -a
// specify a max depth to parse
gotree list --maxdepth 3

// get filestructure and output as tree
gotree tree
// pipe into file
gotree tree > tree_output.txt 2>&1

```

# Run the test
```
go test ./files -v
```

# build and test locall
```
cd gotree

// build a local binary to run and test
go build -o gotree
./gotree list

// install it globally
go install .

```

