package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/coreos/etcd/clientv3"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	//"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"os"
	"time"
)

const prefix = "LINDA"

type configuration struct {
	EtcdEndpoints []string      `envconfig:"etcd_endpoint" default:"localhost:2379,localhost:22379,localhost:32379"`
	Timeout       time.Duration `default:"5s"`
}

func sleep(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	t := args[0].(*zygo.SexpInt)
	time.Sleep(time.Duration(rand.Int31n(int32(t.Val))) * time.Millisecond)
	return zygo.SexpNull, nil
}

func main() {
	var s configuration
	err := envconfig.Process("glinda", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get the Working ID from the cli
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v WORK_ID", os.Args[0])

	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.EtcdEndpoints,
		DialTimeout: s.Timeout,
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	defer cli.Close()

	// The tuple does not exists, watch for a new event until a tuple match
	var msg linda.Message
	rch := cli.Watch(context.Background(), "LINDA-evalc-", clientv3.WithPrefix())
Loop:
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// TODO: Check if ev.Type is PUT otherwise continue
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			// TODO: Remove the tuple from the space is IN is called
			_, err := cli.Delete(context.TODO(), string(ev.Kv.Key), clientv3.WithPrefix())
			if err != nil {
				log.Println(err)
				continue
			}
			// Try to decode the message
			buf := bytes.NewBuffer(ev.Kv.Value)
			dec := gob.NewDecoder(buf)
			err = dec.Decode(&msg)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}
			break Loop
		}
	}
	env := msg.Env
	env.Run()
	/*
		fn := msg.Args[0].(*zygo.SexpFunction)
		expr, err := env.Apply(fn, msg.Args[1:]) // Put the result in the tuplespace
		if err != nil {
			// TODO
			log.Println(err)
		}
		if expr != zygo.SexpNull {
			_, err := cli.Put(context.TODO(), prefix+"-"+uuid.New().String(), expr.SexpString(&zygo.PrintState{}))
			log.Println(err)
			return
		}
	*/
	done := make(chan bool)
	<-done
	//zygo.Repl(env, cfg)
}
