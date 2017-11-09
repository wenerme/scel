package scelfmt

import "github.com/wenerme/scel/genproto/v1/sceldata"

func Unmarshal(b []byte, data *sceldata.ScelData) (err error) {
	parser := NewParser()
	parser.Reset(b)
	parser.Data = data
	_, err = parser.ReadData()
	return
}
