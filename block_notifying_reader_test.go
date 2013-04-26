package eyego

import (
	"testing"
	"bytes"
	"io"
)

func Test8x2(t *testing.T){
	testNotifications([]byte{
		1,2,3,4,5,6,7,8}, 2, []int{1,2,3,5,8,9}, [][]byte{
			[]byte{1,2},
			[]byte{3,4},
			[]byte{5,6},
			[]byte{7,8}}, t)
}

func testNotifications(in []byte, blockSize int, readSizes []int, expected [][]byte, t *testing.T) {
	for i := range readSizes {
		_testNotifications(in, blockSize, readSizes[i], expected, t)
	}
}


func _testNotifications(in []byte, blockSize int, readSize int, expected [][]byte, t *testing.T) {

	Debug("Testing %v size data", blockSize)

	blocks := make([][]byte, 0, 8)

	underlying := bytes.NewBuffer(in)
	bnr := NewBlockNotifyingReader(underlying, blockSize, func(block []byte){
			t.Logf("Appending block %v", block)
			bufCopy := make([]byte, blockSize)
			copy(bufCopy, block)
			blocks = append(blocks, bufCopy)})

	discarder := make([]byte, readSize)

	for {
		_, err := bnr.Read(discarder)

		if err != nil && err != io.EOF {
			t.Errorf("%s", err)
		} else if err == io.EOF {
			break
		}
	}


	if len(expected) != len(blocks) {
		t.Errorf("Wrong number of blocks, expected %d, was %d", len(expected), len(blocks))
	}

	for i := range expected {
		if !bytes.Equal(expected[i], blocks[i]) {
			t.Errorf("Incorrect block number %d. Expected %v, was %v", i, expected[i], blocks[i])
		}
	}
}
