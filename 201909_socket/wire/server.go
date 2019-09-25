package wire

import (
	"context"
	"log"
	"net"
	"time"

	proto "github.com/keiji0/sandbox/201909_socket/wire/protocol"
	pb "github.com/keiji0/sandbox/201909_socket/wire/protocol/pb"
	"google.golang.org/grpc"
)

// Server はWire上でサービスを提供します。
type Server struct {
	// このサーバーのアドレス
	addr string
	// このノードに対するメッセージチャネル
	MessageChannel MessageChannel
	// 接続しているNode
	connectedNodes []*Node
}

// Node はWire上の要素
type Node struct {
	// このノードのアドレス
	addr string
	// 接続しているコネクション
	conn *grpc.ClientConn
	// ノードへのクライアント
	client pb.NodeClient
}

// Message はNode間のメッセージに使われるデータ
type Message interface {
}

// MessageChannel はNode間のメッセージのやり取りするためのチャネル
type MessageChannel chan Message

// NewServer は新たにServerを生成します
func NewServer(addr string) *Server {
	serv := Server{
		addr:           addr,
		MessageChannel: make(MessageChannel),
	}
	return &serv
}

// newNode あ新たにNodeを生成します
func newNode(addr string) *Node {
	node := &Node{
		addr: addr,
	}
	return node
}

// Listen サーバーのポートを開き待ち状態になります
func (serv *Server) Listen() error {
	lis, err := net.Listen("tcp", serv.addr)
	if err != nil {
		return err
	}
	gserv := grpc.NewServer()
	pb.RegisterNodeServer(gserv, serv)
	log.Printf("%s Start Server", serv.addr)
	if err := gserv.Serve(lis); err != nil {
		return err
	}
	return nil
}

// Connect は指定ノードに接続します
func (serv *Server) Connect(addr string) error {
	// ノードを生成し、接続する
	node := newNode(addr)
	if err := node.connect(); err != nil {
		return err
	}
	// 接続に成功した場合は接続済みノードリストへ追加
	serv.connectedNodes = append(serv.connectedNodes, node)
	return nil
}

// Hello implements helloworld.GreeterServer
func (serv *Server) Hello(ctx context.Context, in *pb.Version) (*pb.Version, error) {
	log.Printf("%s Hello: %v", serv.addr, in)
	return &pb.Version{Version: proto.Version}, nil
}

// connect はこのノードに接続する
func (node *Node) connect() error {
	// ノードの接続
	conn, err := grpc.Dial(node.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	node.conn = conn
	node.client = pb.NewNodeClient(node.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := node.client.Hello(ctx, &pb.Version{Version: 2})
	if err != nil {
		conn.Close()
		return err
	}
	log.Printf("%s Receve Hello: %v", node.addr, r)
	return nil
}
