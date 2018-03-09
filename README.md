# help
Parse binary stream and timer

import "help"

func main() {
  
  	timer := help.NewTimer()
  	
	timer.SetTimer(func() bool {
  	
		fmt.Println("timer")
    
 	}, time.second * 10, 0)//0:keep live, >0:timer execution many times
  

	timer.Killed()

  	pkg := help.NewPackage()
	pkg.AddByteArray([]byte{0x12, 0x34, 0x56, 0x78, 0x99})

	log.Printf("0x%x-0x%x\n", pkg.ReadUint8(), pkg.ReadUint32())//0x12-0x99785634

	pkg.AddUint8(0x12)
	pkg.AddUint32(0x99785634)
	log.Printf("%#v", pkg.GetBuffer())//[]byte{0x12, 0x34, 0x56, 0x78, 0x99}
  
}
