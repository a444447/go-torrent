package HandShake

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

// Handshake 该结构用来用于确认peer是否满足我们的握手要求
type Handshake struct {
	Pstr     string   //协议标识符，称为 pstr，它总是 BitTorrent protocol
	InfoHash [20]byte // 用于标识我们希望下载的文件
	PeerID   [20]byte //表示我们自己的PeerID
}

// PeerInfo Peer的结构体
type PeerInfo struct {
	IP   net.IP
	Port uint16
}

func (p PeerInfo) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

// Serialize 返回一个用于握手字符串
/*
字符串格式:
1.协议标识符的长度，它总是 19（用十六进制表示是 0x13）
2.协议标识符，称为 pstr，它总是 BitTorrent protocol
3.8 个保留字节，全部设置为 0。我们可以将它们当中的某些设置为 1 来表示我们支持特定的扩展。但是我们不需要，所以我们只把它们设置为 0 就行
4.我们之前计算过的用于标识我们希望下载的文件的 infohash
5.用于表示我们自己的 Peer ID
*/
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	currPosition := 1
	currPosition += copy(buf[currPosition:], h.Pstr)
	currPosition += copy(buf[currPosition:], make([]byte, 8))
	currPosition += copy(buf[currPosition:], h.InfoHash[:])
	currPosition += copy(buf[currPosition:], h.PeerID[:])
	return buf
}

// Read 处理从服务器回复(与我们发送的握手字符串格式相似)
func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthBuf[0]) // 协议标识符长度, 是19
	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, pstrlen+48)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20]) // 8是保留字节
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	return &h, nil

}
