package ali_utils

import "io"

type XReaderAt interface {
	ReadAt(pos, maxLen int64) io.Reader
}

type rangeReaderBuf struct {
	r             XReaderAt
	chunkSize     int64
	size          int64
	pos           int64
	currentReader io.Reader
}

func (r *rangeReaderBuf) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.currentReader == nil {
		if !r.nextReader() {
			return 0, io.EOF
		}
	}

	if n, err = r.currentReader.Read(p); err != nil {
		if err == io.EOF {
			if !r.nextReader() {
				return n, err
			}

			p = p[n:]
			nn, err2 := r.Read(p)
			nn += n
			if nn != 0 && err2 == io.EOF {
				return nn, nil
			}
			return nn, err2
		}

		return
	}

	return
}

func (r *rangeReaderBuf) nextReader() bool {
	if r.pos >= r.size {
		return false
	}

	chunkSize := r.chunkSize
	if chunkSize > r.size-r.pos {
		chunkSize = r.size - r.pos
	}
	r.currentReader = r.r.ReadAt(r.pos, chunkSize)
	if r.currentReader == nil {
		return false
	}
	r.pos += chunkSize
	return true
}

func NewRangeReader(r XReaderAt, chunkSize, sz int64) io.Reader {
	if chunkSize == 0 {
		chunkSize = 128 * 1024
	}
	return &rangeReaderBuf{r: r, chunkSize: chunkSize, size: sz}
}
