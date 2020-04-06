package profile

import (
	"bufio"
	"io"

	"github.com/cfi2017/bl3-save/internal/pb"
	"github.com/cfi2017/bl3-save/internal/shared"
	"github.com/golang/protobuf/proto"
)

func Deserialize(reader io.Reader) (shared.SavFile, pb.Profile) {

	// deserialise header, decrypt data
	s, data := shared.DeserializeHeader(reader)

	data = shared.Decrypt(data)

	p := pb.Profile{}
	if err := proto.Unmarshal(data, &p); err != nil {
		panic("couldn't unmarshal protobuf data")
	}

	return s, p
}

func Serialize(writer io.Writer, s shared.SavFile, p pb.Profile) {
	w := bufio.NewWriter(writer)
	w.WriteString("GVAS")
}
