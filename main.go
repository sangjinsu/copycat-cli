package main

import (
	"flag"
	"fmt"
	"github.com/copycatcli/utils"
	"log"
)

func main() {

	var src, dst, replace string

	// 입력값 처리
	flag.StringVar(&src, "src", "", "source directory")
	flag.StringVar(&dst, "dst", "", "destination directory")
	flag.StringVar(&replace, "replace", "", "replace string")

	flag.Parse()

	if src == "" || dst == "" || replace == "" {
		fmt.Println("Usage: go run main.go -src=[source directory] -dst=[destination directory] -replace={{old string}:{new string}}")
		return
	}

	replacements := utils.ParseReplacements(replace)
	err := utils.CopyDir(src, dst, replacements)
	if err != nil {
		log.Fatalln(err)
	}
}
