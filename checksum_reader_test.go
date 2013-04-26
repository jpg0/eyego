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
		checksum:"3ba25d0207240171fb70d6465ea846d5"},
	1024: pair{
		data:"O9n7d2lsXTUaWxlfSaqu65nynGPsCOxLheezftx68/q1o59R+JTxl3Jj/MNTNPCsYM2XQAkEppdOLv0EaPpXNabgiwM70tIaeh2IuK0qV7XgrHFrx/R9755kVAnZq9gLWfg0C2fByVQx9LaYyUnZLuh75AYBeO63sUgDjUOWENr55GH546aeNrxqoeSnrQ7uHHCdln3hSrMaZ88SMbpHiecE9i60BeHPdIMF0U35iTGCDbo3Z4GnzHJaWt+/eQCjoh29SYi+2GWq8oMvgvSCUpBB02BmW9J9A4MKgMTyvbi0p+NZC6/qRm+xlMtAVsN1fK05SSm5w/BKB7TETawCX3Iplt37PkxxKaDBwgl+plTb0O3z3PExpqfnR+9LfWFXpWW+UucpmNckVJkVpfRSbYIPvJKhNl2Iy0YoUyXuqECywy6VuwsF85cVe0C0QBfj9GpHsIMNQAh8yUQaOJ7NXRlIKq7Zc4UD8d1GI/+WYSVoSzv5d9DSuOnoqZBsFtqPT8PLKIB/ztEbNmHFMLN9n1vLgWnRbmuCCXsy9Jrqpq5wSHiWSltax6Qyfz0e+wMltumIIeP47jJWox66LtSYHSu8geN0cfbNPJfsaPuFpBrfN2N5zTKE71cXFwESCjFhHMuku0LxHSNvYQCg6HtyfHKTnKpL4hH6PIwyVqoYgoufiYRiPvLkzLPymG9mPfoHJY7avsURoQ1ocJQDtUjH//AL2adds5Sk9urKxq5oen/LqaoEjQso0W25ZF3RSP/pfr1nygmQZRYTdk5J+YWrjS0yFpM73Bqu4fFBd0j+sg4/bInBOHmduHw1QoF2G/gegcbvJSsbJcks5j5vevQ/Ns3FbYqyuQ9BjeYkVWvMRi50AdyRpXDhfoB764Jl8aC/asldiidMHdZHFerHZN2vHkvGHy0uCFQBmu0O+jW63RZpZ98aCTqoGVAFvZdfAU/WdSTAI8020VvUssXZR8QuC7PQGoCn9efGWKCExOQeBaYKxnIVOMFLjWSwAh1HQQhiW16IrhKMV84pmZ3b8SB5/Khz9Q8fQhK6LdPlop+uwa6oEwYmIVyf8XHlDhC93Ow+flrcO+z+NO9mHnIwreE9XU3HDUgF2rALP3qmiDHElFcXq8+r0aGElvlfMUN1YkEVe+6N5kblchOrymrFp90Sph8yCi/hCFU8DXBmLVZ00MNXUWZroyeiMW9V3UWaj10c6xrgjnbjvHaelMqnQc7KsFqmcDuzi1ErthJa7lTGQ4FbW6QuS0JymOHK1/oXCEUvIdEiT7ELEx8CpL3dUdyxXJ0YqZ3XbyxfxAtEBKZhZDP6IZVcZVV3KU66HPwg7Yipkgo8TnFh0dvaqx54T8NUdw==",
		checksum:"6a86e824819019abae3fccfa14a14c5b"}}

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

func Test2Blocks(t *testing.T) {
	testChecksum(1024, t)
}

func testChecksum(size int, t *testing.T) {

	Debug("Testing %v size data for checksum", size)

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

	Debug("Completed testing %v size data for checksum", size)
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
