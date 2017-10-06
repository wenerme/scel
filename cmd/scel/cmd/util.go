package cmd

import (
	"github.com/wenerme/scel/genproto/v1/sceldata"
	"os"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"github.com/wenerme/scel/parser"
	"path/filepath"
	"crypto/sha256"
	"encoding/base64"
	"bytes"
)

var data *sceldata.ScelData

const SHA256_COMMON_PY = "YB4qyXsMZRJ_6krGfnJ1DEosdYYocB5BJt02YWU2F-o"

func open(fn string) {
	switch filepath.Ext(fn) {
	case ".pb":
		openPb(fn)

	default:
		fallthrough
	case ".scel":
		openScel(fn)
	}
}

func write(fn string) {
	switch filepath.Ext(fn) {
	default:
		fallthrough
	case ".pb":
		writePb(fn)
	}
}

func openScel(fn string) {
	p := parser.NewParser()
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	p.Reset(b)
	data, err = p.ReadData()
	if err != nil {
		panic(err)
	}
}

func openPb(fn string) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	data = &sceldata.ScelData{}
	err = proto.Unmarshal(b, data)
	if err != nil {
		panic(err)
	}
}

func writePb(fn string) {
	b, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fn, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func doExcludeExt() {
	for _, v := range data.Words {
		v.Exts = nil
	}
}

func doOptimizeExt() {
	for _, w := range data.Words {
		for k, v := range w.Exts {
			w.Exts[k] = bytes.TrimRight(v, "\x00")
		}
	}
}

func doExcludeCommonPy() {
	hash := sha256.New()
	for _, v := range data.Pinyins {
		hash.Write([]byte(v))
	}
	sum := hash.Sum(nil)

	if SHA256_COMMON_PY == base64.RawURLEncoding.EncodeToString(sum) {
		data.Pinyins = nil
	}
}
