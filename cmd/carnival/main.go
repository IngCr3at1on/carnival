package main

import (
	"context"
	"fmt"
	"os"

	"github.com/IngCr3at1on/x/carnival"
	"github.com/IngCr3at1on/x/carnival/config"
)

func main() {
	cfg := config.Config{
		Address: ":44442",
	}

	if err := carnival.Start(context.Background(), cfg); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
