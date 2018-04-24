package main

import (
	_ "os"
	"fmt"
	"flag"
	 _ "io"
	_ "io/ioutil"
	"bufio"
	"os"
	"log"
	"io"
	"strings"
	"strconv"
	pb "XRaft/xraftpb"
	"google.golang.org/grpc"
	"XRaft"
	"net"
	"google.golang.org/grpc/reflection"
	"time"


	"sync"
	"encoding/gob"
	"bytes"

)

//读取配置信息 以及 节点 id
type ServInfo struct{
	Ip string
	Port int
}
type Config struct{
	Servs [] ServInfo
	Id int
}

type KVServer struct{
	mutux *sync.Mutex
	datafile *os.File
	db map[string]string
}


func MakeKVServer(path string) *KVServer{
	kv := &KVServer{}
	kv.db = make(map[string] string)
	var err error
	kv.datafile,err = os.OpenFile(path,os.O_CREATE|os.O_RDWR,0600)
	if err !=nil{
		log.Fatal("[Error]: data file path do not exsist.")
	}
	return kv
}
func (kv *KVServer) Get(key string) (string,error){
	val, ok := kv.db[key]
	if ok {
		return val, nil
	}
	return val, fmt.Errorf(" key %q do not exsist.",key)

}

func (kv *KVServer) Put(key, val string){
	kv.mutux.Lock()
	defer kv.mutux.Unlock()
	kv.db[key] =val
}

func (kv *KVServer) Append(key, val string) {
	kv.mutux.Lock()
	kv.mutux.Unlock()
	kv.db[key] += val
}

func (kv *KVServer) Delete(key string) error{
	kv.mutux.Lock()
	kv.mutux.Unlock()
	_, err := kv.Get(key)
	if err == nil{
		delete(kv.db ,key)
	}

	return err
}


func (kv *KVServer) Save(){
	buff := bufio.NewWriter(kv.datafile)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(kv.db)
	if err != nil{
		log.Fatal("Encode Error while save data.")
	}

}


func (kv *KVServer) Load(){
	buff := bufio.NewReader(kv.datafile)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(&kv.db)
	if err != nil{
		log.Fatal("Decode Error while load data.")
	}
}

func (cfg * Config) Parse(name string) Config{


	file ,err := os.Open(name); if err != nil{
		log.Fatal("open config file error.")
	}else{
		//ioutil.rea
		bf := bufio.NewReader(file)
		for {
			addr, line_err := bf.ReadString('\n');
			if line_err != nil && line_err != io.EOF {
				log.Fatal("error while read config file.")
			} else if line_err == io.EOF{
				break
			} else {
				//fmt.Println(addr)
				//sv := new ServInfo(addr, 0)
				//fmt.Println(addr[:strings.Index(addr,":")])
				//fmt.Println(strconv.Atoi(strings.Replace( addr[strings.Index(addr,":")+1:],"\n","",-1)))
				idx := strings.Index(addr,":")

				port, _ := strconv.Atoi(strings.Replace( addr[idx+1:],"\n","",-1))
				cfg.Servs = append(cfg.Servs, ServInfo{addr[:idx], port})
			}
		}

	}


	id := flag.Int("id", 0, "Id of kv server in config file.")
	flag.Parse()
	cfg.Id = *id

	return *cfg
}


func InitClients(cfg Config) (peers [] *pb.XRaftClient){
	size := len(cfg.Servs)
	peers = make([] *pb.XRaftClient, size)

	for i := 0; i < size; i++{
		conn, err := grpc.Dial( cfg.Servs[i].Ip +":"+ strconv.Itoa(cfg.Servs[i].Port), grpc.WithInsecure())
		if err != nil{
			log.Fatal(err)
		}
		client := pb.NewXRaftClient(conn)
		peers[i] = &client
	}

	return peers
}
func main(){

	cfgx := (&Config{}).Parse("..//config//config")
	fmt.Print(cfgx)
	//fmt.Println(cfgx.Id)
	for i := range cfgx.Servs{
		fmt.Println(cfgx.Servs[i].Ip)
	}

	peers := InitClients(cfgx)
	fmt.Println(peers)
	ch := make(chan XRaft.ApplyMsg,1)

	//创建日志持久化对象
	persister := XRaft.MakePersister("..//logs")
	
	rft := XRaft.Make(peers, int32(cfgx.Id), persister, ch)

	fmt.Println(len(cfgx.Servs))
	lis, err := net.Listen("tcp",":"+strconv.Itoa( cfgx.Servs[cfgx.Id].Port))

	if err != nil{
		log.Fatal(err)
	}

	//打印当前节点的状态
	go func(){

		for{
			time.Sleep(time.Second)
			fmt.Println(rft.GetState())
		}
	}()



	//listen to the clients
	go func(){

		listen, err := net.Listen("tcp",":6060")
		if err != nil{
			log.Fatal(err)
		}
		for{
			fmt.Println("OK------------------")
			conn, err := listen.Accept()
			if err == nil{

				info_ch := make(chan string)
				is_connected := true
				//接收信息
				go func(info_ch chan string){
					for is_connected{
						var param string
						fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHH")
						raw := make([]byte,2048)
						len, err := conn.Read(raw) //ioutil.ReadAll(conn)
						//raw, err := ioutil.ReadAll(conn)
						if err != nil {
							//log.Fatal(err)
							conn.Close()
							is_connected = false
						}
						buff := bytes.NewBuffer(raw[:len])

						dec := gob.NewDecoder(buff)
						err = dec.Decode(&param)
						if err == nil {
							fmt.Println(param)
							rft.Start(param)//日志
							info_ch <-param
						} else {
							log.Println(err)
							fmt.Println("KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK")
							is_connected = false
						}
					}
				}(info_ch)
				//发送信息
				go func(info_cha chan string){
					for is_connected{
						/*select{
						case msg:= <-info_ch:
							select{
							case apply := <-ch:

							}


						}*/

						//fmt.Println(msg, apply)
						//writer := bufio.NewWriter(conn)
						buff := new(bytes.Buffer)
						enc := gob.NewEncoder(buff)

						err = enc.Encode("server received OK"+ <-info_ch)
						fmt.Println("PPPPPPPPPPPPPPPPPPPPPPPPPPP")
						if err != nil {
							//log.Fatal(err)
							conn.Close()
							is_connected = false
						}
						conn.Write(buff.Bytes())
						//writer.Write(buff.Bytes())
						//writer.Flush()
						//time.Sleep(time.Second)

					}
				}(info_ch)
			}else{
				fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
				log.Println(err)
			}

		}
	}()

	// Logs
	go func(){

		for{

			fmt.Println("msg applyed.", <-ch)

		}
	}()

	//simulate client thread.
	go func(){
		for{
			time.Sleep(time.Second*60)
			//index, term, isleader := rft.Start("Append E 1")
			//fmt.Println("[ client ]",index, term, isleader)
		}
	}()
	//注册xraft服务 开启监听
	serv := grpc.NewServer()
	pb.RegisterXRaftServer(serv, rft)
	reflection.Register(serv)
	serv.Serve(lis)


	//fmt.Println(rft.GetState())
	//fmt.Println(os.Args[:1])
}
