package Bencode

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

func checkNum(data byte) bool {
	return data >= '0' && data <= '9'
}

func WriteDecimal(bw *bufio.Writer, val int) (len int) {
	if val < 0 {
		_ = bw.WriteByte('-')
		len++
	}
	val = int(math.Abs(float64(val)))
	valStr := fmt.Sprintf("%d", val)
	_, _ = bw.Write([]byte(valStr))
	for val > 0 {
		val /= 10
		len++
	}

	return len
}

func ReadDecimal(br *bufio.Reader) (val int, len int) {
	flag := 1
	b, _ := br.ReadByte()
	len++
	if b == '-' {
		flag = -1
		len++
	} else {
		val = val*10 + int(b-'0')
	}

	for {
		b, _ = br.ReadByte()
		if !checkNum(b) {
			return flag * val, len
		}
		val = val*10 + int(b-'0')
		len++
	}

	//return -1, -1
}

func EncodeInt(w io.Writer, val int) int {
	bw := bufio.NewWriter(w)
	WritingLen := 0
	//先写入i
	_ = bw.WriteByte('i')
	WritingLen++
	//再写入数字
	nLen := WriteDecimal(bw, val)
	WritingLen += nLen
	//写入e
	_ = bw.WriteByte('e')
	WritingLen++

	err := bw.Flush() //将buff中的数据写入到io.Writer中
	if err != nil {
		return 0
	}
	return WritingLen
}

func DecodeInt(r io.Reader) (val int, err error) { //用于将Stream中的文本流转换为Bobject
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	b, err := br.ReadByte()
	if b != 'i' {
		return val, err
	}

	val, LEN := ReadDecimal(br)
	if LEN == -1 {
		return val, err
	}
	return
}

func EncodeString(w io.Writer, val string) int {
	strLen := len(val)
	bw := bufio.NewWriter(w)
	wLen := WriteDecimal(bw, strLen)
	_ = bw.WriteByte(':')
	wLen++
	_, _ = bw.WriteString(val)
	wLen += strLen

	err := bw.Flush()
	if err != nil {
		return 0
	}

	return wLen

}

func DecodeString(r io.Reader) (val string, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	num, LEN := 0, 0
	var b byte
	for {
		b, _ = br.ReadByte()
		if !checkNum(b) {
			err = br.UnreadByte()
			if err != nil {
				return val, err
			}
			break
		}
		num = num*10 + int(b-'0')
		LEN++
	}
	if LEN == 0 {
		return val, err
	}
	b, err = br.ReadByte()
	if b != ':' {
		return val, err
	}
	buf := make([]byte, num)
	_, err = io.ReadAtLeast(br, buf, num)
	if err != nil {
		return val, err
	}
	val = string(buf)
	return
}
