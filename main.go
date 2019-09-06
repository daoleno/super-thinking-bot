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
)

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
	chatID    string
	text      string
	parseMode string
}

func getPageSummary(title string) (summary []byte, err error) {

	const (
		wikiAPI string = "https://en.wikipedia.org/api/rest_v1/"
		psAPI   string = wikiAPI + "page/summary/"
	)

	url := psAPI + title
	log.Printf("Wiki url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	summary, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return summary, nil

}

func sendMsg(content wikiContent) error {

	const (
		tgAPI      string = "https://api.telegram.org/bot"
		token      string = ""
		baseURL    string = tgAPI + token
		sendMsgAPI string = baseURL + "/sendMessage"
	)

	var msg message
	msg.chatID = "@superthinking2u"
	msg.text = fmt.Sprintf("<b>%s</b>  %s", content.Title, content.ContentURL.Desktop.Page)
	msg.parseMode = "HTML"

	url := sendMsgAPI + "?chat_id=" + msg.chatID + "&text=" + msg.text + "&parse_mode=" + msg.parseMode
	log.Printf("Telegram url:%s", url)
	_, err := http.Get(url)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	var titles []string

	// Get mental models
	file, err := os.Open("mental-models.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		titles = append(titles, strings.ToLower(scanner.Text()))
	}

	for _, title := range titles {
		// Get summary
		summary, err := getPageSummary(title)
		if err != nil {
			log.Println(err)
		}

		// Generate wikiContent
		var content wikiContent
		err = json.Unmarshal(summary, &content)
		if err != nil {
			log.Println(err)
		}

		// Send to telegram
		err = sendMsg(content)
		if err != nil {
			log.Println(err)
		}
	}

}
