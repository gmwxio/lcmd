package main

import (
	"fmt"
	"os"

	"github.com/jpillora/opts"

	"github.com/wxio/lcmd/internal/grpcsvc"
	"github.com/wxio/lcmd/internal/roll"
)

// Version, Data and Commit set by compile
var (
	Version string
	Date    string
	Commit  string
)

type root struct {
	opts opts.ParsedOpts
}

func (rt *root) Run() {
	fmt.Printf("%s\n%s\n", os.Args[0], rt.opts.Help())
}

func main() {
	rt := &root{}
	ro := opts.New(rt).
		Version(Version).
		EmbedGlobalFlagSet().
		Complete()
	roll.Register(ro)
	grpcsvc.Register(ro)
	rt.opts = ro.Parse()
	rt.opts.RunFatal()
}
