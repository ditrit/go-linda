package main

import (
	"flag"
	"fmt"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	"os"
)

func usage(myflags *flag.FlagSet) {
	fmt.Printf("zygo command line help:\n")
	myflags.PrintDefaults()
	os.Exit(1)
}

func main() {
	cfg := zygo.NewGlispConfig("zygo")
	cfg.DefineFlags()
	err := cfg.Flags.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		usage(cfg.Flags)
	}

	if err != nil {
		panic(err)
	}
	err = cfg.ValidateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "zygo command line error: '%v'\n", err)
		usage(cfg.Flags)
	}
	lda := linda.Linda{}
	env := zygo.NewGlisp()
	env.AddFunction("in", lda.InRd)
	env.AddFunction("rd", lda.InRd)
	env.AddFunction("eval", lda.Eval)
	zygo.Repl(env, cfg)
}
