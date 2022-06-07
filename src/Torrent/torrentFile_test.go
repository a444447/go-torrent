package Torrent

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestParseTorrentFile(t *testing.T) {
	file, err := os.Open("D:\\GolandProject\\go-torrent\\src\\Watch_dogs.torrent")
	if err != nil {
		t.Fatal("open error")
	}
	tf, err := ParseTorrentFile(bufio.NewReader(file))
	if err != nil {
		t.Fatal("Parse error")
	}
	if tf.Announce != "" {
		t.Log("right tracker")
	}
	fmt.Println(tf.Announce)
	fmt.Println(tf.FileName)

}
