// Based on http://golang.org/src/pkg/hash/crc32/crc32.go
// Adapted for CRC-16 by Nick Patavalis (http://npat.efault.net/)

// Package crc16 implements the 16-bit cyclic redundancy check, or
// CRC-16, checksum. See
// http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.
package crc16

import (
	"hash"
	"sync"
)

// The size of a CRC-16 checksum in bytes.
const Size = 2

// Typical CRC-16 configurations. Mostly used are the CCITT (0x1021)
// and the IBM/ANSI (0x8005) polynomials, either bit-reversed or
// not. For more configurations see:
// http://reveng.sourceforge.net/crc-catalogue/
var (
	X25 = &Conf{
		Poly: 0x1021, BitRev: true,
		IniVal: 0xffff, FinVal: 0xffff,
		BigEnd: false,
	}
	PPP    = X25
	Modbus = &Conf{
		Poly: 0x8005, BitRev: true,
		IniVal: 0xffff, FinVal: 0x0,
		BigEnd: false,
	}
	XModem = &Conf{
		Poly: 0x1021, BitRev: false,
		IniVal: 0x0000, FinVal: 0x0,
		BigEnd: true,
	}
	Kermit = &Conf{
		Poly: 0x1021, BitRev: true,
		IniVal: 0x0, FinVal: 0x0,
		BigEnd: false,
	}
)

// Conf is a CRC configuration. It is passed to functions New and
// Checksum and specifies the parameters of the calculated
// checksum. The first time New or Checksum are called with a
// configuration structure c, they calculate the polynomial table for
// this configuration; subsequent calls with the same c use the table
// already calculated. A few commonly used configurations are defined
// as global variables (X25, PPP, Modbus, etc.)
type Conf struct {
	Poly   uint16 // Polynomial to use.
	BitRev bool   // Bit reversed CRC (bit-15 is X^0)?
	IniVal uint16 // Initial value of CRC register.
	FinVal uint16 // XOR CRC with this at the end.
	BigEnd bool   // Emit *bytes* most significant first (see Hash.Sum)?
	once   sync.Once
	table  *Table
	update func(uint16, *Table, []byte) uint16
}

// reverse returns the bit-reversed of v: 0xA001 --> 0x8005
func reverse(v uint16) uint16 {
	r := v
	s := uint(16 - 1)

	for v >>= 1; v != 0; v >>= 1 {
		r <<= 1
		r |= v & 1
		s--
	}
	r <<= s
	return r
}

// makeTable claculates the polynomial table for the given
// configuration structure.
func (c *Conf) makeTable() {
	if c.BitRev {
		c.table = MakeTable(reverse(c.Poly))
		c.update = Update
	} else {
		c.table = MakeTableNBR(c.Poly)
		c.update = UpdateNBR
	}
}

// Table is a 256-word table representing the polynomial for efficient
// processing.
type Table [256]uint16

// MakeTable returns the Table constructed from the specified
// polynomial. The table is calcuated in bit-reversed order (bit-15
// corresponds to the X^0 term). Argument poly must be given
// bit-reversed (e.g. 0xA001 for the 0x8005 polynomial).
func MakeTable(poly uint16) *Table {
	t := new(Table)
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t[i] = crc
	}
	return t
}

// MakeTableNBR returns the Table constructed from the specified
// polynomial. The table is calculated in non-bit-reversed order
// (bit-0 corresponds to the X^0 term).
func MakeTableNBR(poly uint16) *Table {
	t := new(Table)
	for i := 0; i < 256; i++ {
		crc := uint16(i) << 8
		for j := 0; j < 8; j++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t[i] = crc
	}
	return t
}

// Update returns the CRC-16 checksum of p using the polynomial table
// tab constructed by MakeTable (bit-reversed order). The resulting
// CRC is in bit-reversed order (bit-15 corresponds to the X^0
// term). Argument crc is the initial value of the CRC register.
func Update(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab[byte(crc)^v] ^ (crc >> 8)
	}
	return crc
}

// UpdateNBR returns the CRC-16 checksum of p using the polynomial
// table tab constructed by MakeTableNBR (non-bit-reversed order). The
// resulting CRC is in non-bit-reversed order (bit-0 corresponds to
// the X^0 term). Argument CRC is the initial value of the CRC
// register.
func UpdateNBR(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	crc  uint16
	conf *Conf
}

// NOTE(npat): Eventually should be moved in <stdlib>/hash

// Hash16 is the common interface implemented by all 16-bit hash
// functions.
type Hash16 interface {
	hash.Hash
	Sum16() uint16
}

// New creates a new hash.Hash16 computing the CRC-16 checksum using
// the configuration c.
func New(c *Conf) Hash16 {
	c.once.Do(c.makeTable)
	return &digest{crc: c.IniVal, conf: c}
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Reset() { d.crc = d.conf.IniVal }

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = d.conf.update(d.crc, d.conf.table, p)
	return len(p), nil
}

func (d *digest) Sum16() uint16 { return d.crc ^ d.conf.FinVal }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum16()
	if d.conf.BigEnd {
		return append(in, byte(s>>8), byte(s))
	} else {
		return append(in, byte(s), byte(s>>8))
	}
}

// Checksum returns the CRC-16 checksum of data using the
// configuration c.
func Checksum(c *Conf, data []byte) uint16 {
	c.once.Do(c.makeTable)
	return c.update(c.IniVal, c.table, data) ^ c.FinVal
}

// See also:
//   http://en.wikipedia.org/wiki/Computation_of_cyclic_redundancy_checks
//   https://www.kernel.org/doc/Documentation/crc32.txt
//   http://www.zlib.net/crc_v3.txt
//
