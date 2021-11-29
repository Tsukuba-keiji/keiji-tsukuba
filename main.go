// This code is from https://github.com/line/line-bot-sdk-go, arrenged by us.

// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	app, err := NewKitchenSink(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// serve /static/** files
	staticFileServer := http.FileServer(http.Dir("static"))
	http.HandleFunc("/static/", http.StripPrefix("/static/", staticFileServer).ServeHTTP)
	// serve /downloaded/** files
	downloadedFileServer := http.FileServer(http.Dir(app.downloadDir))
	http.HandleFunc("/downloaded/", http.StripPrefix("/downloaded/", downloadedFileServer).ServeHTTP)
	// serve /assets/** files
	assetsFileServer := http.FileServer(http.Dir("assets"))
	http.HandleFunc("/", http.StripPrefix("/", assetsFileServer).ServeHTTP)

	http.HandleFunc("/callback", app.Callback)
	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

// KitchenSink app
type KitchenSink struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

// NewKitchenSink function
func NewKitchenSink(channelSecret, channelToken, appBaseURL string) (*KitchenSink, error) {
	apiEndpointBase := os.Getenv("ENDPOINT_BASE")
	if apiEndpointBase == "" {
		apiEndpointBase = linebot.APIEndpointBase
	}
	bot, err := linebot.New(
		channelSecret,
		channelToken,
		linebot.WithEndpointBase(apiEndpointBase), // Usually you omit this.
	)
	if err != nil {
		return nil, err
	}
	downloadDir := filepath.Join(filepath.Dir(os.Args[0]), "line-bot")
	_, err = os.Stat(downloadDir)
	if err != nil {
		if err := os.Mkdir(downloadDir, 0777); err != nil {
			return nil, err
		}
	}
	return &KitchenSink{
		bot:         bot,
		appBaseURL:  appBaseURL,
		downloadDir: downloadDir,
	}, nil
}

// Callback function for http server
func (app *KitchenSink) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := app.handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		case linebot.EventTypePostback:
			data := event.Postback.Data
			if data == "search" {
				jsondata := loadJson("src/user.json")
				if _, exist := jsondata.(map[string]interface{})[event.Source.UserID]; !exist {
					imageURL := app.appBaseURL + "/static/img/config.png"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "学年を選択", "情報を受け取る学年を全て選んでください",
							linebot.NewPostbackAction("一年生", "grade=1", "", ""),
							linebot.NewPostbackAction("二年生", "grade=2", "", ""),
							linebot.NewPostbackAction("三年生", "grade=3", "", ""),
						),
						linebot.NewCarouselColumn(
							imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
							linebot.NewPostbackAction("一組", "class=1", "", ""),
							linebot.NewPostbackAction("二組", "class=2", "", ""),
							linebot.NewPostbackAction("三組", "class=3", "", ""),
						),
						linebot.NewCarouselColumn(
							imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
							linebot.NewPostbackAction("四組", "class=4", "", ""),
							linebot.NewPostbackAction("五組", "class=5", "", ""),
							linebot.NewPostbackAction("六組", "class=6", "", ""),
						),
					)
					usermap := make(map[string]string)
					slice := []interface{}{usermap}
					jsondata.(map[string]interface{})[event.Source.UserID] = slice
					saveJson(jsondata, "src/user.json")
					if _, err := app.bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("学年とクラスを設定してください。完了したらもう一度検索を選択してください。"),
						linebot.NewTemplateMessage("検索設定", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else {
					_, existg := jsondata.(map[string]interface{})[event.Source.UserID].([]interface{})[0].(map[string]interface{})["grade"]
					_, existc := jsondata.(map[string]interface{})[event.Source.UserID].([]interface{})[0].(map[string]interface{})["class"]
					if existc && existg {
						user := jsondata.(map[string]interface{})[event.Source.UserID].([]interface{})[0]
						grade := user.(map[string]interface{})["grade"].(string)
						class := user.(map[string]interface{})["class"].(string)
						replyurl := "https://keiji-tsukuba.herokuapp.com/search?grade=" + grade + "&class=" + class
						if _, err := app.bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage(replyurl),
							linebot.NewTextMessage("上記URLにて連絡を確認して下さい。"),
						).Do(); err != nil {
							log.Print(err)
						}
					} else {
						imageURL := app.appBaseURL + "/static/img/config.png"
						template := linebot.NewCarouselTemplate(
							linebot.NewCarouselColumn(
								imageURL, "学年を選択", "情報を受け取る学年を全て選んでください",
								linebot.NewPostbackAction("一年生", "grade=1", "", ""),
								linebot.NewPostbackAction("二年生", "grade=2", "", ""),
								linebot.NewPostbackAction("三年生", "grade=3", "", ""),
							),
							linebot.NewCarouselColumn(
								imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
								linebot.NewPostbackAction("一組", "class=1", "", ""),
								linebot.NewPostbackAction("二組", "class=2", "", ""),
								linebot.NewPostbackAction("三組", "class=3", "", ""),
							),
							linebot.NewCarouselColumn(
								imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
								linebot.NewPostbackAction("四組", "class=4", "", ""),
								linebot.NewPostbackAction("五組", "class=5", "", ""),
								linebot.NewPostbackAction("六組", "class=6", "", ""),
							),
						)
						if _, err := app.bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("学年とクラスが設定されていません。完了したらもう一度検索を選択してください。"),
							linebot.NewTemplateMessage("検索設定", template),
						).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			} else if data == "add" {
				if _, err := app.bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("https://keiji-tsukuba.herokuapp.com/add"),
					linebot.NewTextMessage("上のURLでデータを入力してください。"),
				).Do(); err != nil {
					log.Print(err)
				}
			} else if data[:5] == "grade" {
				jsondata := loadJson("src/user.json")
				jsondata.(map[string]interface{})[event.Source.UserID].([]interface{})[0].(map[string]interface{})["grade"] = data[6:]
				saveJson(jsondata, "src/user.json")
				if _, err := app.bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("学年を"+data[6:]+"として設定しました。"),
				).Do(); err != nil {
					log.Print(err)
				}
			} else if data[:5] == "class" {
				jsondata := loadJson("src/user.json")
				jsondata.(map[string]interface{})[event.Source.UserID].([]interface{})[0].(map[string]interface{})["class"] = data[6:]
				saveJson(jsondata, "src/user.json")
				if _, err := app.bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("クラスを"+data[6:]+"として設定しました。"),
				).Do(); err != nil {
					log.Print(err)
				}
			} else {
				log.Printf("Unknown event : %v", data)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *KitchenSink) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "設定を変更":
		imageURL := app.appBaseURL + "/static/img/config.png"
		template := linebot.NewCarouselTemplate(
			linebot.NewCarouselColumn(
				imageURL, "学年を選択", "情報を受け取る学年を全て選んでください",
				linebot.NewPostbackAction("一年生", "grade=1", "", ""),
				linebot.NewPostbackAction("二年生", "grade=2", "", ""),
				linebot.NewPostbackAction("三年生", "grade=3", "", ""),
			),
			linebot.NewCarouselColumn(
				imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
				linebot.NewPostbackAction("一組", "class=1", "", ""),
				linebot.NewPostbackAction("二組", "class=2", "", ""),
				linebot.NewPostbackAction("三組", "class=3", "", ""),
			),
			linebot.NewCarouselColumn(
				imageURL, "クラスを選択", "情報を受け取るクラスを全て選んでください",
				linebot.NewPostbackAction("四組", "class=4", "", ""),
				linebot.NewPostbackAction("五組", "class=5", "", ""),
				linebot.NewPostbackAction("六組", "class=6", "", ""),
			),
		)
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("学年とクラスを設定してください。完了したらもう一度検索を選択してください。"),
			linebot.NewTemplateMessage("検索設定", template),
		).Do(); err != nil {
			log.Print(err)
		}
	}
	return nil
}

//load json
func loadJson(inputPath string) interface{} {
	byteArray, _ := ioutil.ReadFile(inputPath)
	var jsonObj interface{}
	_ = json.Unmarshal(byteArray, &jsonObj)
	return jsonObj
}

//save json
func saveJson(jsonObj interface{}, outputPath string) {
	file, _ := os.Create(outputPath)
	defer file.Close()
	_ = json.NewEncoder(file).Encode(jsonObj)
}
