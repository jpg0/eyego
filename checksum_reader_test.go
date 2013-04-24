package eyego

import (
	"testing"
	"encoding/base64"
	"strings"
	"io/ioutil"
)

type pair struct {
	data     string
	checksum string
}

type tcpPair struct {
	data     string
	checksum uint16
}

var (
	key            = "1234abcd1234abcd"
	dataToChecksum = map[int] pair {
	0: pair{
		data:"",
		checksum:"b9e443175e20fe0c9fa97b317264278d"},
	512: pair{
		data:"vSDhVo42LZJo4JSgn9Hl9q7MRh/ORO6CrS/vR2zGw7cVoWHextPcPot4ugqgnBBK51XmixfOzcrxFve3VZHY9UKPVH9kd+Op9MI/N0QsFKxT8NkInS4ltKal7bVaPQ7DnO5GnWs/hMgskB58thiCHjglwPzDycXbWGFj4QdUaa1Prvwt5FefNdvX+WF9yD+GRw4A+Kp27sXVwIOPUXvyRDmQF0vlu7ogVWhS4YceKaYvKLn1x2SM+5R5pShOAAbSZgyYIqpAzVsqX21yb+2vNBw4jy7cBXsh+jU2EsChgiXDukaUp0XERR1M9XyU/CEVKgPrO18cn2SXFl+Za1t2SShyrCXe2YwXkYRhGNehi+K8ysqdr+SA/Y64lqeb4mXBAzjau+wH61x7DqZcTCJ49QQZ6T4mxd/MMCixD1FOQJjQZ14Len0Gh0U9gJmbEfA65Z2qvMALc2Ru5K52L5cLbFWVcVWMoHP4yLWbB1GAHoSGepTBM6hoBYVkZlP8LB9Omj+u3nPwt+ub6I/ZXgA7Cu0rAnY6mQWxKPEtRpbb2+89TTDXVLfq7bAqczE2WT1shpQSzLpl0JaoCfR+gl42eNElS3TZx4mROR56i/6xogqZVlfB/C0BfrqQRHqPEZO8kwgZuXLhEpiBU8EgbNZPMwnC0fOHCg51j/EVUcWCIcM=",
		checksum:"3ba25d0207240171fb70d6465ea846d5"}}

	tcpChecksums = map[int] tcpPair {
	1: tcpPair{
		data:"vSDhVo42LZJo4JSgn9Hl9q7MRh/ORO6CrS/vR2zGw7cVoWHextPcPot4ugqgnBBK51XmixfOzcrxFve3VZHY9UKPVH9kd+Op9MI/N0QsFKxT8NkInS4ltKal7bVaPQ7DnO5GnWs/hMgskB58thiCHjglwPzDycXbWGFj4QdUaa1Prvwt5FefNdvX+WF9yD+GRw4A+Kp27sXVwIOPUXvyRDmQF0vlu7ogVWhS4YceKaYvKLn1x2SM+5R5pShOAAbSZgyYIqpAzVsqX21yb+2vNBw4jy7cBXsh+jU2EsChgiXDukaUp0XERR1M9XyU/CEVKgPrO18cn2SXFl+Za1t2SShyrCXe2YwXkYRhGNehi+K8ysqdr+SA/Y64lqeb4mXBAzjau+wH61x7DqZcTCJ49QQZ6T4mxd/MMCixD1FOQJjQZ14Len0Gh0U9gJmbEfA65Z2qvMALc2Ru5K52L5cLbFWVcVWMoHP4yLWbB1GAHoSGepTBM6hoBYVkZlP8LB9Omj+u3nPwt+ub6I/ZXgA7Cu0rAnY6mQWxKPEtRpbb2+89TTDXVLfq7bAqczE2WT1shpQSzLpl0JaoCfR+gl42eNElS3TZx4mROR56i/6xogqZVlfB/C0BfrqQRHqPEZO8kwgZuXLhEpiBU8EgbNZPMwnC0fOHCg51j/EVUcWCIcM=",
		checksum: 60593}}
)

func TestEmpty(t *testing.T) {
	testChecksum(0, t)
}

func TestSingleBlock(t *testing.T) {
	testChecksum(512, t)
}

func testChecksum(size int, t *testing.T) {
	cr := NewChecksumReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(dataToChecksum[size].data)))

	b, err := ioutil.ReadAll(cr)

	if err != nil {
		t.Error("%s", err)
	}

	if len(b) != size {
		t.Error("wrong size, expected %d, was %d", size, len(b))
	}

	cs := cr.Checksum(key)

	if cs != dataToChecksum[size].checksum {
		t.Fatalf("checksum failed, expected %s, was %s", dataToChecksum[size].checksum, cs)
	}
}

func TestTcpChecksum(t *testing.T) {
	data, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(tcpChecksums[1].data)))

	if err != nil {
		t.Error("%s", err)
	}

	tcs := tcp_checksum(data)
	if tcs != tcpChecksums[1].checksum {
		t.Fatalf("checksum failed, expected %s, was %s", tcpChecksums[1].checksum, tcs)

	}
}
