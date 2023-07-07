package rpc

import (
	"context"
	"distributedStorage/service/download/config"
	download_proto "distributedStorage/service/download/proto"
)

type Downloader struct {
}

func (d *Downloader) DownloadEntry(ctx context.Context, req *download_proto.ReqEntry, resp *download_proto.RespEntry) error {

	resp.Entry = config.DownloadEntry
	return nil
}
