package tfsum

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"time"

	"path"

	"github.com/golang/protobuf/proto"
	summary "github.com/helinwang/tfsum/tensorflow/core/framework"
	event "github.com/helinwang/tfsum/tensorflow/core/util"
)

// Writer writes string of serialized tensorflow.Summary into file
// can be used for tensorboard
type Writer struct {
	Dir  string
	Name string
	f    *os.File
}

// Close closes file created by writer
func (w *Writer) Close() {
	if w.f != nil {
		w.f.Close()
	}
}

// AddEvent adds a tensorflow.Summary event
func (w *Writer) AddEvent(s string, step int64) error {
	var spro summary.Summary
	err := proto.Unmarshal([]byte(s), &spro)
	if err != nil {
		return err
	}
	start := time.Now().Unix()
	e := event.Event{
		WallTime: float64(start),
		Step:     step,
		What:     &event.Event_Summary{Summary: &spro},
	}
	b, err := proto.Marshal(&e)
	if err != nil {
		return err
	}
	if w.f == nil {
		dir := path.Join(w.Dir, w.Name)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		name := path.Join(w.Dir, w.Name, "events.out.tfevents."+strconv.Itoa(int(start))+".user")
		w.f, err = os.Create(name)
		if err != nil {
			return err
		}
	}

	// reference: https://github.com/tensorflow/tensorflow/blob/5a5a25ea3ebef623e07fb9a46419a9df377a37a5/tensorflow/core/lib/io/record_writer.cc#L98
	// Format of a single record:
	//  uint64    length
	//  uint32    masked crc of length
	//  byte      data[length]
	//  uint32    masked crc of data
	l := make([]byte, 8)
	binary.LittleEndian.PutUint64(l, uint64(len(b)))
	if err != nil {
		return err
	}

	table := crc32.MakeTable(crc32.Castagnoli)
	ck := mask(crc32.Checksum(l, table))
	c := make([]byte, 4)
	binary.LittleEndian.PutUint32(c, ck)
	ck = mask(crc32.Checksum(b, table))
	cc := make([]byte, 4)
	binary.LittleEndian.PutUint32(cc, ck)

	_, err = w.f.Write(l)
	if err != nil {
		return err
	}
	_, err = w.f.Write(c)
	if err != nil {
		return err
	}
	_, err = w.f.Write(b)
	if err != nil {
		return err
	}
	_, err = w.f.Write(cc)
	if err != nil {
		return err
	}
	return nil
}

func mask(crc uint32) uint32 {
	// reference: https://github.com/tensorflow/tensorflow/blob/754048a0453a04a761e112ae5d99c149eb9910dd/tensorflow/core/lib/hash/crc32c.h#L40
	// Rotate right by 15 bits and add a constant.
	return ((crc >> 15) | (crc << 17)) + 0xa282ead8
}
