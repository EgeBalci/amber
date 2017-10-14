package main


import "math/rand"
import "io/ioutil"
import "os/exec"
import "time"
import "os"

func crypt() {

  verbose("Ciphering payload...","*")

  if len(peid.key) != 0 {
    payload, err := ioutil.ReadFile("Payload")
    ParseError(err,"[!] ERROR: Can't open payload file.","")

    progress()
    payload = xor(payload,peid.key)
    payload_xor, err2 := os.Create("Payload.xor")
    ParseError(err2,"[!] ERROR: Can't create payload.xor file.","")

    progress()
    payload_key, err3 := os.Create("Payload.key")
    ParseError(err3,"[!] ERROR: Can't create payload.xor file.","")
    payload_xor.Write(payload)
    payload_xor.Write(peid.key)

    payload_xor.Close()
    payload_key.Close()
    progress()
  }else{
    key := generateKey(peid.keySize)
    progress()
    payload, err := ioutil.ReadFile("Payload")
    ParseError(err,"[!] ERROR: Can't open payload file.","")
    progress()
    payload = xor(payload,key)
    payload_xor, err2 := os.Create("Payload.xor")
    ParseError(err2,"[!] ERROR: Can't create payload.xor file.","")
    progress()
    payload_key, err3 := os.Create("Payload.key")
    ParseError(err3,"[!] ERROR: Can't create payload.xor file.","")
    payload_xor.Write(payload)
    payload_key.Write(key)

    payload_xor.Close()
    payload_key.Close()
  }
  progress()

  hex, _ := exec.Command("sh", "-c", "xxd -i Payload.key").Output()
  verbose("Payload ciphered with: \n","*")
  verbose(string(hex),"B")
}


func xor(Data []byte, Key []byte) ([]byte){
  for i := 0; i < len(Data); i++{
    Data[i] = (Data[i] ^ (Key[(i%len(Key))]))
  }
  return Data
}


func generateKey(Size int) ([]byte){
  Key := make([]byte, Size)
  rand.Seed(time.Now().UTC().UnixNano())
  for i := 0; i < Size; i++{
    Key[i] = byte(rand.Intn(255))
  }
  return Key
}

// Implement RC4...