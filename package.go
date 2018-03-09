package help

import (
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
)

type Package struct {
	Buf []byte
}

func NewPackage(params ...[]byte) *Package {
	p := &Package{[]byte{}}
	for _, data := range params {
		p.AddByteArray(data)
	}
	return p
}

//新桥牌不再使用GLID_XXX的方式，使用字符串代替，例如："chat"
//长度为：1-8个字符（可用字符为数字、大小写字母、下划线），不足8个的会在后面补空格
func (this *Package) AddCmd(t string) {
	matched, err := regexp.MatchString("^[0-9a-zA-Z_]{1,8}$", t)
	if !matched || err != nil {
		log.Panicf("类型必须为1-8字节字符串，可用字符为数字、大小写字母、下划线,cmd = [%s]", t)
	}
	t = t + "       "
	hex := []byte(t)
	var n = binary.LittleEndian.Uint64(hex)
	this.AddUint64(n)
}

func (this *Package) AddUint8(n uint8) {
	this.Buf = append(this.Buf, n)
}

func (this *Package) AddUint16(n uint16) {
	buf := make([]byte, 2)
	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) Reset() {
	this.Buf = nil
}

//对于新桥牌内部协议，设置消息类型应使用AddType
func (this *Package) AddUint32(n uint32) {
	buf := make([]byte, 4)
	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	buf[2] = uint8(n >> 16)
	buf[3] = uint8(n >> 24)
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) AddUint64(n uint64) {
	buf := make([]byte, 8)
	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	buf[2] = uint8(n >> 16)
	buf[3] = uint8(n >> 24)
	buf[4] = uint8(n >> 32)
	buf[5] = uint8(n >> 40)
	buf[6] = uint8(n >> 48)
	buf[7] = uint8(n >> 56)
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) AddInt16(n int16) {
	buf := make([]byte, 2)
	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) AddInt32(n int32) {
	buf := make([]byte, 4)
	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	buf[2] = uint8(n >> 16)
	buf[3] = uint8(n >> 24)
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) AddBool(isTrue bool) {
	if isTrue {
		this.AddInt32(1)
	} else {
		this.AddInt32(0)
	}
}

func (this *Package) AddStr(str string) {
	buf := []byte(str)
	this.Buf = append(this.Buf, buf...)
	this.Buf = append(this.Buf, 0)
}

func (this *Package) AddString(str string, Length int) {
	buf := []byte(str)
	if len(buf) > Length {
		buf = buf[:Length]
	} else if len(buf) < Length {
		zeroCount := Length - len(buf)
		for i := 0; i < zeroCount; i++ {
			buf = append(buf, 0)
		}
	}
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) ReadCmd() string {
	if len(this.Buf) < 8 {
		log.Panicf("期望读取8字节，实际只有%d", len(this.Buf))
	}
	var cmd = string(this.Buf[0:8])
	for i := 0; i < 7; i++ {
		if cmd[len(cmd)-1] == 0 || cmd[len(cmd)-1] == ' ' {
			cmd = cmd[0 : len(cmd)-1]
		}
	}
	this.Buf = this.Buf[8:]
	return cmd
}

func (this *Package) ReadStr() string {
	for i := 0; i < len(this.Buf); i++ {
		if this.Buf[i] == 0 {
			str := string(this.Buf[:i])
			this.Buf = this.Buf[i+1:]
			return str
		}
	}
	return ""
}

func (this *Package) ReadString(Length int) string {
	buf := this.Buf[:Length]
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			buf = buf[:i]
			break
		}
	}
	this.Buf = this.Buf[Length:]
	return string(buf)
}

func (this *Package) ReadIp() string {
	ip := fmt.Sprintf("%d.%d.%d.%d", this.Buf[0], this.Buf[1], this.Buf[2], this.Buf[3])
	this.Buf = this.Buf[4:]
	return ip
}

func (this *Package) ReadUint8() uint8 {
	n := this.Buf[0]
	this.Buf = this.Buf[1:]
	return n
}

func (this *Package) ReadUint16() uint16 {
	var n uint16
	n = uint16(this.Buf[0])
	n |= uint16(this.Buf[1]) << 8
	this.Buf = this.Buf[2:]
	return n
}

func (this *Package) ReadUint32() uint32 {
	var n uint32
	n = uint32(this.Buf[0])
	n |= uint32(this.Buf[1]) << 8
	n |= uint32(this.Buf[2]) << 16
	n |= uint32(this.Buf[3]) << 24
	this.Buf = this.Buf[4:]
	return n
}

func (this *Package) ReadInt16() int16 {
	var n int16
	n = int16(this.Buf[0])
	n |= int16(this.Buf[1]) << 8
	this.Buf = this.Buf[2:]
	return n
}

func (this *Package) ReadInt32() int32 {
	var n int32
	n = int32(this.Buf[0])
	n |= int32(this.Buf[1]) << 8
	n |= int32(this.Buf[2]) << 16
	n |= int32(this.Buf[3]) << 24
	this.Buf = this.Buf[4:]
	return n
}

func (this *Package) ReadBool() bool {
	if this.ReadInt32() == 0 {
		return false
	}
	return true
}

func (this *Package) ReadUint64() uint64 {
	var n uint64
	n = uint64(this.Buf[0])
	n |= uint64(this.Buf[1]) << 8
	n |= uint64(this.Buf[2]) << 16
	n |= uint64(this.Buf[3]) << 24
	n |= uint64(this.Buf[4]) << 32
	n |= uint64(this.Buf[5]) << 40
	n |= uint64(this.Buf[6]) << 48
	n |= uint64(this.Buf[7]) << 56
	this.Buf = this.Buf[8:]
	return n
}

func (this *Package) ReadInt64() int64 {
	var n int64
	n = int64(this.Buf[0])
	n |= int64(this.Buf[1]) << 8
	n |= int64(this.Buf[2]) << 16
	n |= int64(this.Buf[3]) << 24
	n |= int64(this.Buf[4]) << 32
	n |= int64(this.Buf[5]) << 40
	n |= int64(this.Buf[6]) << 48
	n |= int64(this.Buf[7]) << 56
	this.Buf = this.Buf[8:]
	return n
}

func (this *Package) AddByteArray(buf []byte) {
	this.Buf = append(this.Buf, buf...)
}

func (this *Package) AddPackage(pkg *Package) {
	this.Buf = append(this.Buf, pkg.Buf...)
}

func (this *Package) GetBuffer() []byte {
	return this.Buf
}

func (this *Package) SetBuffer(buf []byte) {
	this.Buf = buf
}
