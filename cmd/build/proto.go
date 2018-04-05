package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/IngCr3at1on/x/build"
)

func init() {
	var out bytes.Buffer
	if err := build.Exec(&out, "which", "protoc"); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func regenerate(dir, path string) error {
	return build.Exec(
		nil,
		"protoc", "-I",
		fmt.Sprintf("%s/", dir),
		// FIXME: *.proto doesn't seem to work here...
		fmt.Sprintf("%s/carnival.proto", dir),
		fmt.Sprintf("--gogo_out=plugins=grpc:%s/", dir),
		fmt.Sprintf("--proto_path=%s:.", path),
	)
}
