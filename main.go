package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
)


var (
	readStdin = flag.Bool("i", false, "read file from stdin")
	fileName = flag.String("f", "", "Go source filename")
	format = flag.String("format", "emacs", "output format [emacs]")
)

func main() {
	flag.Parse()
	
	var (
		src []byte
		err error
	)
	if *readStdin {
		src, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else if *fileName != "" {
		src, err = ioutil.ReadFile(*fileName)
                if err != nil {
                        panic(err)
                }
	} else {
		flag.Usage()
		return
	}
	items, err := index(src)
	if err != nil {
		panic(err)
	}
	output(items)
}

func output(items []item) {
	switch *format {
        case "emacs":
                for _, i := range items {
			fmt.Printf("%s,,%s,,%d\n", i.Title, i.Element, i.Pos) 
                }
        }
}
