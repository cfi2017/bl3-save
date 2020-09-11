package server

import (
	"bytes"
	"errors"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/cfi2017/bl3-save-core/pkg/shared"
	"io"
	"io/ioutil"
)

var ErrInvalidSavFile = errors.New("invalid save file")

type Platforms map[string]shared.Magic

type DeserializeFunc func(reader io.Reader, magic shared.Magic) (shared.SavFile, pb.Character, error)

func TryDeserialize(deserializeFunc DeserializeFunc, platforms Platforms, reader io.Reader) (s shared.SavFile, char pb.Character, platform string, err error) {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return shared.SavFile{}, pb.Character{}, "", err
	}
	for p, magic := range platforms {
		reader = bytes.NewReader(bs)
		s, char, err = deserializeFunc(reader, magic)
		platform = p
		if err == nil {
			return
		}
	}
	return shared.SavFile{}, pb.Character{}, "", ErrInvalidSavFile
}
