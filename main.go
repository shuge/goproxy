package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

var (
	version        bool
	verbose        bool
	pid            string
	buildTimestamp string
	addr           string
)

func main() {
	flag.BoolVar(&verbose, "verbose", false, "print proxy request to stdout")
	flag.BoolVar(&verbose, "version", false, "print version")
	flag.StringVar(&addr, "addr", ":8118", "proxy listen address")
	flag.StringVar(&pid, "pid", "", "full path to pid")
	flag.Parse()

	if version {
		fmt.Println(buildTimestamp)
		os.Exit(0)
	}

	if pid != "" {
		parent := path.Dir(pid)
		_, err := os.Stat(parent)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(parent, 0755)
		}

		content := fmt.Sprintf("%#v", os.Getpid())
		err = ioutil.WriteFile(pid, []byte(content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		//syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigs

		if pid != "" {
			err := os.Remove(pid)
			if err != nil {
				log.Println("os.Remove failed")
			}
		}

		os.Exit(0)
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	err := http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Println(err)
	}

}
