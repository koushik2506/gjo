//
// gjo - An attempt of a go-version of jo
// author - koushik.narayanan@gmail.com
//

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Version struct {
	Program string `json:"program"`
	Author  string `json:"author"`
	Repo    string `json:"repo"`
	Version string `json:"version"`
}

var version Version = Version{Program: "gjo", Author: "nkoushik", Repo: "https://github.com/koushik2506/gjo", Version: "0.1"}

var out io.Writer = os.Stdout
var err io.Writer = os.Stderr
var in io.Reader = os.Stdin

var arrayflag bool
var boolflag bool
var verflag bool
var jsonverflag bool
var prettyflag bool

func init() {
	flag.BoolVar(&arrayflag, "a", false, "creates an array of words")
	flag.BoolVar(&boolflag, "B", false, "disables boolean true/false")
	flag.BoolVar(&verflag, "v", false, "Show version")
	flag.BoolVar(&jsonverflag, "V", false, "Show version in JSON")
	flag.BoolVar(&prettyflag, "p", false, "pretty-prints JSON on output")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage: jo [-a] [-B] [-p] [-v] [-V] [word...]
word is a key=value or key@value`)
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		flag.PrintDefaults()
	}

}

func GuessValue(input string) interface{} {
	// try int
	i, err := strconv.ParseInt(input, 0, 64)
	if err == nil {
		return i
	}

	// try bool first if -B is not set
	if !boolflag {
		b, err := strconv.ParseBool(input)
		if err == nil {
			return b
		}
	}

	//try float
	f, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return f
	}

	// Empty strings are a short cut to null
	if len(input) == 0 {
		return nil
	}

	//give up
	return input
}

func ToJson(input []string) ([]byte, error) {
	if arrayflag {
		return JsonArrayEncode(input)
	} else {
		return JsonKeyValEncode(input)
	}
}

func JsonArrayEncode(input []string) ([]byte, error) {
	inputarray := make([]interface{}, len(input), len(input))

	for i, args := range input {
		inputarray[i] = GuessValue(args)
	}

	return JsonEncode(inputarray)
}

func JsonKeyValEncode(input []string) ([]byte, error) {
	inputmap := make(map[string]interface{})
	for _, args := range input {
		kvArray := strings.Split(args, "=")
		if len(kvArray) == 2 {
			inputmap[kvArray[0]] = GuessValue(kvArray[1])
		} else if len(kvArray) == 1 {
			inputmap[kvArray[0]] = nil
		} else {
			return nil, fmt.Errorf("Invalid key/value combination")
		}
	}

	return JsonEncode(inputmap)
}

func JsonEncode(input interface{}) ([]byte, error) {
	if prettyflag {
		data, err := json.MarshalIndent(input, "", "   ")
		return data, err
	} else {
		data, err := json.Marshal(input)
		return data, err
	}
}

func ReadInput() []string {
	scanner := bufio.NewScanner(in)

	var inputstr []string

	for scanner.Scan() {
		inputstr = append(inputstr, scanner.Text())
	}

	return inputstr
}

func main() {
	flag.Parse()

	if verflag {
		fmt.Fprintf(out, "%s %s\n", version.Program, version.Version)
		return
	}

	if jsonverflag {
		data, err := JsonEncode(version)
		if err != nil {
			log.Fatalf("JSON Marshalling failed: %s", err)
			return
		}
		fmt.Fprintf(out, "%s\n", data)
		return
	}

	var input []string

	if flag.NArg() > 0 {
		input = flag.Args()
	} else {
		input = ReadInput()
	}

	var data []byte

	data, err := ToJson(input)

	if err != nil {
		log.Fatalf("JSON Marshalling failed: %s", err)
		return
	}

	fmt.Fprintf(out, "%s\n", data)

}
