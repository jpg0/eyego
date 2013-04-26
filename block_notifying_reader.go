package eyego

import "io"

type BlockNotifyingReader struct {
	delegate io.Reader
	state *blockState
	blockSize int
	notifier func([]byte)
}

type blockState struct {
	buf []byte
	ptr int
}

func NewBlockNotifyingReader(_delegate io.Reader, _blockSize int, _notifier func([]byte)) BlockNotifyingReader {
	return BlockNotifyingReader{
		delegate: _delegate,
		blockSize: _blockSize,
		notifier: _notifier,
		state: &blockState{
			buf: make([]byte, _blockSize),
			ptr: 0}}
}

func (r BlockNotifyingReader) Read(p []byte) (n int, err error){
	n, err = r.delegate.Read(p)

	if err == nil && n > 0 {
		r.appendBytes(p, n)
	}

	return n, err
}

func (r BlockNotifyingReader) appendBytes(b []byte, len int) {

	Trace("Appending %d bytes", len)

	srcPtr := 0

	for ;r.state.ptr + (len - srcPtr) >= r.blockSize; {
		Trace("Adding block")
		Trace("Filling from source[%d:%d] to buf[%d:%d]", srcPtr, r.blockSize-r.state.ptr + srcPtr, r.state.ptr, r.blockSize)
		toAdd := copy(r.state.buf[r.state.ptr:r.blockSize], b[srcPtr:r.blockSize-r.state.ptr + srcPtr]) //copy bytes to fill temp buffer
//		Trace("Notifying %v", r.state.buf)
		r.notifier(r.state.buf)
//		r.state.buf = r.state.buf[:0]
		r.state.ptr = 0
		srcPtr += toAdd
	}

	Trace("srcPtr: %d", srcPtr)

	//copy remaining
	Trace("Adding %d bytes to (%d byte) buffer", len-srcPtr, r.state.ptr)
	Trace("Copying from source[%d:%d] to buf[%d:%d]", srcPtr, len-srcPtr, r.state.ptr, r.state.ptr+len-srcPtr)
	r.state.ptr += copy(r.state.buf[r.state.ptr:r.state.ptr+len-srcPtr], b[srcPtr:len])
}
