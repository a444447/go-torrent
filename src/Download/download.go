// Package Download Download包负责处理同时与多个 peers 进行通信时的并发，以及管理这些正在与我们进行交互的 peers 的状态
package Download

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"go-torrent/src/Client"
	"go-torrent/src/HandShake"
	"go-torrent/src/Message"
	"log"
	"runtime"
	"time"
)

const (
	MaxBlockSize = 16384 //MaxBlockSize 是一个 request 所能请求的最大字节数
	MaxBacklog   = 5     //MaxBacklog 是管道中最多存储的未完成的请求数量
)

// TorrentData 一系列数据包括torrent找到的peer
type TorrentData struct {
	Peers       []HandShake.PeerInfo
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHash   [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// pieceProgress 该结构用来跟踪每一个peer, peer 那里下载了哪些数据、我们向指定 peer 请求了哪些数据和我们是否处于 choked 状态
type pieceProgress struct {
	index      int
	client     *Client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

// revStateFromMsg 读取消息确定当前peer的状态
func (s *pieceProgress) revStateFromMsg() error {
	msg, err := s.client.Read() //请求消息
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}

	switch msg.ID {
	case Message.MsgUnchoke:
		s.client.Choked = false
	case Message.MsgChoke:
		s.client.Choked = true
	case Message.MsgHave:
		index, err := Message.ParseHave(msg)
		if err != nil {
			return err
		}
		s.client.Bitfield.SetPiece(index)
	case Message.MsgPiece:
		length, err := Message.ParsePiece(s.index, s.buf, msg)
		if err != nil {
			return err
		}
		s.downloaded += length
		s.backlog--
	}
	return nil

}

// 管道化策略。并且我们将每个作为block进行下载
func pieceDownload(c *Client.Client, pw *pieceWork) ([]byte, error) {
	//设定现在的peer状态
	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}
	// 设置一个Deadline防止一个peer无响应而导致堵塞
	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{}) //Deadline不会自己重置，需要手动重置

	for state.downloaded < pw.length {
		// 当peer处于unchoked状态，就不断发送请求直到管道已满
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := MaxBlockSize
				//注意: 最后一个block的大小可能和典型的block大小不一样
				if pw.length-state.requested < blockSize {
					blockSize = pw.length - state.requested
				}

				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++ //管道数量又增加
				state.requested += blockSize

			}
		}

		err := state.revStateFromMsg()
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil

}

// 检查下载的buf内容与piece的hash是否相同
func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("index %d failed integrity check", pw.index)
	}

	return nil
}

func (t *TorrentData) startDownloadWorker(peer HandShake.PeerInfo, workQueue chan *pieceWork, results chan *pieceResult) {
	c, err := Client.NewDial(peer, t.PeerID, t.InfoHash) //与peer建立tcp链接
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer c.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke() //告诉peer unchoke信息，才能开始发送
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw //如果该peer中没有想要的piece，就将其放回workQueue中
			continue
		}

		//开始下载Piece
		buf, err := pieceDownload(c, pw)
		if err != nil {
			log.Println("exiting", err)
			workQueue <- pw
			return
		}

		err = checkIntegrity(pw, buf)
		if err != nil {
			log.Printf("piece #%d failed integrity check\n", pw.index)
			workQueue <- pw
			continue
		}

		c.SendHave(pw.index) //告诉peer我们拥有找个piece
		results <- &pieceResult{pw.index, buf}
	}

}

func (t *TorrentData) calculateBounds(index int) (begin, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

// Download 下载torrent.并且将整个文件存储在内存中
func (t *TorrentData) Download() ([]byte, error) {
	log.Println("start download:", t.Name)
	//现在我们用channels来同步我们的两个worker
	workQueue := make(chan *pieceWork, len(t.PieceHash)) // 用来从不同的peer中下载piece
	results := make(chan *pieceResult)                   //收集已经下载的Piece,将其放入缓存buf中
	for index, hash := range t.PieceHash {
		begin, end := t.calculateBounds(index)
		length := end - begin
		workQueue <- &pieceWork{index, hash, length}
	}
	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results) //为每个peer创建一个goroutine
	}

	//将下载的结果装在一个buffer里面
	buf := make([]byte, t.Length)
	done := 0 // 已经下载完成的piece
	for done < len(t.PieceHash) {
		res := <-results
		begin, end := t.calculateBounds(res.index)
		copy(buf[begin:end], res.buf)
		done++

		percent := float64(done) / float64(len(t.PieceHash)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // 当前正在运行的worker数量，减去1是减去main goroutine
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	close(workQueue)

	return buf, nil

}
