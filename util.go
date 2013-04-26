package eyego

import (
	"io"
	"bufio"
	"strconv"
	"os"
)

func EachLine(r io.Reader, f func(string) interface {}) (rv []interface {}, err error){
	br := bufio.NewReader(r)

	rv = make([]interface {}, 10)

	for ;err != io.EOF; {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}
		rv = append(rv, f(line[:len(line) - 1]))
	}

	return rv, err
}

func Atoi(s string) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func CopyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}
