package message

import (
	"bytes"
	"io"
	"reflect"

	"github.com/keiji0/sandbox/201904_p2p/crypto/hash"
	"github.com/keiji0/sandbox/201904_p2p/net"
	"github.com/keiji0/sandbox/201904_p2p/net/internal"
	"github.com/pkg/errors"
)

const (
	// HelloType はHelloMessageのタイプ
	HelloType net.MessageType = 1
)

// メッセージコマンドの一覧
var messages = []Message{
	&Hello{},
}

// Message はノードへ送信するコマンドのインターフェース
// 各コマンドはこのインターフェースを実装します
type Message interface {
	// Messageのタイプを返す
	MessageType() net.MessageType
	// MessageのPayloadをシリアライズする
	Serialize(w io.Writer) error
	// MessageのPayloadをデシリアライズする
	Deserialize(r io.Reader) error
}

// 送信するメッセージのヘッダー
type messageHeader struct {
	// magic はメッセージを判別するための識別子
	magic net.MessageMagic
	// messageType はメッセージのタイプ
	messageType net.MessageType
	// length はメッセージの長さ
	length uint32
	// checksum はメッセージのチェックサム
	checksum [net.MessageChecksumSize]byte
}

// メッセージとMessageTypeのマップ、ちょっとでも早くアクセスするため
var messageMap = map[net.MessageType]reflect.Type{}

func init() {
	for _, msg := range messages {
		messageMap[msg.MessageType()] = reflect.TypeOf(msg).Elem()
	}
}

// Send はメッセージをネットワークに送信します
func Send(w io.Writer, msg Message) (err error) {
	h := messageHeader{}

	// NetTypeからmagicを取得
	h.magic = net.MessageMagicNumber

	// ProtocolにあったCommandに変換
	h.messageType = msg.MessageType()

	// コマンド本体のバイト列を取得
	payload := &bytes.Buffer{}
	if err := msg.Serialize(payload); err != nil {
		return err
	}

	h.length = uint32(len(payload.Bytes()))
	if net.MessageMaxSize < h.length {
		return errors.Errorf("MessageのPayloadのサイズが規定値より大きいです: type=%d, size=%d", msg.MessageType(), h.length)
	}

	copy(h.checksum[:], hash.Sha256x2(payload.Bytes())[0:net.MessageChecksumSize])

	// メッセージヘッダーを送信
	if err := internal.BulkSerialize(w, h.magic, h.messageType, h.length, h.checksum); err != nil {
		return err
	}

	// Payloadを送信
	if _, err := w.Write(payload.Bytes()); err != nil {
		return errors.Wrapf(err, "Payloadの送信に失敗しました")
	}

	return nil
}

// Receive はネットワークからメッセージを受信します
func Receive(r io.Reader) (Message, error) {
	h, err := readMessageHeader(r)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, h.length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, errors.Wrapf(err, "Payloadの読み込みに失敗しました: type=%d", h.messageType)
	}

	checksum := hash.Sha256x2(payload)[0:net.MessageChecksumSize]
	if !bytes.Equal(checksum, h.checksum[:]) {
		return nil, errors.Errorf("Payloadのチェックサムが一致しません: type=%d", h.messageType)
	}

	msg, err := newMessage(h.messageType)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(payload)
	if err := msg.Deserialize(buf); err != nil {
		return nil, err
	}

	return msg, nil
}

// newMessage は指定したコマンド名のメッセージを生成します
func newMessage(messageType net.MessageType) (Message, error) {
	t, ok := messageMap[messageType]
	if !ok {
		return nil, errors.Errorf("メッセージタイプが見つかりませんでした: %d", messageType)
	}
	i, ok := reflect.New(t).Interface().(Message)
	if !ok {
		return nil, errors.Errorf("Message Interfaceにキャストできませんでした: %q", t)
	}
	return i, nil
}

// readMessageHeader はMessageHeaderをネットワークから読み込みます
func readMessageHeader(r io.Reader) (*messageHeader, error) {
	h := &messageHeader{}

	if err := internal.BulkDeserialize(r, &h.magic, &h.messageType, &h.length, &h.checksum); err != nil {
		return nil, err
	}

	if net.MessageMagicNumber != h.magic {
		return nil, errors.Errorf("Messageのマジックナンバーが一致しません")
	}

	if net.MessageMaxSize < h.length {
		return nil, errors.Errorf("MessageのPayloadのサイズが規定値より大きいです: messageType=%d, size=%d", h.messageType, h.length)
	}

	return h, nil
}
