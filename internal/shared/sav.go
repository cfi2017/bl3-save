package shared

type SavFile struct {
	SgVersion          int
	PkgVersion         int
	EngineMajorVersion int
	EngineMinorVersion int
	EnginePatchVersion int
	EngineBuildVersion int
	BuildId            string
	FmtVersion         int
	FmtCount           int
	CustomFmtData      []CustomFormatData
	SgType             string
}

type CustomFormatData struct {
	Guid  string
	Entry int
}
