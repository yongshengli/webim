package models

import (
	"testing"
	"webim/comet/common"
)

func TestInsertDevice(t *testing.T) {
	dId := common.Uuid()
	d := &Device{DeviceId: dId}
	res := InsertDevice(dId, d)

	if res.Error != nil {
		t.Error(res.Error)
	}

}

func init() {
	ConnectTestMysql()
}
