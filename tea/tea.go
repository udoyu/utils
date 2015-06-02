package tea

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Tea struct {
	key []uint32
}

func NewTea(key []uint32) *Tea {
	if len(key) != 4 {
		return nil
	}
	return &Tea{key[:]}
}

func NewTeaFromBytes(key []byte) *Tea {
	if len(key) > 16 {
		key = key[:16]
	}else if len(key) < 16 {
		buf := make([]byte, 16)
		copy(buf[:], key[:])
		key = buf
	}
	t := &Tea{}
	t.key = make([]uint32, 4)
	buf := bytes.NewBuffer(key)
	binary.Read(buf, binary.LittleEndian, t.key[:])
	return t
}

func (this *Tea) Encrypt(plain []byte)(crypt []byte, err error) {
	return TEAEncrypt(plain, this.key)
}

func (this *Tea) Decrypt(crypt []byte)(plain []byte, err error) {
	return TEADecrypt(crypt, this.key)
}

func (this *Tea) EncryptWrapper(plain []byte)(crypt []byte, err error) {
	return TEAEncryptWrapper(plain, this.key)
}

func (this *Tea) DecryptWrapper(crypt []byte)(plain []byte, err error) {
	return TEADecryptWrapper(crypt, this.key)
}

func TEADecrypt(crypt []byte, key []uint32) (plain []byte, err error) {
	crypt_len := len(crypt)
	plain_len := len(plain)

	if crypt_len < 1 || crypt_len%8 != 0 {
		err = errors.New(fmt.Sprintf("crypt_len=%d", crypt_len))
		return plain, err
	}
	//convert to uint32[]
	var tcrypt []uint32
	var tplain []uint32
	tcrypt = make([]uint32, crypt_len/4)
	buf := bytes.NewBuffer(crypt)
	binary.Read(buf, binary.LittleEndian, tcrypt[:])
	tplain = make([]uint32, crypt_len/4)

	length := crypt_len
	pre_plain := []uint32{0, 0}
	p_buf := make([]uint32, 2)
	c_buf := make([]uint32, 2)
	tkey := key
	tinyDecrypt(tcrypt, tkey, p_buf, 4)
	copy(pre_plain[:], p_buf)
	copy(tplain[:2], p_buf)
	for i := 2; i < length/4; i += 2 {
		c_buf[0] = tcrypt[i] ^ pre_plain[0]
		c_buf[1] = tcrypt[i+1] ^ pre_plain[1]
		tinyDecrypt(c_buf, tkey, p_buf, 4)
		copy(pre_plain, p_buf)
		tplain[i] = p_buf[0] ^ tcrypt[i-2]
		tplain[i+1] = p_buf[1] ^ tcrypt[i-1]
	}
	if uint32(tplain[length/4-1]) != 0 || uint32(tplain[length/4-2])&0xffffff00 != 0 {
		err = errors.New(fmt.Sprintf("length=%d|tplain[length/4-1]=%x|tplain[length/4-2]=%x", length, tplain[length/4-1], tplain[length/4-2]))
		return plain, err
	}
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, tplain)
	tplain_buf := buf.Bytes()

	padLentgh := int(tplain_buf[0] & 0x07)
	plain_len = crypt_len - (padLentgh + 3) - 7
	plain = make([]byte, 1024)
	copy(plain[:], tplain_buf[(padLentgh+3):(length+(padLentgh+3))])
	plain = plain[:plain_len]
	return plain, err
}

func TEAEncrypt(plain []byte, key []uint32) (crypt []byte, err error) {
	plain_len := len(plain)
	if plain_len < 1 {
		err = errors.New(fmt.Sprintf("plain_len=%d", plain_len))
		return crypt, err
	}
	pad := [9]byte{0xad, 0xad, 0xad, 0xad, 0xad, 0xad, 0xad, 0xad, 0xad}
	var tcrypt []uint32
	var tplain []uint32
	if len(plain)%4 != 0 {
		tmp := make([]byte, (len(plain)/4+1)*4)
		copy(tmp[:], plain)
		plain = tmp
	}
	tplain = make([]uint32, plain_len/4)
	buf := bytes.NewBuffer(plain)
	binary.Read(buf, binary.LittleEndian, tplain[:])
	pre_plain := []uint32{0, 0}
	pre_crypt := []uint32{0, 0}
	p_buf := make([]uint32, 2)
	c_buf := make([]uint32, 2)
	padLentgh := (plain_len + 10) % 8
	if padLentgh != 0 {
		padLentgh = 8 - padLentgh
	}
	length := padLentgh + 3 + plain_len + 7

	tcrypt_buf := make([]byte, length)
	tcrypt_buf[0] = 0xa8 | byte(padLentgh)
	copy(tcrypt_buf[1:], pad[:padLentgh+2])
	copy(tcrypt_buf[padLentgh+3:], plain)
	tcrypt = make([]uint32, len(tcrypt_buf)/4)
	buf = bytes.NewBuffer(tcrypt_buf)
	binary.Read(buf, binary.LittleEndian, tcrypt[:len(tcrypt_buf)/4])
	for i := 0; i < length/4; i += 2 {
		p_buf[0] = tcrypt[i] ^ pre_crypt[0]
		p_buf[1] = tcrypt[i+1] ^ pre_crypt[1]
		tinyEncrypt(p_buf, key, c_buf, 4)
		tcrypt[i] = c_buf[0] ^ pre_plain[0]
		tcrypt[i+1] = c_buf[1] ^ pre_plain[1]
		copy(pre_crypt[:], tcrypt[i:i+2])
		copy(pre_plain[:], p_buf[:2])
	}
	fmt.Println(tcrypt)
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, tcrypt[:])
	tcrypt_buf = buf.Bytes()

	crypt = tcrypt_buf[:length]

	return crypt, err
}

func TEADecryptWrapper(crypt []byte, key []uint32) (plain []byte, err error) {
	crypt_len := len(crypt)
	if crypt_len%2 != 0 {
		err = errors.New(fmt.Sprintf("crypt_len=%d", crypt_len))
		return plain, err
	}
	real_crypt_len := crypt_len / 2
	real_crypt := make([]byte, real_crypt_len)
	for i := 0; i < real_crypt_len; i++ {
		b1 := HexCharToInt(crypt[i*2])
		b2 := HexCharToInt(crypt[i*2+1])
		real_crypt[i] = byte(b2<<4 | b1)
	}
	return TEADecrypt(real_crypt, key)
}

func TEAEncryptWrapper(plain []byte, key []uint32) (crypt []byte, err error) {
	crypt, err = TEAEncrypt(plain, key)
	if err != nil {
		return crypt, err
	}
	crypt_len := len(crypt)
	real_crypt_len := 2 * len(crypt)
	real_crypt := make([]byte, real_crypt_len)
	for i := 0; i < crypt_len; i++ {
		b1 := crypt[i] & 0x0F
		b2 := (crypt[i] & 0xF0) >> 4
		real_crypt[i*2] = IntToHexChar(int(b1))
		real_crypt[i*2+1] = IntToHexChar(int(b2))
	}
	crypt = real_crypt[:]
	return crypt, err
}


func HexCharToInt(c byte) int {
	if c <= '9' && c >= '0' {
		return int(c - '0')
	} else if c <= 'f' && c >= 'a' {
		return int(c-'a') + 10
	} else if c <= 'F' && c >= 'A' {
		return int(c-'A') + 10
	}
	return 0
}

func IntToHexChar(i int) byte {
	return fmt.Sprintf("%x", i)[0]
}

var (
	TEA_DELTA uint32 = 0x9E3779B9
	TEA_SUM   uint32 = 0xE3779B90
)

func tinyEncrypt(plain, key, crypt []uint32, power uint32) {
	rounds := uint32(1 << power)
	sum := uint32(0)
	y := plain[0]
	z := plain[1]
	a := key[0]
	b := key[1]
	c := key[2]
	d := key[3]
	for i := uint32(0); i < rounds; i++ {
		sum += TEA_DELTA
		y += ((z << 4) + a) ^ (z + sum) ^ ((z >> 5) + b)
		z += ((y << 4) + c) ^ (y + sum) ^ ((y >> 5) + d)
	}
	crypt[0] = y
	crypt[1] = z
}

func tinyDecrypt(crypt, key, plain []uint32, power uint32) {
	rounds := uint32(1 << power)
	sum := uint32(TEA_DELTA << power)
	y := crypt[0]
	z := crypt[1]
	a := key[0]
	b := key[1]
	c := key[2]
	d := key[3]
	for i := uint32(0); i < rounds; i++ {
		z -= ((y << 4) + c) ^ (y + sum) ^ ((y >> 5) + d)
		y -= ((z << 4) + a) ^ (z + sum) ^ ((z >> 5) + b)
		sum -= TEA_DELTA
	}
	plain[0] = y
	plain[1] = z
}
