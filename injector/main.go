package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type configuration struct {
	EtcdEndpoints []string      `envconfig:"etcd_endpoint" default:"localhost:2379,localhost:22379,localhost:32379"`
	Timeout       time.Duration `default:"5s"`
}

func main() {
	var s configuration
	err := envconfig.Process("glinda", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get the Working ID from the cli
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v path/to/zygofile", os.Args[0])

	}
	lisp, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Cannot read zygo file: ", err)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.EtcdEndpoints,
		DialTimeout: s.Timeout,
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	defer cli.Close()
	id := uuid.New()
	// Putting the content of the lisp file in the tuple space
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	_, err = cli.Put(ctx, id.String(), string(lisp))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

	//zygo.Repl(env, cfg)
}
