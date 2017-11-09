package cmd

import (
	"bytes"
	"fmt"
	"github.com/wenerme/scel/genproto/v1/sceldata"
	"github.com/wenerme/scel/util"
)

var data *sceldata.ScelData

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
	if data.Pinyins == nil {
		fmt.Println("No pinyin table")
		return
	}

	if scelutil.IsPinyinTableCommon(data.Pinyins) {
		data.Pinyins = nil
	}
}
