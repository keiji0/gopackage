package node

import "github.com/keiji0/sandbox/201904_p2p/net/message"

// Node はP2Pのノード
type Node struct {
	// connect 通信しているコネクション
	connect Connect
}

// NewNode はノードを生成する
func NewNode() *Node {
	node := &Node{}
	return node
}

// Send はノードにメッセージを送信する
func (n *Node) Send(msg *message.Message) *message.Message {
	return nil
}
