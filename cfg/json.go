// REF: https://stackoverflow.com/a/16466189
package cfg

import (
    "encoding/json"
	"flag"
    "os"
    "fmt"
)

type Config struct {
	Server struct {
		Addr string
	}
	Ably struct {
		Key string
	}
	Path struct {
		Asset string
	}
}

var config Config
var cfgMap map[string]any

func LoadJSON() {
	cfgPath := flag.String("config", "conf.json", "config file path")
	addr := flag.String("addr", "", "server address")
	flag.Parse()

	file, _ := os.Open(*cfgPath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
	  fmt.Println("error:", err)
	}
	if "" != *addr {
		config.Server.Addr = *addr
	}

	obj, _ := json.Marshal(config)
	json.Unmarshal(obj, &cfgMap)
}

func GetStr(seqs []string) string {
	obj := cfgMap
	for _, k := range seqs[:len(seqs)-1] {
		obj = obj[k].(map[string]any)
	}

	return obj[seqs[len(seqs)-1]].(string)
}

func Get(seqs []string) any {
	obj := cfgMap
	for _, k := range seqs[:len(seqs)-1] {
		obj = obj[k].(map[string]any)
	}

	return obj[seqs[len(seqs)-1]]
}
