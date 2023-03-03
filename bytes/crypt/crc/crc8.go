package crc

import "math/bits"

// Param8 represents parameters of a CRC-8 algorithm including polynomial and initial value.
// More information about algorithms parametrization and parameter descriptions
// can be found here - http://www.zlib.net/crc_v3.txt
type Param8 struct {
	Poly   uint8
	Init   uint8
	RefIn  bool
	RefOut bool
	XorOut uint8
	Check  uint8
	Name   string
}

// Predefined CRC-8 algorithms.
// List of algorithms with their parameters borrowed from here - http://reveng.sourceforge.net/crc-catalogue/1-15.htm#crc.cat-bits.8
//
// The variables can be used to create Table for the selected algorithm.
var (
	CRC8          = Param8{0x07, 0x00, false, false, 0x00, 0xF4, "CRC-8"}
	CRC8_CDMA2000 = Param8{0x9B, 0xFF, false, false, 0x00, 0xDA, "CRC-8/CDMA2000"}
	CRC8_DARC     = Param8{0x39, 0x00, true, true, 0x00, 0x15, "CRC-8/DARC"}
	CRC8_DVB_S2   = Param8{0xD5, 0x00, false, false, 0x00, 0xBC, "CRC-8/DVB-S2"}
	CRC8_EBU      = Param8{0x1D, 0xFF, true, true, 0x00, 0x97, "CRC-8/EBU"}
	CRC8_I_CODE   = Param8{0x1D, 0xFD, false, false, 0x00, 0x7E, "CRC-8/I-CODE"}
	CRC8_ITU      = Param8{0x07, 0x00, false, false, 0x55, 0xA1, "CRC-8/ITU"}
	CRC8_MAXIM    = Param8{0x31, 0x00, true, true, 0x00, 0xA1, "CRC-8/MAXIM"}
	CRC8_ROHC     = Param8{0x07, 0xFF, true, true, 0x00, 0xD0, "CRC-8/ROHC"}
	CRC8_WCDMA    = Param8{0x9B, 0x00, true, true, 0x00, 0x25, "CRC-8/WCDMA"}
)

// Table8 is a 256-byte table representing polynomial and algorithm settings for efficient processing.
type Table8 struct {
	params Param8
	data   [256]uint8
}

// MakeTable8 returns the Table constructed from the specified algorithm.
func MakeTable8(params Param8) *Table8 {
	table := new(Table8)
	table.params = params
	for n := 0; n < 256; n++ {
		crc := uint8(n)
		for i := 0; i < 8; i++ {
			bit := (crc & 0x80) != 0
			crc <<= 1
			if bit {
				crc ^= params.Poly
			}
		}
		table.data[n] = crc
	}
	return table
}

// Init8 returns the initial value for CRC register corresponding to the specified algorithm.
func Init8(table *Table8) uint8 {
	return table.params.Init
}

// Update8 returns the result of adding the bytes in data to the crc.
func Update8(crc uint8, data []byte, table *Table8) uint8 {
	if table.params.RefIn {
		for _, d := range data {
			d = bits.Reverse8(d)
			crc = table.data[crc^d]
		}
	} else {
		for _, d := range data {
			crc = table.data[crc^d]
		}
	}
	return crc
}

// Complete8 returns the result of CRC calculation and post-calculation processing of the crc.
func Complete8(crc uint8, table *Table8) uint8 {
	if table.params.RefOut {
		crc = bits.Reverse8(crc)
	}

	return crc ^ table.params.XorOut
}

// Checksum8 returns CRC checksum of data using specified algorithm represented by the Table.
func Checksum8(data []byte, table *Table8) uint8 {
	crc := Init8(table)
	crc = Update8(crc, data, table)
	return Complete8(crc, table)
}
