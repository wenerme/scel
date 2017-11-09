package scelutil

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wenerme/letsgo/compress"
	"github.com/wenerme/scel/db"
	"github.com/wenerme/scel/fmt"
	"github.com/wenerme/scel/genproto/v1/sceldata"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(file string) (data *sceldata.ScelData, err error) {
	return ReadWithFormat(file, "")
}
func ReadWithFormat(file string, format string) (data *sceldata.ScelData, err error) {
	if format == "" {
		switch filepath.Ext(wcompress.FinalName(file)) {
		case ".json":
			format = "json"
		case ".pb":
			format = "pb"
		case ".scel":
			format = "scel"
		default:
			err = errors.New("can not detect format from filename")
			return
		}
	}
	var b []byte
	if _, b, err = wcompress.DecompressAll(file); err != nil {
		return
	}
	data = &sceldata.ScelData{}
	switch format {
	case "pb":
		err = proto.Unmarshal(b, data)
	case "scel":
		err = scelfmt.Unmarshal(b, data)
	case "json":
		err = json.Unmarshal(b, data)
	default:
		err = errors.New("invalid format")
	}
	return
}

func Write(data *sceldata.ScelData, file string) (err error) {
	return WriteWithFormat(data, file, "")
}
func WriteWithFormat(data *sceldata.ScelData, file string, format string) (err error) {
	if format == "" {
		switch filepath.Ext(file) {
		case ".pb":
			format = "protobuf"
		case ".db":
			format = "sqlite"
		case ".csv":
			format = "csv"
		case ".json":
			format = "json"
		default:
			err = errors.New("can not detect format from filename")
			return
		}
	}
	var b []byte
	switch format {
	case "csv":
	case "sqlite":
		return writeDb(data, file, "sqlite3")
	case "json":
		b, err = json.Marshal(data)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(file, b, os.ModePerm)
	case "pb":
		fallthrough
	case "protobuf":
		b, err = proto.Marshal(data)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(file, b, os.ModePerm)
	default:
		err = errors.New("invalid format")
	}
	return
}
func writeDb(data *sceldata.ScelData, arg string, dialet string) (err error) {
	db, err := gorm.Open(dialet, arg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.AutoMigrate(sceldb.Scel{}, sceldb.Word{}).Error; err != nil {
		return err
	}
	{
		db := db.Begin()
		defer func() {
			if err != nil {
				db.Rollback()
			} else {
				db.Commit()
			}
		}()
		scel := &sceldb.Scel{
			Name:        data.Info.Name,
			Type:        data.Info.Type,
			Description: data.Info.Description,
			Example:     data.Info.Example,
		}
		if err = db.Create(scel).Error; err != nil {
			return
		}

		for _, v := range data.Words {
			var pinyins []string
			for _, p := range v.Pinyins {
				pinyins = append(pinyins, data.Pinyins[p])
			}
			b, _ := json.Marshal(pinyins)
			py := string(b)

			for _, w := range v.Words {
				r := &sceldb.Word{
					Scel:   scel,
					Pinyin: py,
					Word:   w,
				}
				if err = db.Create(r).Error; err != nil {
					err = errors.Wrap(err, "failed to insert word")
					return
				}
			}

		}
	}

	return
}
