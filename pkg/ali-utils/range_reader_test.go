package ali_utils

import (
	"errors"
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"io/ioutil"
	"testing"
)

func TestRangeReader(t *testing.T) {
	s := suppress.PanicAll{}
	var xb = make([]byte, 257*1024)
	for i := range xb {
		xb[i] = byte(i & 0xff)
	}

	var gg = NewBytesRandReader(xb)

	var gg2 = NewRangeReader(&dispatcher.RandReaderUploadRequest{
		BaseUploadRequest: dispatcher.BaseUploadRequest{
			XFileName:  "",
			XSize:      257 * 1024,
			XChunkHint: 128 * 1024,
		},
		R: gg,
		S: s,
	}, 128*1024, 257*1024)

	bb, err := ioutil.ReadAll(gg2)
	s.Suppress(err)
	if len(bb) != len(xb) {
		s.Suppress(errors.New("gg len"))
	}
	for i := range xb {
		if bb[i] != xb[i] {
			s.Suppress(errors.New("gg content"))
		}
	}
}
