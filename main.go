package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Config struct {
	URL         string `json:"url"`
	Interval    int    `json:"interval"`
	Threads     int    `json:"threads"`
	Format      string `json:"format"`
	FolderName  string `json:"folderName"`
	Qtd         int    `json:"qtd"`
	RampaInicio int    `json:"rampaInicio"`
	Forever     bool   `json:"forever"`
	Persist     bool   `json:"persist"`
}
type Tester struct {
	config      Config
	currentDir  string
	wg          *sync.WaitGroup
	count       Counter
	namingMutex *sync.Mutex
}
type Counter struct {
	count uint32
	m     *sync.RWMutex
}

func main() {
	const maxThreads = 16
	prev := runtime.GOMAXPROCS(maxThreads)
	fmt.Println("Definindo número máximo de threads do SO para", maxThreads, "- Anterior:", prev)

	configName := flag.String("config", "config.json", "Caminho para o arquivo de configuração")
	flag.Parse()

	tester := NewTester(*configName)
	sleepTime := time.Millisecond * time.Duration(tester.config.RampaInicio/tester.config.Threads)

	for i := 0; i < tester.config.Threads; i++ {
		tester.wg.Add(1)

		go tester.Run()

		time.Sleep(sleepTime)
	}

	tester.wg.Wait()
}

func (te *Tester) readConfig(configName string) {
	raw, err := ioutil.ReadFile(configName)

	if err != nil {
		log.Println("Erro ao abrir aquivo json")
		os.Exit(1)
	}

	if err := json.Unmarshal(raw, &te.config); err != nil {
		log.Println("Erro ao decodificar config.json:", err.Error())
		os.Exit(1)
	}
}

func (te *Tester) createDirectory() {
	err := os.MkdirAll("data/"+te.config.FolderName, os.ModePerm)
	if err != nil {
		log.Print("Erro ao criar diretório para salvar os dados:", err.Error())
		os.Exit(1)
	}
}

func (te *Tester) getFileName(t time.Time) string {
	te.namingMutex.Lock()
	defer te.namingMutex.Unlock()

	tStr := t.Format("2006_01_02_15_04_05.000000")
	timeString := strings.ReplaceAll(tStr, ".", "_")
	seq := fmt.Sprintf("%06v", te.Counter())

	return filepath.Join(te.currentDir, "data", te.config.FolderName, timeString+"_"+seq+"."+te.config.Format)
}

func (te *Tester) DoRequest(t time.Time) {
	statusCode, body, err := fasthttp.Get(nil, te.config.URL)
	if err != nil {
		log.Println("Erro na requisição:", err.Error())
		log.Println("Status:", statusCode)
	}

	fileName := te.getFileName(t)

	if te.config.Persist {
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(body)
		if err != nil {
			log.Println("Erro ao escrever resposta no arquivo:", err.Error())
		}

		f.Close()
	}
}
func (te *Tester) Run() {
	for i := 0; i < te.config.Qtd; i++ {
		if te.config.Forever {
			i--
		}

		t := time.Now()
		te.DoRequest(t)

		sleepTime := time.Millisecond*time.Duration(te.config.Interval) - time.Since(t)
		if sleepTime < 0 {
			sleepTime = time.Duration(0)
		}

		time.Sleep(sleepTime)
	}
	te.wg.Done()
}
func NewTester(configName string) *Tester {
	te := &Tester{}
	te.config = Config{}
	te.wg = &sync.WaitGroup{}
	te.namingMutex = &sync.Mutex{}
	te.count = Counter{
		count: 0,
		m:     &sync.RWMutex{},
	}

	te.readConfig(configName)
	fmt.Println("Configurações carregadas com sucesso.")

	if te.config.Forever {
		te.config.Qtd = 10
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	te.currentDir = filepath.Dir(ex)

	te.createDirectory()

	return te
}

func (te *Tester) Counter() (val uint32) {
	te.count.m.Lock()
	defer te.count.m.Unlock()
	val = te.count.count
	te.count.count++

	return
}
