package main

import (
	"github.com/mike-zenith/kmstool/src"
	"os"
)

func main() {
	a := kmstool.NewApp()
	a.Run(os.Args)
}
