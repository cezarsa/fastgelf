package fastgelf

import (
	"encoding/json"
	"net"
	"reflect"
	"testing"
	"time"
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

func Test_UDPWriter_WriteMessage(t *testing.T) {
	listener := createListener(t)
	defer listener.Close()

	msg := &Message{
		Version:  "1.1",
		TimeUnix: float64(time.Unix(1000, 100000000).UnixNano()) / float64(time.Second),
		Host:     "myhost",
		Short:    "123456",
		Facility: "kernel",
		Full:     "full msg",
		Level:    3,
		Extra: map[string]interface{}{
			"abcdef": "ghijkl",
			"foo":    "bar",
			"xyz":    "baz",
		},
		RawExtra: []byte(`{"ABCD": "EFGH"}`),
	}

	w, err := NewUDPWriter(listener.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	err = w.WriteMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	listener.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := listener.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	var gotMsg map[string]interface{}
	err = json.Unmarshal(buffer[:n], &gotMsg)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"version":       "1.1",
		"host":          "myhost",
		"short_message": "123456",
		"timestamp":     1000.1,
		"facility":      "kernel",
		"level":         3.0,
		"full_message":  "full msg",
		"abcdef":        "ghijkl",
		"foo":           "bar",
		"xyz":           "baz",
		"ABCD":          "EFGH",
	}
	if !reflect.DeepEqual(gotMsg, expected) {
		t.Fatalf("unexpected message received.\nexpected: %#v\ngot: %#v", expected, gotMsg)
	}
}
