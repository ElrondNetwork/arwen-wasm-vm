package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const suffix = ".toml"

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fs, _ := ioutil.ReadDir(".")
	out, _ := os.Create("tomlfiles.go")
	out.Write([]byte("package gasschedules \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), suffix) {
			out.Write([]byte(strings.TrimSuffix(f.Name(), suffix) + " = `"))
			f, _ := os.Open(f.Name())
			io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
