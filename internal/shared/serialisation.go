package shared

import (
	"bufio"
	"io"
)

func DeserializeHeader(reader io.Reader) (SavFile, []byte) {
	r := bufio.NewReader(reader)
	s := SavFile{}

	// ensure file starts with GVAS
	if header := ReadNBytes(r, 4); string(header) != "GVAS" {
		panic("error reading header")
	}
	s.SgVersion = ReadInt(r)
	s.PkgVersion = ReadInt(r)
	s.EngineMajorVersion = ReadShort(r)
	s.EngineMinorVersion = ReadShort(r)
	s.EnginePatchVersion = ReadShort(r)
	s.EngineBuildVersion = ReadInt(r)
	s.BuildId = ReadString(r)
	s.FmtVersion = ReadInt(r)
	s.FmtCount = ReadInt(r)
	s.CustomFmtData = make([]CustomFormatData, s.FmtCount)
	for i := 0; i < s.FmtCount; i++ {
		data := CustomFormatData{}
		data.Guid = ReadGuid(r)
		data.Entry = ReadInt(r)
		s.CustomFmtData[i] = data
	}
	s.SgType = ReadString(r)

	l := ReadInt(r)
	data := ReadNBytes(r, l)
	if _, err := r.ReadByte(); err.Error() != "EOF" {
		panic("didn't get eof though expecting eof")
	}

	return s, data
}

func SerializeHeader(writer io.Writer, s SavFile, content []byte) {

	w := bufio.NewWriter(writer)
	WriteBytes(w, []byte("GVAS"))
	WriteInt(w, s.SgVersion)
	WriteInt(w, s.PkgVersion)
	WriteShort(w, s.EngineMajorVersion)
	WriteShort(w, s.EngineMinorVersion)
	WriteShort(w, s.EnginePatchVersion)
	WriteInt(w, s.EngineBuildVersion)
	WriteString(w, s.BuildId)
	WriteInt(w, s.FmtVersion)
	WriteInt(w, len(s.CustomFmtData))
	for _, d := range s.CustomFmtData {
		WriteGuid(w, d.Guid)
		WriteInt(w, d.Entry)
	}
	WriteString(w, s.SgType)

	WriteInt(w, len(content))
	WriteBytes(w, content)
}
