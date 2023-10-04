package itswizard_m_s3bucket

// WriteBuffer is a simple type that implements io.WriterAt on an in-memory buffer.
// The zero value of this type is an empty buffer ready to use.
type WriteBuffer struct {
	d []byte
	m int
}

// NewWriteBuffer creates and returns a new WriteBuffer with the given initial size and
// maximum. If maximum is <= 0 it is unlimited.
func newWriteBuffer(size, max int) *WriteBuffer {
	if max < size && max >= 0 {
		max = size
	}
	return &WriteBuffer{make([]byte, size), max}
}

// Bytes returns the WriteBuffer's underlying data. This value will remain valid so long
// as no other methods are called on the WriteBuffer.
func (wb *WriteBuffer) Bytes() []byte {
	return wb.d
}

func (wb *WriteBuffer) WriteAt(dat []byte, off int64) (int, error) {
	wb.d = dat
	wb.m = int(off)
	return len(dat), nil
}
