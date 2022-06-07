package Torrent

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBuildTrackerURL(t *testing.T) {
	to := TorrentFile{
		Announce: "http://bttracker.debian.org:6969/announce",
		InfoHash: [20]byte{216, 247, 57, 206, 195, 40, 149, 108, 204, 91, 191, 31, 134, 217, 253, 207, 219, 168, 206, 182},
		PieceHashes: [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106},
			{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
		},
		PieceLength: 262144,
		FileLength:  351272960,
		FileName:    "debian-10.2.0-amd64-netinst.iso",
	}
	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	const port uint16 = 6881
	url, err := to.buildTrackerURL(peerID, port)
	expected := "http://bttracker.debian.org:6969/announce?compact=1&downloaded=0&info_hash=%D8%F79%CE%C3%28%95l%CC%5B%BF%1F%86%D9%FD%CF%DB%A8%CE%B6&left=351272960&peer_id=%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14&port=6881&uploaded=0"
	assert.Nil(t, err)
	assert.Equal(t, url, expected)
}

func TestRequestPeers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []byte(
			"d" +
				"8:interval" + "i900e" +
				"5:peers" + "12:" +
				string([]byte{
					192, 0, 2, 123, 0x1A, 0xE1, // 0x1AE1 = 6881
					127, 0, 0, 1, 0x1A, 0xE9, // 0x1AE9 = 6889
				}) + "e")
		w.Write(response)
	}))
	defer ts.Close()
	tf := TorrentFile{
		Announce: ts.URL,
		InfoHash: [20]byte{216, 247, 57, 206, 195, 40, 149, 108, 204, 91, 191, 31, 134, 217, 253, 207, 219, 168, 206, 182},
		PieceHashes: [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106},
			{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
		},
		PieceLength: 262144,
		FileLength:  351272960,
		FileName:    "debian-10.2.0-amd64-netinst.iso",
	}
	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	const port uint16 = 6882
	expected := []PeerInfo{
		{IP: net.IP{192, 0, 2, 123}, Port: 6881},
		{IP: net.IP{127, 0, 0, 1}, Port: 6889},
	}
	p, err := tf.RequestPeers(peerID, port)
	assert.Nil(t, err)
	assert.Equal(t, expected, p)

	for i, v := range expected {
		fmt.Printf("Peer:%d  ip:%s  port:%d\n", i, v.IP, v.Port)
	}

}

func TestIfGetRight(t *testing.T) {
	file, _ := os.Open("D:\\GolandProject\\go-torrent\\src\\debian-iso.torrent")
	tf, _ := ParseTorrentFile(bufio.NewReader(file))
	//tf.Announce = "http://207.241.226.111:6969/announce"

	var peerId [20]byte
	_, _ = rand.Read(peerId[:])

	peers, _ := tf.RequestPeers(peerId, 6881)
	for i, p := range peers {
		fmt.Printf("Peer %d, Ip: %s, Port: %d\n", i, p.IP, p.Port)
	}
}
