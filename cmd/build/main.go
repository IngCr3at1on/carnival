package main

import (
	"fmt"
	"os"

	"github.com/IngCr3at1on/x/build"
	"github.com/IngCr3at1on/x/build/wrapper"
)

var (
	_pkg = fmt.Sprintf("github.com/IngCr3at1on/x/carnival/cmd/carnival")

	_out string
)

func init() {
	_wd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_out := fmt.Sprintf("%s/out", _wd)
	if err := os.MkdirAll(_out, os.ModeDir|0700); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	wrapper.Builder = &build.Builder{
		Flows: []*build.Flow{
			&build.Flow{
				Name: "build",
				Steps: []func() error{
					func() error {
						return build.Exec(nil, "go", "build", "-o", fmt.Sprintf("%s/app", _out), _pkg)
					},
				},
			},
			&build.Flow{
				Name: "regenerate",
				Steps: []func() error{
					func() error {
						return regenerate("../../proto", "../../../..")
					},
				},
			},
		},
	}
}

func main() {
	if err := wrapper.Run(os.Args[1:]...); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
