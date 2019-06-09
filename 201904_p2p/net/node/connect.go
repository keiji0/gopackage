package node

// Connect はNode間をつなぐコネクター
// メッセージを交換できるようにソケットを抽象化されたクラス
type Connect struct {
}

// NewConnect はコネクターを生成する
func NewConnect() *Connect {
	conn := &Connect{}
	return conn
}
