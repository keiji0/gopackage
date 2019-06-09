package message

import (
	"io"

	"github.com/keiji0/sandbox/201904_p2p/net"
	"github.com/keiji0/sandbox/201904_p2p/net/internal"
)

// Hello ははじめにノードに接続するためのメッセージ
// ノード接続時にお互いのノードがHelloメッセージを送りハンドシェイクを行う
type Hello struct {
	// Version は送り主がサポートしているプロトコルのバージョンが格納されている
	ProtocolVersion net.Version
}

// MessageType はこのメッセージのコマンド名を返します
func (m *Hello) MessageType() net.MessageType {
	return HelloType
}

// NewHello はMsgVersionメッセージを生成します
func NewHello() *Hello {
	v := &Hello{
		ProtocolVersion: net.CurrentVersion,
	}
	return v
}

// Serialize はMessageのPayloadをシリアライズする
func (m *Hello) Serialize(w io.Writer) error {
	return internal.BulkSerialize(
		w,
		m.ProtocolVersion,
	)
}

// Deserialize はMessageのPayloadをデシリアライズする
func (m *Hello) Deserialize(r io.Reader) error {
	return internal.BulkDeserialize(
		r,
		&m.ProtocolVersion,
	)
}
