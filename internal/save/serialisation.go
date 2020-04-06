package save

import (
	"io"

	"github.com/cfi2017/bl3-save/internal/pb"
	"github.com/cfi2017/bl3-save/internal/shared"
	"google.golang.org/protobuf/proto"
)

func Deserialize(reader io.Reader) (shared.SavFile, pb.Character) {
	// deserialise header, decrypt data
	s, data := shared.DeserializeHeader(reader)

	data = shared.Decrypt(data)
	p := pb.Character{}
	if err := proto.Unmarshal(data, &p); err != nil {
		panic("couldn't unmarshal protobuf data")
	}

	return s, p
}
