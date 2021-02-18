package buffer

import (
	"fmt"
	"os"
	"encoding/binary"
)

type Buffer struct {
	Path string
	fd *os.File
}

func (buf *Buffer) Load() error {
	fd, err := os.OpenFile(buf.Path, os.O_RDWR | os.O_CREATE, 0755)

	fmt.Println("Load result", fd, err);
	
	if err != nil {
		fmt.Print("Error", err)

		panic(err)
	}

	buf.fd = fd;

	return nil
}

func (buf *Buffer) read(numByte int) ([]byte, error) {
	tmp := make([]byte, numByte)
	
	_, err := buf.fd.Read(tmp)

	return tmp, err
}

func (buf *Buffer) ReadByte() (byte, error) {
	b, err := buf.read(1)

	if err != nil {
		return 0, err
	}

	return b[0], err
}

func (buf *Buffer) ReadShort() (int, error) {
	b, err := buf.read(2)

	if err != nil {
		return 0, err
	}

	num := binary.LittleEndian.Uint16(b)

	return int(num), err
}

func (buf *Buffer) ReadInt() (int, error) {
	b, err := buf.read(4)

	if err != nil {
		return 0, err
	}

	num := binary.LittleEndian.Uint32(b)

	return int(num), err
}

func (buf *Buffer) ReadLong() (int, error) {
	b, err := buf.read(8)

	if err != nil {
		return 0, err
	}

	num := binary.LittleEndian.Uint64(b)

	return int(num), err
}

func (buf *Buffer) ReadULeb128() (int, error) {
	total := 0
	shift := 0

	var b byte

	for b, _ = buf.ReadByte(); (int(b) & 0x80) != 0; b, _ = buf.ReadByte() {
		total |= (int(b) & 0x7F) << shift
		shift += 7
	}

	total |= (int(b) & 0x7F) << shift

	return total, nil
}

func (buf *Buffer) ReadString() (string, error) {
	b, err := buf.ReadByte()

	if err != nil {
		return "", err
	}

	if b == 0x00 {
		return "", err
	}

	len, err := buf.ReadULeb128()
	res, err := buf.read(len)

	return string(res), err
}

func (buf *Buffer) write(bytes []byte) error {
	_, err := buf.fd.Write(bytes)

	return err
}

func (buf *Buffer) WriteByte(b byte) error {
	bytes := []byte{b}

	return buf.write(bytes);
}

func (buf *Buffer) WriteShort(val int) error {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(val))

	return buf.write(bytes);
}

func (buf *Buffer) WriteInt(val int) error {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(val))

	return buf.write(bytes);
}

func (buf *Buffer) WriteLong(val int) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(val))

	return buf.write(bytes);
}

func (buf *Buffer) WriteULeb128(val int) error {
	bytes := []byte{}

	for val > 0 {
		b := val & 0x7F
		val = val >> 7
		if val != 0 {
			b = b | 0x80
		}

		bytes = append(bytes, byte(b))
	}

	return buf.write(bytes)
}

func (buf *Buffer) WriteString(val string) error {
	bytes := []byte(val)

	if len(bytes) == 0 {
		return buf.write([]byte{0})
	}

	buf.write([]byte{11})
	buf.WriteULeb128(len(bytes))
	buf.write(bytes)

	return nil
}