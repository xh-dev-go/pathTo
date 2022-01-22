package main

import (
	_ "embed"
	"flag"
	"github.com/atotto/clipboard"
	"github.com/xh-dev-go/xhUtils/flagUtils"
	"github.com/xh-dev-go/xhUtils/flagUtils/FlagSets"
	"github.com/xh-dev-go/xhUtils/flagUtils/StringOptions"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagBool"
	"github.com/xh-dev-go/xhUtils/logical"
	"github.com/xh-dev-go/xhUtils/osDetection"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)


//go:embed version
var version string
func main() {
	showVersion := flagUtils.Version().BindCmd()

	fromOptions := StringOptions.New().
		Add("in-clipboard", "from clipboard").
		Add("in-stdin", "from std in").
		Bind(FlagSets.CommandFlag)

	clipboardOut := flagBool.New("out-clipboard", "copy result to clipboard").BindCmd()

	winFormat := flagBool.New("win", "convert path to windows format").BindCmd()
	unixFormat := flagBool.New("unix", "convert path to unix format").BindCmd()
	flag.Parse()

	if showVersion.Value() {
		println(version)
		os.Exit(0)
	}

	if !logical.ExclusiveOr(winFormat.Value(), unixFormat.Value()){
		panic("either win or unix is selectable, but not all")
	}

	option, err :=fromOptions.Value()
	if err == StringOptions.ExEmptySelection {
		panic("Nothing select")
	} else if err == StringOptions.ExDuplicateSelection {
		panic("Multiple option selected")
	} else if err != nil{
		panic(err)
	}

	var msg string
	switch option {
	case "in-clipboard":
		msg, err = clipboard.ReadAll()
		if err != nil {
			panic(err)
		}
	case "in-stdin":
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		msg = string(b)
	}

	var finalMsg=""

	if winFormat.Value() {
		if strings.HasPrefix(msg, "\\\\"){
			finalMsg = msg
		}else if strings.HasPrefix(msg, ".") {
			finalMsg = strings.ReplaceAll(msg, "/", "\\")
		} else if regexp.MustCompile(`^/[a-zA-Z]/$`).MatchString(msg[0:3]){
			finalMsg = strings.ReplaceAll(msg[1:2]+":\\"+msg[3:], "/", "\\")
		} else if strings.HasPrefix(msg, "~") {
			var home, err = os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			if osDetection.CurOS == osDetection.OS_WIN && strings.HasSuffix(home, "\\"){
				home = home[0: len(home)-1]
			}
			if osDetection.CurOS == osDetection.OS_LINUX && !strings.HasSuffix(home, "/"){
				home = home[0: len(home)-1]
			}
			finalMsg = home + strings.ReplaceAll(msg[1:], "/", "\\")
		} else {
			finalMsg = strings.ReplaceAll(msg, "/", "\\")
		}
	}

	if unixFormat.Value(){
		if strings.HasPrefix(msg, "\\\\"){
			finalMsg = strings.ReplaceAll(msg, "\\", "/")
		}else if strings.HasPrefix(msg, ".") {
			finalMsg = strings.ReplaceAll(msg, "\\", "/")
		} else if strings.HasPrefix(msg, "/") {
			finalMsg = msg
		} else {
			finalMsg = strings.ReplaceAll(msg, "\\", "/")
		}
	}

	if clipboardOut.Value(){
		clipboard.WriteAll(finalMsg)
	}
	println(finalMsg)
	return

}