package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/wenerme/scel/genproto/v1/sceldata"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const (
	OFFSET_PINGYIN = 0x1540
	OFFSET_CHINESE = 0x2628
)

var (
	MAGIC    = []byte{0x40, 0x15, 0x00, 0x00, 0x44, 0x43, 0x53, 0x01, 0x01, 0x00, 0x00, 0x00}
	PY_MAGIC = []byte{0x9D, 0x01, 0x00, 0x00}
)

type Parser struct {
	Bytes   []byte
	Data    *sceldata.ScelData
	decoder *encoding.Decoder
	pyRemap map[uint16]int32
}

func NewParser() *Parser {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	decoder := utf16.NewDecoder()
	return &Parser{
		decoder: decoder,
	}
}
func (self *Parser) Reset(b []byte) {
	self.Bytes = b
	self.Data = &sceldata.ScelData{}
	self.decoder.Reset()
}
func (self *Parser) ReadInfo() *sceldata.Info {
	self.Data.Info = &sceldata.Info{
		Name:        self.str(0x130, 0x338),
		Type:        self.str(0x338, 0x540),
		Description: self.str(0x540, 0xD40),
		Example:     self.str(0xD40, OFFSET_PINGYIN),
	}
	return self.Data.Info
}
func (self *Parser) IsMagicMatch() bool {
	if len(self.Bytes) <= OFFSET_PINGYIN+len(PY_MAGIC) {
		return false
	}
	return bytes.Equal(self.Bytes[0:len(MAGIC)], MAGIC) && bytes.Equal(self.Bytes[OFFSET_PINGYIN:OFFSET_PINGYIN+len(PY_MAGIC)], PY_MAGIC)
}

func (self *Parser) str(a, b int) string {
	return self.readString(self.Bytes[a:b])
}

func (self *Parser) readString(b []byte) string {
	i := 0
	for ; i < len(b); i += 2 {
		if b[i] == 0 && b[i+1] == 0 {
			break
		}
	}
	if i > 0 {
		dst, err := self.decoder.Bytes(b[:i])
		defer self.decoder.Reset()
		if err != nil {
			panic(err)
		}
		return string(dst)
	}
	return ""
}

func (self *Parser) ReadPinyin() {
	b := self.Bytes[OFFSET_PINGYIN:OFFSET_CHINESE]
	b = b[len(PY_MAGIC):]
	pyRemap := make(map[uint16]int32)
	pys := make([]string, 0)
	for len(b) > 0 {
		idx := binary.LittleEndian.Uint16(b)
		l := binary.LittleEndian.Uint16(b[2:])
		b = b[4:]
		s := self.readString(b[:l])

		b = b[l:]

		pyRemap[idx] = int32(len(pyRemap))
		pys = append(pys, s)
	}

	self.pyRemap = pyRemap
	self.Data.Pinyins = pys
}

func (self *Parser) ReadWord() {
	b := self.Bytes[OFFSET_CHINESE:]
	for len(b) > 0 {
		w := &sceldata.Word{}
		// 同音词
		same := int(binary.LittleEndian.Uint16(b))
		b = b[2:]

		// 拼音
		pyLen := int(binary.LittleEndian.Uint16(b))
		b = b[2:]
		// 2 per py, pyLen/2
		for i := 0; i < pyLen/2; i++ {
			w.Pinyins = append(w.Pinyins, self.pyRemap[binary.LittleEndian.Uint16(b[i*2:])])
		}
		b = b[pyLen:]

		for i := 0; i < same; i++ {
			// 词组
			wordLen := int(binary.LittleEndian.Uint16(b))
			b = b[2:]
			word := self.readString(b[:wordLen])

			b = b[wordLen:]

			// 扩展
			extLen := int(binary.LittleEndian.Uint16(b))
			b = b[2:]
			ext := b[:extLen]

			b = b[extLen:]

			w.Words = append(w.Words, word)
			w.Exts = append(w.Exts, ext)
		}

		self.Data.Words = append(self.Data.Words, w)
	}
}

func (self *Parser) ReadData() (*sceldata.ScelData, error) {
	if self.IsMagicMatch() {
		return nil, errors.New("Invalid data")
	}
	self.ReadInfo()
	self.ReadPinyin()
	self.ReadWord()
	return self.Data, nil
}
