package Torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
	"go-torrent/src/Download"
	"io"
	"math/rand"
	"os"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// TorrentFile 一个更加平坦的结构
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte   //整个Bencode 编码后的 info 字典（包含了name、size 和 piece hashes）的 SHA-1 哈希值;在tracker和peer通信时唯一表示我们希望下载的文件
	PieceHashes [][20]byte //将pieces转换为长度为20字节的哈希值切片， 方便我们的访问
	PieceLength int
	FileLength  int
	FileName    string
}

// ParseTorrentFile 将输入流中的bencode内容全部转换为TorrentFile格式
//大体过程为 Parse(input)->Unmarshal(). Unmarshal()后就已经将bencode信息填充到bencodeTorrent结构中了,然后我们再赋给我们的扁平化结构TorrentFile
func ParseTorrentFile(r io.Reader) (*TorrentFile, error) {
	raw := new(bencodeTorrent)       // new一个bencodeInfo机构的指针
	err := bencode.Unmarshal(r, raw) // 将r中的内容填充到raw结构
	if err != nil {
		fmt.Println("Fail Parse torrent File")
		return nil, err
	}

	ret := new(TorrentFile) //new一个TorrentFile结构指针
	ret.Announce = raw.Announce
	ret.FileName = raw.Info.Name
	ret.FileLength = raw.Info.Length
	ret.PieceLength = raw.Info.PieceLength

	//计算整个文件的SHA-1哈希值
	buf := new(bytes.Buffer)
	err = bencode.Marshal(buf, raw.Info)
	if err != nil {
		fmt.Println("raw file info error")
	}
	ret.InfoHash = sha1.Sum(buf.Bytes())

	//计算pieces哈希
	bytePieces := []byte(raw.Info.Pieces)
	cnt := len(bytePieces) / 20
	hashes := make([][20]byte, cnt)
	for i := 0; i < cnt; i++ {
		copy(hashes[i][:], bytePieces[i*20:(i+1)*20])
	}
	ret.PieceHashes = hashes
	return ret, nil
}

// DownloadToFile 下载torrent文件并且将其写到文件中
func (t *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := t.RequestPeers(peerID, 6881)
	torrent := Download.TorrentData{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHash:   t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.FileLength,
		Name:        t.FileName,
	}

	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	//bencodeTo := new(bencodeTorrent)
	//err = bencode.Unmarshal(file, bencodeTo)
	//if err != nil {
	//	return TorrentFile{}, err
	//}

	finalTorrent, _ := ParseTorrentFile(file)

	return *finalTorrent, err

}
