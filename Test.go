package XRaft

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"

	"log"
	"io/ioutil"
)

func main(){
	time.Sleep(20)
	fmt.Println("end sleep")
	x := make([]int ,5)
	fmt.Println(len(x))

	kv := make(map[string] string)

	kv["abc"] = "123"
	fmt.Println(kv)
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(kv)
	if err != nil{
		log.Fatal("error.")
	}
	ioutil.WriteFile("testfile",buffer.Bytes(), 0666)

	kv = make(map[string] string)

	raw, err := ioutil.ReadFile("testfile")
	if err != nil{
		panic(err)
	}
	buffer = bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&kv)
	if err != nil{
		panic(err)
	}
	fmt.Println(kv)


}
