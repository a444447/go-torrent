package Client

import (
	"bytes"
	"fmt"
	"go-torrent/src/BitField"
	"go-torrent/src/HandShake"
	"go-torrent/src/Message"
	"net"
	"time"
)

// Client Client是与peer的TCP连接
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield BitField.Bitfield
	peer     HandShake.PeerInfo
	infoHash [20]byte
	peerID   [20]byte
}

// 进行一次握手检查peer能否完成我们的要求
func checkHandshake(conn net.Conn, infoHash, peerID [20]byte) (*HandShake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second)) // 设置Deadline time 超过这个时间不会再读取数据
	defer conn.SetDeadline(time.Time{})               //手动重置这个值

	req := HandShake.New(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	res, err := HandShake.Read(conn)
	if err != nil {
		return nil, err
	}
	//比较服务器返回的infoHash和我们要下载的文件的infoHash是否一致
	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infoHash)
	}

	return res, nil

}

func recvBitField(conn net.Conn) (BitField.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := Message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("excepted bitfield but got %s", msg)
	}
	if msg.ID != Message.MsgBitfield {
		return nil, fmt.Errorf("expected bitfield but got ID %d", msg.ID)
	}

	return msg.Payload, nil
}

func NewDial(peer HandShake.PeerInfo, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second) // 建立一个TCP连接，并且设定过期时间为3s,这样就不会浪费过多时间在不允许连接的peer上
	if err != nil {
		return nil, err
	}
	_, err = checkHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := recvBitField(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

// 从connection中读取消息
func (c *Client) Read() (*Message.Message, error) {
	msg, err := Message.Read(c.Conn)
	return msg, err
}

// SendRequest 向peer发送Request Message
func (c *Client) SendRequest(index, begin, length int) error {
	req := Message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested Message to the peer
func (c *Client) SendInterested() error {
	msg := Message.Message{ID: Message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested Message to the peer
func (c *Client) SendNotInterested() error {
	msg := Message.Message{ID: Message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke Message to the peer
func (c *Client) SendUnchoke() error {
	msg := Message.Message{ID: Message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have Message to the peer
func (c *Client) SendHave(index int) error {
	msg := Message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
