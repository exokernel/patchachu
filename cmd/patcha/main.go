package main

import (
	"github.com/exokernel/patchachu/patchachu"
)

func main() {
	pdb := patchachu.NewPatchastore()
	pdb.Init()
}
