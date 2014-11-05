// Based on http://golang.org/src/pkg/hash/crc32/crc32_test.go
// Adapted for CRC-16 by Nick Patavalis (http://npat.efault.net/)

package crc16

import (
	"io"
	"testing"
)

type crcCheck struct {
	cfg *Conf
	sum uint16
	in  string
}

var cTbl = []crcCheck{
	{X25, 0x906e, "123456789"},
	{PPP, 0x906e, "123456789"},
	{Modbus, 0x4b37, "123456789"},
	{XModem, 0x31c3, "123456789"},
	{Kermit, 0x2189, "123456789"},
}

func TestCRCCheck(t *testing.T) {
	for i := range cTbl {
		sum := Checksum(cTbl[i].cfg, []byte(cTbl[i].in))
		if sum != cTbl[i].sum {
			t.Errorf("C%d: 0x%04x want 0x%04x", i, sum, cTbl[i].sum)
		}
	}
}

type test struct {
	ppp, modbus uint16
	in          string
}

var golden = []test{
	{0x0000, 0xffff, ""},
	{0x82f7, 0xa87e, "a"},
	{0x33de, 0xc9a9, "ab"},
	{0x9e25, 0x5749, "abc"},
	{0xa36b, 0x1d97, "abcd"},
	{0x19a5, 0x859c, "abcde"},
	{0x04f6, 0x4305, "abcdef"},
	{0x757c, 0xe9c2, "abcdefg"},
	{0xa6a8, 0x7f69, "abcdefgh"},
	{0x275b, 0x007f, "abcdefghi"},
	{0xd055, 0xcfc1, "abcdefghij"},
	{0x7025, 0x21a5, "Discard medicine more than two years old."},
	{0x2be0, 0x70e7, "He who has a shady past knows that nice guys finish last."},
	{0x81d3, 0x1974, "I wouldn't marry him with a ten foot pole."},
	{0xf471, 0xc315, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0xdc73, 0x0eab, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x6a62, 0xa782, "Nepal premier won't resign."},
	{0xd860, 0x9201, "For every action there is an equal and opposite government program."},
	{0xee04, 0xefaf, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x1687, 0xac2a, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x497d, 0xad2e, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x0073, 0x9205, "size:  a.out:  bad magic"},
	{0xccb3, 0xd7d3, "The major problem is with sendmail.  -Mark Horton"},
	{0x9464, 0x39a3, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x318e, 0x24fd, "If the enemy is within range, then so are you."},
	{0x2cbb, 0xc7cf, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0xde8b, 0x1f3b, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xfe32, 0x86c7, "C is as portable as Stonehedge!!"},
	{0x9186, 0xb89a, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xd304, 0xee28, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0xfd73, 0xb025, "How can you write a big system without C++?  -Paul Glick"},
}

func TestGolden(t *testing.T) {
	for _, g := range golden {
		h := New(PPP)
		io.WriteString(h, g.in)
		s := h.Sum16()
		if s != g.ppp {
			t.Errorf("PPP(%s) = 0x%x want 0x%x", g.in, s, g.ppp)
		}

		h = New(Modbus)
		io.WriteString(h, g.in)
		s = h.Sum16()
		if s != g.modbus {
			t.Errorf("MODBUS(%s) = 0x%x want 0x%x", g.in, s, g.modbus)
		}
	}
}

func bench(b *testing.B, sz int64) {
	b.SetBytes(sz)
	data := make([]byte, sz)
	for i := range data {
		data[i] = byte(i)
	}
	h := New(Modbus)
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(in)
	}
}

func BenchmarkCrc32B(b *testing.B) {
	bench(b, 32)
}

func BenchmarkCrc128B(b *testing.B) {
	bench(b, 128)
}

func BenchmarkCrc256B(b *testing.B) {
	bench(b, 256)
}

func BenchmarkCrcKB(b *testing.B) {
	bench(b, 1024)
}
