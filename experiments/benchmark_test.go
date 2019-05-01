package experiments

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/cezarsa/fastgelf"
)

func createListener(t testing.TB) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return udpConn
}

func Benchmark_FastGELF_WriteMessageExtra(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		runFastGELFWriteMessage(b, &fastgelf.Message{
			Version:  "1.1",
			TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
			Host:     "myhost",
			Short:    "123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456",
		})
	})
	b.Run("with extra", func(b *testing.B) {
		runFastGELFWriteMessage(b, &fastgelf.Message{
			Version:  "1.1",
			TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
			Host:     "myhost",
			Short:    "123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456",
			Extra: map[string]interface{}{
				"abcdef": "ghijkl",
				"foo":    "bar",
				"xyz":    "baz",
			},
			RawExtra: []byte(`{"ABCD": "EFGH"}`),
		})
	})
}

func Benchmark_gogelf_WriteMessageExtra(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		runGogelfWriteMessage(b, &gelf.Message{
			Version:  "1.1",
			TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
			Host:     "myhost",
			Short:    "123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456",
		})
	})
	b.Run("with extra", func(b *testing.B) {
		runGogelfWriteMessage(b, &gelf.Message{
			Version:  "1.1",
			TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
			Host:     "myhost",
			Short:    "123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456 123456",
			Extra: map[string]interface{}{
				"abcdef": "ghijkl",
				"foo":    "bar",
				"xyz":    "baz",
			},
			RawExtra: []byte(`{"ABCD": "EFGH"}`),
		})
	})
}

func runFastGELFWriteMessage(b *testing.B, m *fastgelf.Message) {
	oldLogger := fastgelf.Logger
	logBuf := bytes.NewBuffer(nil)
	fastgelf.Logger = log.New(logBuf, "", 0)
	defer func() { fastgelf.Logger = oldLogger }()
	listener := createListener(b)
	defer listener.Close()
	w, err := fastgelf.NewUDPWriter(listener.LocalAddr().String())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = w.WriteMessage(m)
		if err != nil {
			b.Fatal(err)
		}
	}
	err = w.Close()
	if err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
	if logBuf.String() != "" {
		b.Fatalf("expected empty error logs, got: %v", logBuf.String())
	}
}

func runGogelfWriteMessage(b *testing.B, m *gelf.Message) {
	listener := createListener(b)
	defer listener.Close()
	w, err := gelf.NewUDPWriter(listener.LocalAddr().String())
	if err != nil {
		b.Fatal(err)
	}
	w.CompressionType = gelf.CompressNone
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = w.WriteMessage(m)
		if err != nil {
			b.Fatal(err)
		}
	}
	err = w.Close()
	if err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
}
