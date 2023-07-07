package rpc

import (
	"context"
	"distributedStorage/service/upload/config"
	"distributedStorage/service/upload/proto"
)

type Upload struct {
}

func (u *Upload) UploadEntry(ctx context.Context, req *go_micro_service_upload.ReqEntry, res *go_micro_service_upload.RespEntry) error {
	res.Entry = config.UploadEntry
	return nil
}
