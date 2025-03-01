package astits

import (
	"bytes"
	"testing"

	"github.com/icza/bitio"
	"github.com/stretchr/testify/assert"
)

func TestParseData(t *testing.T) {
	// Init
	pm := newProgramMap()
	ps := []*Packet{}

	// Custom parser
	cds := []*DemuxerData{{PID: 1}}
	c := func(ps []*Packet) (o []*DemuxerData, skip bool, err error) {
		o = cds
		skip = true
		return
	}
	ds, err := parseData(ps, c, pm)
	assert.NoError(t, err)
	assert.Equal(t, cds, ds)

	// Do nothing for CAT
	ps = []*Packet{{Header: &PacketHeader{PID: PIDCAT}}}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Empty(t, ds)

	// PES
	p := pesWithHeaderBytes()
	ps = []*Packet{
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[:33],
		},
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[33:],
		},
	}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Equal(t, []*DemuxerData{{FirstPacket: ps[0], PES: pesWithHeader(), PID: uint16(256)}}, ds)

	// PSI
	pm.set(uint16(256), uint16(1))
	p = psiBytes()
	ps = []*Packet{
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[:33],
		},
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[33:],
		},
	}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Equal(t, psi.toData(ps[0], uint16(256)), ds)
}

func TestIsPSIPayload(t *testing.T) {
	pm := newProgramMap()
	var pids []int
	for i := 0; i <= 255; i++ {
		if isPSIPayload(uint16(i), pm) {
			pids = append(pids, i)
		}
	}
	assert.Equal(t, []int{0, 16, 17, 18, 19, 20, 30, 31}, pids)
	pm.set(uint16(1), uint16(0))
	assert.True(t, isPSIPayload(uint16(1), pm))
}

func TestIsPESPayload(t *testing.T) {
	buf := &bytes.Buffer{}
	w := bitio.NewWriter(buf)
	WriteBinary(w, "0000000000000001")
	assert.False(t, isPESPayload(buf.Bytes()))
	buf.Reset()
	WriteBinary(w, "000000000000000000000001")
	assert.True(t, isPESPayload(buf.Bytes()))
}
