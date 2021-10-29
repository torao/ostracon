// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Modified from original GoGo Protobuf to return number of bytes written.

package protoio

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/line/ostracon/proto/ostracon/p2p"
	"github.com/line/ostracon/proto/ostracon/types"
	"io"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
)

// NewDelimitedWriter writes a varint-delimited Protobuf message to a writer. It is
// equivalent to the gogoproto NewDelimitedWriter, except WriteMsg() also returns the
// number of bytes written, which is necessary in the p2p package.
func NewDelimitedWriter(w io.Writer) WriteCloser {
	return &varintWriter{w, make([]byte, binary.MaxVarintLen64), nil}
}

type varintWriter struct {
	w      io.Writer
	lenBuf []byte
	buffer []byte
}

func (w *varintWriter) WriteMsg(msg proto.Message) (int, error) {
	if m, ok := msg.(marshaler); ok {
		n, ok := getSize(m)
		if ok {
			if n+binary.MaxVarintLen64 >= len(w.buffer) {
				w.buffer = make([]byte, n+binary.MaxVarintLen64)
			}
			lenOff := binary.PutUvarint(w.buffer, uint64(n))
			_, err := m.MarshalTo(w.buffer[lenOff:])
			if err != nil {
				return 0, err
			}
			_, err = w.w.Write(w.buffer[:lenOff+n])
			logMessage(">>>", msg, w.buffer[:lenOff], w.buffer[lenOff:lenOff+n])
			return lenOff + n, err
		}
	}

	// fallback
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, err
	}
	length := uint64(len(data))
	n := binary.PutUvarint(w.lenBuf, length)
	_, err = w.w.Write(w.lenBuf[:n])
	if err != nil {
		return 0, err
	}
	_, err = w.w.Write(data)
	logMessage(">>>", msg, w.lenBuf[:n], data)
	return len(data) + n, err
}

func (w *varintWriter) Close() error {
	if closer, ok := w.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func MarshalDelimited(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	_, err := NewDelimitedWriter(&buf).WriteMsg(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func logMessage(io string, msg proto.Message, varLen []byte, data []byte){
	var name string
	if t := reflect.TypeOf(msg); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}
	switch name {
	case "Packet":
		packet := msg.(*p2p.Packet)
		if ping := packet.GetPacketPing(); ping != nil {
			name += ".Ping"
		} else if pong := packet.GetPacketPong(); pong != nil {
			name += ".Pong"
		} else if msg := packet.GetPacketMsg(); msg != nil {
			data := hex.EncodeToString(msg.Data)
			if len(data) > 8 {
				data = data[:8] + "..."
			}
			name += fmt.Sprintf(".Msg(%d,%v,%s)", msg.GetChannelID(), msg.GetEOF(), data)
		} else {
			name += fmt.Sprintf(".%d", packet.Sum.Size())
		}
		break
	case "CanonicalVote":
		canonicalVote := msg.(*types.CanonicalVote)
		msgType := types.SignedMsgType_name[int32(canonicalVote.GetType())]
		if strings.HasPrefix(msgType, "SIGNED_MSG_TYPE_") {
			msgType = msgType[len("SIGNED_MSG_TYPE_"):]
		}
		name += fmt.Sprintf("(%s,%d,%d,%X,%s,%s)", msgType, canonicalVote.GetHeight(), canonicalVote.GetRound(), canonicalVote.GetBlockID().GetHash()[:4], canonicalVote.GetTimestamp().Format("20060102030405MST"), canonicalVote.GetChainID())
	}
	length := len(varLen) + len(data)
	fmt.Printf("%s %s (%d) %X %X\n", io, name, length, varLen, data)
}
