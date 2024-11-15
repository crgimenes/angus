package angus

import (
	"fmt"
	"math/rand"
	"testing"

	"angus/constants"
)

func TestEncodeDecode(t *testing.T) {
	in := []byte("hello")
	out := make([]byte, MaxPackageSize)

	nout, err := Encode(out, in, 0x01)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, constants.BufferSize)
	cmd, n, err := Decode(data, out[:nout])
	if err != nil {
		t.Fatal(err)
	}

	t.Log("data:", string(data[:n]))

	if cmd != 0x01 {
		t.Errorf("cmd = %v, want 0x01", cmd)
	}

	if string(data[:n]) != string(in) {
		t.Errorf("data = %q, want %q", string(data[:n]), string(in))
	}

	// test invalid size
	_, _, err = Decode(data, out[:nout-1])
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, constants.BufferSize+1), 0x01)
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, 0), 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, constants.BufferSize), 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, constants.BufferSize-1), 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, 1), 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, 0), 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// test invalid size
	_, err = Encode(out, make([]byte, MaxPackageSize+1), 0x01)
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}

	// decode invalid size
	_, _, err = Decode(data, make([]byte, 1))
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}

	// decode invalid size
	_, _, err = Decode(data, make([]byte, 0))
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}

	data = []byte(randonPayload(100))
	_, err = Encode(out, data, 0x01)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	// change size in data to FFFFFFFF
	out[2] = 0xFF
	out[3] = 0xFF
	out[4] = 0xFF
	out[5] = 0xFF
	_, _, err = Decode(data, out)
	if err != ErrInvalidSize {
		t.Errorf("err = %v, want %v", err, ErrInvalidSize)
	}
}

func randonPayload(n int) string {
	ascii := make([]rune, 256)
	for i := range ascii {
		ascii[i] = rune(i)
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = ascii[rand.Intn(len(ascii))]
	}
	return string(b)
}

func TestEncodeDecodeLoop(t *testing.T) {
	data := make([]byte, constants.BufferSize)
	for i := 0; i < 255; i++ {
		in := []byte(randonPayload(rand.Intn(10 + i)))
		out := make([]byte, MaxPackageSize)

		n, err := Encode(out, in, uint8(i))
		if err != nil {
			t.Fatal(err)
		}

		cmd, n, err := Decode(data, out[:n])
		if err != nil {
			t.Fatal(err)
		}

		if cmd != uint8(i) {
			t.Errorf("cmd = %v, want %v", cmd, i)
		}

		if string(data[:n]) != string(in) {
			t.Errorf("data = %v, want %v", string(data), string(in))
		}

		// test size
		if len(in) != len(data[:n]) {
			t.Errorf("size = %v, want %v", len(data), len(in))
		}
	}
}

// Testable example

func ExampleEncode() {
	data := []byte("hello")
	out := make([]byte, MaxPackageSize)

	n, err := Encode(out, data, 0x01)
	if err != nil {
		panic(err)
	}

	fmt.Printf("data: %02X\n", string(out[:n]))
	// Output:
	// data: 010000000568656C6C6F
}

func ExampleDecode() {
	data := []byte("hello")
	out := make([]byte, MaxPackageSize)

	n, err := Encode(out, data, 0x01)
	if err != nil {
		panic(err)
	}

	cmd, n, err := Decode(data, out[:n])
	if err != nil {
		panic(err)
	}

	fmt.Printf("cmd: %02X\n", cmd)
	fmt.Printf("data: %v\n", string(data[:n]))
	// Output:
	// cmd: 01
	// data: hello
}
