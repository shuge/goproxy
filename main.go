package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/elazarl/goproxy"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/shuge/goproxy/g"
)


const (
	defaultProfHTTP = "localhost:6060"
)


var (
	cfg string
	buildTimestamp string
	printBuildTimestamp bool

)

func main() {
	flag.BoolVar(&printBuildTimestamp, "v", false, "print build timestamp ")
	flag.StringVar(&cfg, "c", "", "full path to configuration file")
	flag.Parse()

	if printBuildTimestamp {
		fmt.Println(buildTimestamp)
		os.Exit(0)
	}

	if cfg == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := g.ParseConfig(cfg)
	if err != nil {
		log.Fatalln(fmt.Sprintf("[fatal] g.ParseConfig -%s- failed", cfg), err)
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	log.SetFlags(flags)
	if g.Cfgs.Debug {
		log.SetFlags(flags)
	}

	var output *lumberjack.Logger
	if g.Cfgs.Logpath != "" {
		output = &lumberjack.Logger{
			Filename:   g.Cfgs.Logpath,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		}
		log.SetOutput(output)
	}

	if g.Cfgs.Prof {
		go func() {
			var addr string
			if g.Cfgs.ProfHTTP == "" {
				addr = defaultProfHTTP
			} else {
				addr = g.Cfgs.ProfHTTP
			}

			err := http.ListenAndServe(addr, nil)
			if err != nil {
				log.Println(fmt.Sprintf("[error] http.ListenAndServe %s", addr), err)
			}
		}()
	}



	if g.Cfgs.Pidpath != "" {
		parent := path.Dir(g.Cfgs.Pidpath)
		_, err := os.Stat(parent)
		if err != nil && os.IsNotExist(err) {
			err := os.MkdirAll(parent, 0755)
			if err != nil {
				log.Println(fmt.Sprintf("[error] os.MkdirAll -%s- failed", parent))
			}
		}

		content := fmt.Sprintf("%#v", os.Getpid())
		err = ioutil.WriteFile(g.Cfgs.Pidpath, []byte(content), 0644)
		if err != nil {
			log.Println("[error] ioutil.WriteFile failed", err)
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

		if g.Cfgs.Pidpath != "" {
			err := os.Remove(g.Cfgs.Pidpath)
			if err != nil {
				log.Println(fmt.Sprintf("[error] os.Remove -%s-", g.Cfgs.Pidpath), err)
			}
		}

		os.Exit(0)
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = g.Cfgs.Debug
	proxy.Logger = log.New(output, "", log.LstdFlags)
	err = http.ListenAndServe(g.Cfgs.ListenHTTP, proxy)
	if err != nil {
		log.Println(fmt.Sprintf("[error] http.ListenAndServe -%s-", g.Cfgs.ListenHTTP), err)
	}

}
