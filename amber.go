package main


import "gopkg.in/cheggaaa/pb.v1"
import "github.com/fatih/color"
import "strconv"
import "os/exec"
import "strings"
import "runtime"
import "fmt"
import "os"

const VERSION string = "1.0.0"


type PARAMETERS struct {
  fileName string 
  keySize int
  key []byte
  staged bool
  resource bool
  verbose bool
}

type PE struct {
  imageBase string
  subSystem string
}

var Red *color.Color = color.New(color.FgRed)
var BoldRed *color.Color = Red.Add(color.Bold)
var	Blue *color.Color = color.New(color.FgBlue)
var	BoldBlue *color.Color = Blue.Add(color.Bold)
var	Yellow *color.Color = color.New(color.FgYellow)
var	BoldYellow *color.Color = Yellow.Add(color.Bold)
var	Green *color.Color = color.New(color.FgGreen)
var	BoldGreen *color.Color = Green.Add(color.Bold)


var progressBar *pb.ProgressBar
var parameters PARAMETERS
var pe PE

func main() {

  runtime.GOMAXPROCS(runtime.NumCPU())

  parameters.keySize = 8
  parameters.staged = false
  parameters.resource = true
  parameters.verbose = false

 	ARGS := os.Args[1:]
  if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h"{
    Banner()
    Help()
    os.Exit(0)
  }

  Banner()
  parameters.fileName = ARGS[0]

  for i := 0; i < len(ARGS); i++{
  	if ARGS[i] == "-ks" || ARGS[i] == "--keysize" {
  		ks, Err := strconv.Atoi(ARGS[i+1])
      if Err != nil {
        BoldRed.Println("\n[!] ERROR: Invalid key size.\n")
        fmt.Println(Err)
        os.Exit(1)
      }else{
        parameters.keySize = ks
      } 
  	}
  	if ARGS[i] == "-k" || ARGS[i] == "--key" {
  		parameters.key = []byte(ARGS[i+1]) 
  	}
  	if ARGS[i] == "--staged" {
  		parameters.staged = true 
  	}
  	if ARGS[i] == "--no-resource" {
  		parameters.resource = false 
  	}
    if ARGS[i] == "-v" || ARGS[i] == "--verbose" {
      parameters.verbose = true 
    }
  }

  BoldYellow.Print("\n[*] File: ")
  BoldBlue.Println(parameters.fileName)
  BoldYellow.Print("[*] Staged: ")
  BoldBlue.Println(parameters.staged)
  if len(parameters.key) != 0 {
    BoldYellow.Print("[*] Key: ")
    BoldBlue.Println(parameters.key)
  }else{
    BoldYellow.Print("[*] Key Size: ")
    BoldBlue.Println(parameters.keySize)   
  }
  BoldYellow.Print("[*] Verbose: ")
  BoldBlue.Println(parameters.verbose,"\n")


  if parameters.staged == true {
  	progressBar = pb.New(21)
  }else{
  	progressBar = pb.New(27)
  }
  progressBar.SetWidth(80)
  progressBar.Start()  


  if CheckRequirments() == false {
    BoldRed.Println("\n\n[!] ERROR: Amber is not installed properly (missing dependencies)")
    os.Exit(1)
  }

  progressBar.Increment()
  InspectPE()
  BuildPayload()
  if parameters.staged == true {
  	exec.Command("sh", "-c", string("mv Payload "+parameters.fileName+".stage")).Run()
  }else{
  	CryptPayload()
  	CompileStub()
  }
  CleanFiles()
  progressBar.FinishPrint("\n")

  if parameters.staged == true {
  	BoldGreen.Println("[+] Stage generated !\n")	
  }else{
  	BoldGreen.Println("[+] File successfully crypted !\n")	
  }
  

}


func InspectPE() {

  if parameters.verbose == true {
    BoldYellow.Println("[*] Striping pe file... ")
  }

  exec.Command("sh", "-c", string("strip "+parameters.fileName)).Run()
  progressBar.Increment()

  ls, Err := exec.Command("sh", "-c", string("ls  "+parameters.fileName)).Output()
  if (!strings.Contains(string(ls), parameters.fileName)) || (Err != nil)  {
    BoldRed.Println("\n[!] ERROR: Unable to locate file :(")
    Red.Println(string(ls))
    Red.Println(Err)
    os.Exit(1)
  }

  progressBar.Increment()

	magic, _ := exec.Command("sh", "-c", string("objdump -x "+parameters.fileName+"|grep Magic|tr -d \"\\n\"")).Output()
	if !strings.Contains(string(magic), "010b") {
		BoldRed.Println("\n[!] ERROR: File is not a valid PE")
		BoldRed.Println(string(magic))
		os.Exit(1)
	}
  progressBar.Increment()
	arch, _ := exec.Command("sh", "-c", string("objdump -x "+parameters.fileName+"|grep architecture|tr -d \"\\n\"")).Output()
	if !strings.Contains(string(arch), "i386"){
		BoldRed.Println("\n[!] ERROR: Unsupported file architecture (only 32 PE files supported)")
		BoldYellow.Println(string(arch))
		os.Exit(1)		
	}
  progressBar.Increment()
	imageBase, _ := exec.Command("sh", "-c", string("objdump -x "+parameters.fileName+"| grep ImageBase|tr -d \"ImageBase		\"|tr -d \"\\n\" ")).Output()
  pe.imageBase = string(imageBase)
  progressBar.Increment()
	subSystem, _ := exec.Command("sh", "-c", string("objdump -x "+parameters.fileName)).Output()
	if strings.Contains(string(subSystem), "(Windows GUI)") {
		pe.subSystem = "(Windows GUI)"
	}else{
		pe.subSystem = "(Windows CUI)"
	}
	progressBar.Increment()
	boundImports, _ := exec.Command("sh", "-c", string("objdump -x "+parameters.fileName+"|grep \"Bound Import Directory\" |tr -d \"Entry b \"|tr -d \"BoudImpoDieco\"")).Output()
	if string(boundImports) != "0000000000000000\n" {
		BoldRed.Println("\n[!] ERROR: Incompatible PE file (file has bounded imports)")
		BoldYellow.Println(string(boundImports))
		os.Exit(1)
	}
  progressBar.Increment()
  if parameters.verbose == true {
    BoldYellow.Println("\n[*] "+string(magic))
    BoldYellow.Println("[*] "+string(arch))
    BoldYellow.Println("[*] ImageBase: 0x"+pe.imageBase)
    BoldYellow.Println("[*] SubSystem: "+pe.subSystem)
  }
}


func BuildPayload() {

  MapPE, _ := exec.Command("sh", "-c", string("wine MapPE.exe "+parameters.fileName)).Output()
  if strings.Contains(string(MapPE), "[!]") {
    BoldRed.Println("\n[!] ERROR: While mapping pe file :(")
    BoldRed.Println(string(MapPE))
    CleanFiles()
    os.Exit(1)    	
  }
  progressBar.Increment()

  moveMap, moveMapErr := exec.Command("sh", "-c", "mv Mem.map ReplaceProcess/").Output()
  if moveMapErr != nil {
    BoldRed.Println("\n[!] ERROR: While moving the file map")
    BoldRed.Println(string(moveMap))
    CleanFiles()
    os.Exit(1)    	
  }

  progressBar.Increment()
  nasm, Err := exec.Command("sh", "-c", "cd ReplaceProcess && nasm -f bin ReplaceProcess.asm -o Payload").Output()
  if Err != nil {
    BoldRed.Println("\n[!] ERROR: While assembling payload :(")
    BoldRed.Println(string(nasm))
    BoldRed.Println(Err)
    CleanFiles()
    os.Exit(1)    
  }

  movePayload, movePayErr := exec.Command("sh", "-c", "mv ReplaceProcess/Payload ./").Output()
  if movePayErr != nil {
    BoldRed.Println("\n[!] ERROR: While moving the payload")
    BoldRed.Println(string(movePayload))
    BoldRed.Println(Err)
    CleanFiles()
    os.Exit(1)    
  }


  progressBar.Increment()

  if parameters.verbose == true {
    _MapPE := strings.Split(string(MapPE), "github.com/egebalci/mappe")
    fmt.Println(string(_MapPE[1]))
  }

}


func CryptPayload() {
  if parameters.verbose == true {
    BoldYellow.Println("[*] Ciphering payload...")    
  }
  if len(parameters.key) != 0 {
    exec.Command("sh", "-c", string("./bitbender ^=\""+string(parameters.key)+"\" Payload")).Run()  
  }else{
    ks := strconv.Itoa(parameters.keySize)
    exec.Command("sh", "-c", "./bitbender ^ "+ks+" Payload").Run()
  }
  progressBar.Increment()  

  xxd := exec.Command("sh", "-c", "rm Payload && mv Payload.xor Payload && xxd -i Payload > Stub/PAYLOAD.h")
  xxd.Stdout = os.Stdout
  xxd.Stderr = os.Stderr
  xxd.Run()

  progressBar.Increment()  

  _xxd := exec.Command("sh", "-c", "xxd -i Payload.key > Stub/KEY.h")
  _xxd.Stdout = os.Stdout
  _xxd.Stderr = os.Stderr
  _xxd.Run()

  progressBar.Increment()  

  if parameters.verbose == true {
    key, _ := exec.Command("sh", "-c", "xxd -i Payload.key").Output() 
    BoldYellow.Println("[*] Payload ciphered with: ")
    BoldBlue.Println(string(key))    
  }
  
}

func CompileStub() {

  if parameters.verbose == true {
    BoldYellow.Println("[*] Compiling Stub... ")
  }


  progressBar.Increment()
  var CompileCommand string = ""

  if parameters.resource == true {
  	mingwObj, Err := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 -c Stub/Stub.cpp").Output()
  	if Err != nil {
    	BoldRed.Println("\n[!] ERROR: While compiling the stub object :(")
    	Red.Println(string(mingwObj))
    	Red.Println(Err)
    	CleanFiles()
    	os.Exit(1)
  	}

  	CompileCommand += string("i686-w64-mingw32-g++-win32 Stub.o Stub/Resource.o -Wl,--image-base=0x"+pe.imageBase)
  }else{
  	CompileCommand += string("i686-w64-mingw32-g++-win32 Stub.cpp -Wl,--image-base=0x"+pe.imageBase)
  }


  if pe.subSystem == "(Windows GUI)"{
    CompileCommand += string("-mwindows -o "+parameters.fileName)  
  }else{
  	CompileCommand += string(" -o "+parameters.fileName)
  }

  mingw, Err2 := exec.Command("sh", "-c", CompileCommand).Output()
  if Err2 != nil {
    BoldRed.Println("\n[!] ERROR: While compiling the stub :(")
    Red.Println(string(mingw))
    Red.Println(Err2)
    CleanFiles()
    os.Exit(1)
  }
  progressBar.Increment()

  if parameters.verbose == true {
    BoldYellow.Println("[*] "+CompileCommand)
    BoldYellow.Println("[*] Striping crypted file... ")
  }

  exec.Command("sh", "-c", string("strip "+parameters.fileName)).Run()
  progressBar.Increment()
}


func CleanFiles() {

  exec.Command("sh", "-c", "rm ReplaceProcess/Mem.map").Run()
  exec.Command("sh", "-c", "rm Stub.o").Run()
  exec.Command("sh", "-c", "rm Payload").Run()
  exec.Command("sh", "-c", "rm Payload.xor").Run()
  exec.Command("sh", "-c", "rm Payload.key").Run()


  exec.Command("sh", "-c", "echo   > Stub/PAYLOAD.h").Run()
  exec.Command("sh", "-c", "echo   > Stub/KEY.h").Run()

  progressBar.Increment() 
}

func CheckRequirments() (bool){

  CheckMingw, _ := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 --version").Output()
  if (!strings.Contains(string(CheckMingw), "Copyright")) {
    return false
  }
  progressBar.Increment()
  CheckMingwDress, _ := exec.Command("sh", "-c", "i686-w64-mingw32-windres -V").Output()
  if (!strings.Contains(string(CheckMingwDress), "Copyright")) {
    return false
  }
  progressBar.Increment()
 	CheckNasm, _ := exec.Command("sh", "-c", "nasm -h").Output()
  if (!strings.Contains(string(CheckNasm), "usage:")) {
    return false
  }
  progressBar.Increment()
  CheckStrip, _ := exec.Command("sh", "-c", "strip -V").Output()
  if (!strings.Contains(string(CheckStrip), "Copyright")) {
    return false
  }
  progressBar.Increment()
  CheckWine, _ := exec.Command("sh", "-c", "wine --help").Output()
  if (!strings.Contains(string(CheckWine), "Usage:")) {
    return false
  }
  progressBar.Increment()
  Checkbitbender, _ := exec.Command("sh", "-c", "./bitbender").Output()
  if (!strings.Contains(string(Checkbitbender), "USAGE:")) {
    return false
  }
  progressBar.Increment()
  CheckMapPE, _ := exec.Command("sh", "-c", "ls MapPE.exe").Output()
  if (!strings.Contains(string(CheckMapPE), "MapPE.exe")) {
    return false
  }
  progressBar.Increment()
 	CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
  if (!strings.Contains(string(CheckXXD), "Amber")) {
    return false
  }
  progressBar.Increment()
  CheckMultiLib, _ := exec.Command("sh", "-c", "apt-cache policy gcc-multilib").Output()
  if (strings.Contains(string(CheckMultiLib), "(none)")) {
    return false
  }
  progressBar.Increment()
	CheckMultiLibPlus, _ := exec.Command("sh", "-c", "apt-cache policy g++-multilib").Output()
  if (strings.Contains(string(CheckMultiLibPlus), "(none)")) {
    return false
  }
  progressBar.Increment()
	return true

}


func Banner() {

  	var BANNER string = `

//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  POC Crypter For ReplaceProcess                                             
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
  amber file.exe [options]


OPTIONS:
  
  -k, --key       [string]        Custom cipher key
  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:100/Min:4)
  --staged                        Generated a staged payload
  --no-resource                   Don't add any resource
  -v, --verbose                   Verbose output mode
  -h, --help                      Show this massage

EXAMPLE:
  (Default settings if no option parameter passed)
  amber file.exe -ks 8 -o crypted.exe
`
  color.Green(Help)

}