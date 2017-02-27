package main

import (
	"flag"
	"fmt"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	"log"
	"math/rand"
	"os"
	"time"
)

func usage(myflags *flag.FlagSet) {
	fmt.Printf("zygo command line help:\n")
	myflags.PrintDefaults()
	os.Exit(1)
}

func sleep(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	t := args[0].(*zygo.SexpInt)
	time.Sleep(time.Duration(rand.Int31n(int32(t.Val))) * time.Millisecond)
	return zygo.SexpNull, nil
}

func main() {
	cfg := zygo.NewGlispConfig("zygo")
	cfg.DefineFlags()
	err := cfg.Flags.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		usage(cfg.Flags)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err != nil {
		panic(err)
	}
	err = cfg.ValidateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "zygo command line error: '%v'\n", err)
		usage(cfg.Flags)
	}
	lda := linda.New(nil)
	env := zygo.NewGlisp()
	env.AddFunction("in", lda.InRd)
	env.AddFunction("rd", lda.InRd)
	env.AddFunction("out", lda.Out)
	env.AddFunction("eval", lda.Eval)
	env.AddFunction("sleep", sleep)
	env.SourceFile(f)
	_, err = env.Run()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	<-done
	//zygo.Repl(env, cfg)
}
