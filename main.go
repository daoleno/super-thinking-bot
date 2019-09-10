package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"gopkg.in/yaml.v2"
)

type config struct {
	Telegram struct {
		Channel string `yaml:"channel"`
		Token   string `yaml:"token"`
	} `yaml:"telegram"`
}

type contentURL struct {
	Desktop struct {
		Page string `json:"page"`
	} `json:"desktop"`

	Mobile struct {
		Page string `json:"page"`
	} `json:"mobile"`
}

type wikiContent struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ContentURL  contentURL `json:"content_urls"`
	Summary     string     `json:"extract"`
}

type message struct {
	ChatID    string
	Text      string
	ParseMode string
}

func getPageSummary(title string) (summary []byte, err error) {
	const (
		wikiAPI string = "https://en.wikipedia.org/api/rest_v1/"
		psAPI   string = wikiAPI + "page/summary/"
	)

	url := psAPI + title
	log.Printf("wikiURL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Status code:%d \t wikiAPI:%s \t title:%s", resp.StatusCode, url, title)
	}

	summary, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func sendMsg(content wikiContent) error {

	var (
		tgAPI      = "https://api.telegram.org/bot"
		Token      = globalConfig.Telegram.Token
		baseURL    = tgAPI + Token
		sendMsgAPI = baseURL + "/sendMessage"
	)

	var msg message
	msg.ChatID = globalConfig.Telegram.Channel
	msg.Text = fmt.Sprintf("<b>%s</b>  %s", content.Title, content.ContentURL.Desktop.Page)
	msg.ParseMode = "HTML"

	url := sendMsgAPI + "?chat_id=" + msg.ChatID + "&text=" + msg.Text + "&parse_mode=" + msg.ParseMode
	// log.Printf("telegram URL:%s", url)
	_, err := http.Get(url)
	if err != nil {
		return err
	}

	return nil
}

func readFile(cfg *config) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	globalConfig config
)

func main() {

	// Read config
	readFile(&globalConfig)

	// Get mental models
	file, err := os.Open("mental-models.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var titles []string
	for scanner.Scan() {
		titles = append(titles, strings.ToLower(scanner.Text()))
	}

	wikiChan := make(chan wikiContent)
	go func() {
		for _, title := range titles {
			// Get summary
			summary, err := getPageSummary(title)
			if err != nil {
				log.Println(err)
				continue
			}

			// Generate wikiContent
			var content wikiContent
			err = json.Unmarshal(summary, &content)
			if err != nil {
				log.Println(err)
				continue
			}

			wikiChan <- content

		}
	}()

	// Send to telegram
	loc, _ := time.LoadLocation("Asia/Shanghai")
	gocron.Every(1).Day().At("08:00").Loc(loc).Do(func() {
		log.Println("Sending to telegram at ", time.Now().In(loc))
		err := sendMsg(<-wikiChan)
		if err != nil {
			log.Println(err)
		}
	})
	<-gocron.Start()

}
