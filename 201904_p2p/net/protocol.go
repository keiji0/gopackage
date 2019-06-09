package net

import "encoding/binary"

// MessageMagic はネットワークで利用するマジックタイプを表す型
type MessageMagic uint32

// Version はプロトコルのバージョンを表す型
type Version int32

// ByteOrder はProtocolで使用するバイトオーダーを指定
var ByteOrder = binary.LittleEndian

// MessageMagicNumber はメッセージのマジックナンバー
const MessageMagicNumber MessageMagic = 0xfefefefe

// CurrentVersion はサポートしているプロトコルのバージョンになります
const CurrentVersion Version = 1

// MaxStringLength はStringのサイズの最大値
const MaxStringLength = 0xffff

// VarUint は可変長の符号なし数値を表す型
type VarUint uint64

// 可変長数値の識別子を定義
const (
	VarUint8Max        = 0xfc
	VarUint16Tag uint8 = 0xfd
	VarUint32Tag uint8 = 0xfe
	VarUint64Tag uint8 = 0xff
)

// MessageType はメッセージのタイプを表す型
type MessageType uint8

// MessageChecksumSize はメッセージ内のチェックサムのバイトサイズ
const MessageChecksumSize = 4

// MessageMaxSize はメッセージの最大サイズ
const MessageMaxSize = 0x02000000
