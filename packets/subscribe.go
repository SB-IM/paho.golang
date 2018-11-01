package packets

import (
	"bytes"
	"io"
	"net"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	PacketID      uint16
	Properties    *Properties
	Subscriptions map[string]SubOptions
}

// SubOptions is the struct representing the options for a subscription
type SubOptions struct {
	QoS               byte
	NoLocal           bool
	RetainAsPublished bool
	RetainHandling    byte
}

// Pack is the implementation of the interface required function for a packet
func (s *SubOptions) Pack() byte {
	var ret byte
	ret |= s.QoS & 0x03
	if s.NoLocal {
		ret |= 1 << 2
	}
	if s.RetainAsPublished {
		ret |= 1 << 3
	}
	ret |= s.RetainHandling & 0x30

	return ret
}

// Unpack is the implementation of the interface required function for a packet
func (s *Subscribe) Unpack(r *bytes.Buffer) error {
	var err error
	s.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}

	err = s.Properties.Unpack(r, SUBSCRIBE)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (s *Subscribe) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(s.PacketID, &b)
	var subs bytes.Buffer
	for t, o := range s.Subscriptions {
		writeString(t, &subs)
		subs.WriteByte(o.Pack())
	}
	idvp := s.Properties.Pack(SUBSCRIBE)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp, subs.Bytes()}
}

// WriteTo is the implementation of the interface required function for a packet
func (s *Subscribe) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: SUBSCRIBE, Flags: 2}}
	cp.Content = s

	return cp.WriteTo(w)
}