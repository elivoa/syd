package main

import (
	"github.com/elivoa/syd"
	"flag"
	"fmt"
	"github.com/elivoa/got"
	"github.com/elivoa/got/config"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("\n >> start syd project <<")
	config.Config.RegisterModule(syd.Module)
	got.BuildStart()
}
