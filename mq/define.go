package mq

type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType int
}
