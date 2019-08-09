package astits

import (
	"testing"

	"github.com/asticode/go-astitools/binary"
	"github.com/stretchr/testify/assert"
	"github.com/asticode/go-astitools/byte"
)

var eit = &EITData{
	Events: []*EITDataEvent{{
		Descriptors:    descriptors,
		Duration:       dvbDurationSeconds,
		EventID:        6,
		HasFreeCSAMode: true,
		RunningStatus:  7,
		StartTime:      dvbTime,
	}},
	LastTableID:              5,
	OriginalNetworkID:        3,
	SegmentLastSectionNumber: 4,
	ServiceID:                1,
	TransportStreamID:        2,
}

func eitBytes() []byte {
	w := astibinary.New()
	w.Write(uint16(2))               // Transport stream ID
	w.Write(uint16(3))               // Original network ID
	w.Write(uint8(4))                // Segment last section number
	w.Write(uint8(5))                // Last table id
	w.Write(uint16(6))               // Event #1 id
	w.Write(dvbTimeBytes)            // Event #1 start time
	w.Write(dvbDurationSecondsBytes) // Event #1 duration
	w.Write("111")                   // Event #1 running status
	w.Write("1")                     // Event #1 free CA mode
	descriptorsBytes(w)              // Event #1 descriptors
	return w.Bytes()
}

func TestParseEITSection(t *testing.T) {
	var b = eitBytes()
	d, err := parseEITSection(astibyte.NewIterator(b), len(b), uint16(1))
	assert.Equal(t, d, eit)
	assert.NoError(t, err)
}
