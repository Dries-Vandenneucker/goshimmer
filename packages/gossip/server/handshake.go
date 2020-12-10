package server

import (
	"bytes"
	"time"

	pb "github.com/iotaledger/goshimmer/packages/gossip/server/proto"
	"github.com/iotaledger/hive.go/autopeering/server"
	"google.golang.org/protobuf/proto"
)

const (
	versionNum          = 0
	handshakeExpiration = 20 * time.Second
)

// isExpired checks whether the given UNIX time stamp is too far in the past.
func isExpired(ts int64) bool {
	return time.Since(time.Unix(ts, 0)) >= handshakeExpiration
}

func newHandshakeRequest(toAddr string) ([]byte, error) {
	m := &pb.HandshakeRequest{
		Version:   versionNum,
		To:        toAddr,
		Timestamp: time.Now().Unix(),
	}
	return proto.Marshal(m)
}

func newHandshakeResponse(reqData []byte) ([]byte, error) {
	m := &pb.HandshakeResponse{
		ReqHash: server.PacketHash(reqData),
	}
	return proto.Marshal(m)
}

func (t *TCP) validateHandshakeRequest(reqData []byte) bool {
	m := new(pb.HandshakeRequest)
	if err := proto.Unmarshal(reqData, m); err != nil {
		t.log.Infow("invalid handshake",
			"err", err,
		)
		return false
	}
	if m.GetVersion() != versionNum {
		t.log.Infow("invalid handshake",
			"version", m.GetVersion(),
			"want", versionNum,
		)
		return false
	}
	if isExpired(m.GetTimestamp()) {
		t.log.Infow("invalid handshake",
			"timestamp", time.Unix(m.GetTimestamp(), 0),
		)
	}

	return true
}

func (t *TCP) validateHandshakeResponse(resData []byte, reqData []byte) bool {
	m := new(pb.HandshakeResponse)
	if err := proto.Unmarshal(resData, m); err != nil {
		t.log.Infow("invalid handshake",
			"err", err,
		)
		return false
	}
	if !bytes.Equal(m.GetReqHash(), server.PacketHash(reqData)) {
		t.log.Infow("invalid handshake",
			"hash", m.GetReqHash(),
		)
		return false
	}

	return true
}
