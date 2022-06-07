package Torrent

import (
	"encoding/binary"
	"fmt"
	"github.com/jackpal/bencode-go"
	"go-torrent/src/HandShake"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce) //将string解析为*url格式
	if err != nil {
		fmt.Println("Announce error :", t.Announce)
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},      // 20 字节, 将 .torrent 文件中的 info 键对应的值生成的 SHA1 哈希, 该哈希值可作为所要请求的资源的标识符
		"peer_id":    []string{string(peerID[:])},          //终端生成的 20 个字符的唯一标识符, 每个进行 BT 下载的终端随机生成的 20 个字符的字符串作为其标识符 (终端应在每次开始一个新的下载任务时重新随机生成一个新的 peer_id)
		"port":       []string{strconv.Itoa(int(6881))},    //该终端正在监听的端口 (因为 BT 协议是 P2P 的, 所以每一个下载终端也都会暴露一个端口, 供其它结点下载)
		"uploaded":   []string{"0"},                        //当前已经上传的文件的字节数 (十进制数字表示)
		"downloaded": []string{"0"},                        //当前已经下载的文件的字节数 (十进制数字表示)
		"compact":    []string{"1"},                        //
		"left":       []string{strconv.Itoa(t.FileLength)}, //当前仍需要下载的文件的字节数 (十进制数字表示)
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

// 从buffer中将peer的IP地址与port unmarshal
func buildPeerInfo(peersBuffer []byte) ([]HandShake.PeerInfo, error) {
	const peerSize = 6                      // 4(IP length) + 2(Port length)
	numsPeer := len(peersBuffer) / peerSize // 求peer的数量
	fmt.Println(numsPeer)
	if len(peersBuffer)%peerSize != 0 {
		err := fmt.Errorf("received wrong form peersBuffer")
		return nil, err
	}
	peers := make([]HandShake.PeerInfo, numsPeer)
	for i := 0; i < numsPeer; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBuffer[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBuffer[offset+4 : offset+6]))
	}

	return peers, nil

}

func (t *TorrentFile) RequestPeers(peerID [20]byte, port uint16) ([]HandShake.PeerInfo, error) {
	url, err := t.buildTrackerURL(peerID, port)
	fmt.Println(url)
	if err != nil {
		return nil, err
	}
	cli := &http.Client{Timeout: 15 * time.Second}
	resp, err := cli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackResp := new(TrackerResp)
	err = bencode.Unmarshal(resp.Body, trackResp)
	fmt.Println(trackResp.Interval)
	if err != nil {
		return nil, err
	}

	return buildPeerInfo([]byte(trackResp.Peers))

}
