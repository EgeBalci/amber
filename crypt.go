package main

import (
	"crypto/rc4"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func crypt() {

	verbose("Ciphering payload...", "*")

	key := GenerateKey(target.KeySize)
	if target.KeySize != len(key) {
		verbose(string("Key size rounded to "+strconv.Itoa(len(key))), "*")
	}
	payload, err := ioutil.ReadFile(target.workdir + "/payload")
	parseErr(err)
	payload = RC4(payload, key)
	encPayload, err := os.Create("payload.enc")
	parseErr(err)
	payloadKey, err := os.Create("payload.key")
	parseErr(err)
	encPayload.Write(payload)
	payloadKey.Write(key)
	defer encPayload.Close()
	defer payloadKey.Close()
	verbose("Payload encrypted with RC4 algorithm", "*")
	verbose("Key \n", "*")
	verbose(xxd("payload.key"), "B")
	defer progress()
}

func xxd(fileName string) string {
	file, err := ioutil.ReadFile(fileName)
	parseErr(err)
	out := "{"
	for i, j := range file {
		out += fmt.Sprintf("0x%02X", j)
		if i != len(file)-1 {
			out += ", "
		}
		if i+1%12 == 0 {
			out += "\n"
		}
	}
	out += "}\n"
	defer progress()
	return out
}

func xor(Data []byte, Key []byte) []byte {
	for i := 0; i < len(Data); i++ {
		Data[i] = (Data[i] ^ (Key[(i % len(Key))]))
	}
	defer progress()
	return Data
}

// RC4 .
func RC4(data []byte, key []byte) []byte {
	c, e := rc4.NewCipher(key)
	parseErr(e)
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	defer progress()
	return dst
}

// GenerateKey .
func GenerateKey(Size int) []byte {

	if target.reflective == true && (Size%8) != 0 && Size >= 8 {
		Size += (8 - (Size % 8))
	}

	Key := make([]byte, Size)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < Size; i++ {
		Key[i] = byte(rand.Intn(255))
	}
	defer progress()
	return Key
}

// RandomString .
func RandomString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}
	return string(b)
}
