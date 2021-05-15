// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// https://github.com/line/line-bot-sdk-go/tree/master/linebot

package main

import (
	"strconv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Constants
var timeFormat = "01/02 PM03:04:05"
var user_mosen = "u2023c2d6c4de3dc7c266f3f07cfabdcc"
var user_yaoming = "U3aaab6c6248bb38f194134948c60f757"
var user_jackal = "U3effab06ddf5bcf0b46c1c60bcd39ef5"
var user_shane = "U2ade7ac4456cb3ca99ffdf9d7257329a"

// Global Settings
var channelSecret = os.Getenv("CHANNEL_SECRET")
var channelToken = os.Getenv("CHANNEL_TOKEN")
//var baseURL = os.Getenv("APP_BASE_URL")
var baseURL = "https://line-talking-bot-go.herokuapp.com"
var endpointBase = os.Getenv("ENDPOINT_BASE")
var tellTimeInterval int = 15
var answers_TextMessage = []string{
		"",
	}
var answers_ImageMessage = []string{
		"",
	}
var answers_StickerMessage = []string{
		"",
	}
var answers_VideoMessage = []string{
		"",
	}
var answers_AudioMessage = []string{
		"",
	}
var answers_LocationMessage = []string{
		"",
	}
var answers_ReplyCurseMessage = []string{
		"",
	}

var silentMap = make(map[string]bool) // [UserID/GroupID/RoomID]:bool

//var echoMap = make(map[string]bool)

var loc, _ = time.LoadLocation("Asia/Tehran")
var bot *linebot.Client


func tellTime(replyToken string, doTell bool){
	var silent = false
	now := time.Now().In(loc)
	nowString := now.Format(timeFormat)
	
	if doTell {
		log.Println("ساعت بوقت(Tehran): " + nowString)
		bot.ReplyMessage(replyToken, linebot.NewTextMessage("ساعت بوقت(Tehran): " + nowString)).Do()
	} else if silent != true {
		log.Println("ساعت بوقت(Tehran): " + nowString)
		bot.PushMessage(replyToken, linebot.NewTextMessage("ساعت بوقت(Tehran): " + nowString)).Do()
	} else {
		log.Println("tell time misfired")
	}
}

func tellTimeJob(sourceId string) {
	for {
		time.Sleep(time.Duration(tellTimeInterval) * time.Minute)
		now := time.Now().In(loc)
		log.Println("time to tell time to : " + sourceId + ", " + now.Format(timeFormat))
		tellTime(sourceId, false)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	/*
	go func() {
		tellTimeJob(user_zchien);
	}()
	go func() {
		for {
			now := time.Now().In(loc)
			log.Println("keep alive at : " + now.Format(timeFormat))
			//http.Get("https://line-talking-bot-go.herokuapp.com")
			time.Sleep(time.Duration(rand.Int31n(29)) * time.Minute)
		}
	}()
	*/

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	log.Print("URL:"  + r.URL.String())
	
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		var replyToken = event.ReplyToken

		var source = event.Source //EventSource		
		var userId = source.UserID
		var groupId = source.GroupID
		var roomId = source.RoomID
		log.Print("callbackHandler to source UserID/GroupID/RoomID: " + userId + "/" + groupId + "/" + roomId)
		
		var sourceId = roomId
		if sourceId == "" {
			sourceId = groupId
			if sourceId == "" {
				sourceId = userId
			}
		}
		
		if event.Type == linebot.EventTypeMessage {
			_, silent := silentMap[sourceId]
			
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				log.Print("ReplyToken[" + replyToken + "] TextMessage: ID(" + message.ID + "), Text(" + message.Text  + "), current silent status=" + strconv.FormatBool(silent) )
				//if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}
				
				if source.UserID != "" && source.UserID != user_mosen {
					profile, err := bot.GetProfile(source.UserID).Do()
					if err != nil {
						log.Print(err)
					} else if _, err := bot.PushMessage(user_mosen, linebot.NewTextMessage(profile.DisplayName + ": "+message.Text)).Do(); err != nil {
							log.Print(err)
					}
				}
				
				if strings.Contains(message.Text, "1") {
					silentMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
				if strings.Contains(message.Text, "2") {
					silentMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
				} else if strings.Contains(message.Text, "time") {
					tellTime(replyToken, true)
				} else if "1" == message.Text {
					silentMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("，1、2、3... OK")).Do()
				} else if "profile" == message.Text {
					if source.UserID != "" {
						profile, err := bot.GetProfile(source.UserID).Do()
						if err != nil {
							log.Print(err)
						} else if _, err := bot.ReplyMessage(
							replyToken,
							linebot.NewTextMessage("Display name: "+profile.DisplayName + ", Status message: "+profile.StatusMessage)).Do(); err != nil {
								log.Print(err)
						}
					} else {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("Bot can't use profile API without user ID")).Do()
					}
				} else if "buttons" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-IVJ0bg14co4/YBq4zQOEN0I/AAAAAAAAL6Q/ojEHrB9Uju8Cj4nQ1FTHun-6XKHYZd_vACK8BGAsYHg/s340/2021-02-03.gif"
					//log.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> "+imageURL)
					template := linebot.NewButtonsTemplate(
						imageURL, "My button sample", "«ϻơsɛɳ»",
						linebot.NewURITemplateAction("Go to line.me",  "line://ti/p/~M_BW"),
						linebot.NewMessageTemplateAction("Say message", "«ϻơsɛɳ»"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Buttons alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "buttons1" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-IVJ0bg14co4/YBq4zQOEN0I/AAAAAAAAL6Q/ojEHrB9Uju8Cj4nQ1FTHun-6XKHYZd_vACK8BGAsYHg/s340/2021-02-03.gif"
					//log.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> "+imageURL)
					template := linebot.NewButtonsTemplate(
						imageURL, "My button sample", "Hello, my button",
						linebot.NewURITemplateAction("Go to line.me",  "line://ti/p/~M_BW"),
						linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Buttons alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "vpn" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-xHqQP4wTZDU/YBq5AgqjvCI/AAAAAAAAL6c/TmVGaX4tgIk07K5bZIPDtV9Ct49xEwaxwCK8BGAsYHg/s512/2021-02-03.gif"
					//log.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> "+imageURL)
					template := linebot.NewButtonsTemplate(
						imageURL, "فیلتر شکن", "روی یکی از منوهای زیر کلیک کنید",
						linebot.NewURITemplateAction("supervpnfree", "https://play.google.com/store/apps/details?id=com.jrzheng.supervpnfree&hl=fa"),
						linebot.NewURITemplateAction("totally", "https://play.google.com/store/apps/details?id=net.rejinderi.totallyfreevpn&hl=fa"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Buttons alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "datetime" == message.Text {
					template := linebot.NewButtonsTemplate(
						"", "", "Select date / time !",
						linebot.NewDatetimePickerAction("date", "DATE", "date", "", "", ""),
						linebot.NewDatetimePickerAction("time", "TIME", "time", "", "", ""),
						linebot.NewDatetimePickerAction("datetime", "DATETIME", "datetime", "", "", ""),
					)
					if _, err := app.bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Datetime pickers alt text", template),
					).Do(); err != nil {
						return err
					}
				} else if "mee2"== message.Text {
					imageURL := app.appBaseURL + "https://imgurl.ir/uploads/g643845_.gif"
					template := linebot.NewImageCarouselTemplate(
						linebot.NewImageCarouselColumn(
							imageURL,
							linebot.NewURIAction("Go to LINE", "https://line.me/ti/p/~m_bw"),
						),
					)
					if _, err := app.bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("[┅═★ᖼᗱO꓅★═┅]", template),
					).Do(); err != nil {
						return err
					}
				} else if "ربات" == message.Text {
					template := linebot.NewConfirmTemplate(
						"از ربات راضی هستید?",
						linebot.NewMessageTemplateAction("Yes", "Yes!"),
						linebot.NewMessageTemplateAction("No", "No!"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "carousel1" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-CKPfi57SLOs/YGtXrTQ30ZI/AAAAAAAAMUU/SkJbo6DV4S0m7QmM3Dpsbl9BWpgA6uWJwCK8BGAsYHg/s500/2021-04-05.gif"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "мosᴇɴ", "ʙoт",
							linebot.NewURITemplateAction("Go to line.me", "line://ti/p/~M_BW"),
						),
						linebot.NewCarouselColumn(
							imageURL, "мosᴇɴ", "ʙoт",
							linebot.NewMessageTemplateAction("мosᴇɴ", "✿мosᴇɴ👿ʙoт✿"),
						),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Carousel alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "me2" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-CKPfi57SLOs/YGtXrTQ30ZI/AAAAAAAAMUU/SkJbo6DV4S0m7QmM3Dpsbl9BWpgA6uWJwCK8BGAsYHg/s500/2021-04-05.gif"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "",
							linebot.NewMessageTemplateAction("★ᖼOᗱᗴℕ★", "ᖼOᗱᗴℕ"),
						),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Carousel alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
                    
				} else if "carousel" == message.Text {
					imageURL := "https://lh3.googleusercontent.com/-CKPfi57SLOs/YGtXrTQ30ZI/AAAAAAAAMUU/SkJbo6DV4S0m7QmM3Dpsbl9BWpgA6uWJwCK8BGAsYHg/s500/2021-04-05.gif"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "мosᴇɴ", "ʙoт",
							linebot.NewURITemplateAction("Go to line.me",  "line://ti/p/~M_BW"),
						),
						linebot.NewCarouselColumn(
							imageURL, "мosᴇɴ", "ʙoт",
							linebot.NewPostbackTemplateAction("мosᴇɴ", "hello ", "hello こんにちは"),
						),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Carousel alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "imagemap" == message.Text {
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewImagemapMessage(
							baseURL + "https://8pic.ir/uploads/045ceedfbc8efc4bb05d612cebe22ae0.jpg",
							"Imagemap alt text",
							linebot.ImagemapBaseSize{1040, 1040},
							linebot.NewURIImagemapAction("https://store.line.me/family/manga/en", linebot.ImagemapArea{0, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/music/en", linebot.ImagemapArea{520, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/play/en", linebot.ImagemapArea{0, 520, 520, 520}),
							linebot.NewMessageImagemapAction("URANAI!", linebot.ImagemapArea{520, 520, 520, 520}),
						),
					).Do(); err != nil {
						log.Print(err)
				} else if "/bye" == message.Text {
					switch source.Type {
					case linebot.EventSourceTypeUser:
						return app.replyText(replyToken, "Bot can't leave from 1:1 chat")
					case linebot.EventSourceTypeGroup:
						if err := app.replyText(replyToken, "Leaving group"); err != nil {
							return err
						}
						if _, err := app.bot.LeaveGroup(source.GroupID).Do(); err != nil {
							return app.replyText(replyToken, err.Error())
						}
					case linebot.EventSourceTypeRoom:
						if err := app.replyText(replyToken, "Leaving room"); err != nil {
							return err
						}
						if _, err := app.bot.LeaveRoom(source.RoomID).Do(); err != nil {
							return app.replyText(replyToken, err.Error())
						}
					}
				default:
					log.Printf("Echo message to %s: %s", replyToken, message.Text)
					if _, err := app.bot.ReplyMessage(
						replyToken,{
						linebot.NewTextMessage(message.Text),
					).Do(); err != nil {
						return err
					}
                                {    
				return nil
				} else if "無恥" == message.Text {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_ReplyCurseMessage[rand.Intn(len(answers_ReplyCurseMessage))])).Do()
				} else if silentMap[sourceId] != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_TextMessage[rand.Intn(len(answers_TextMessage))])).Do()
				}
			case *linebot.ImageMessage :
				log.Print("ReplyToken[" + replyToken + "] ImageMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_ImageMessage[rand.Intn(len(answers_ImageMessage))])).Do()
				}
			case *linebot.VideoMessage :
				log.Print("ReplyToken[" + replyToken + "] VideoMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_VideoMessage[rand.Intn(len(answers_VideoMessage))])).Do()
				}
			case *linebot.AudioMessage :
				log.Print("ReplyToken[" + replyToken + "] AudioMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), Duration(" + strconv.Itoa(message.Duration) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_AudioMessage[rand.Intn(len(answers_AudioMessage))])).Do()
				}
			case *linebot.LocationMessage:
				log.Print("ReplyToken[" + replyToken + "] LocationMessage[" + message.ID + "] Title(" + message.Title  + "), Address(" + message.Address + "), Latitude(" + strconv.FormatFloat(message.Latitude, 'f', -1, 64) + "), Longitude(" + strconv.FormatFloat(message.Longitude, 'f', -1, 64) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_LocationMessage[rand.Intn(len(answers_LocationMessage))])).Do()
				}
			case *linebot.StickerMessage :
				log.Print("ReplyToken[" + replyToken + "] StickerMessage[" + message.ID + "] PackageID(" + message.PackageID + "), StickerID(" + message.StickerID + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_StickerMessage[rand.Intn(len(answers_StickerMessage))])).Do()
				}
			}
		} else if event.Type == linebot.EventTypePostback {
		} else if event.Type == linebot.EventTypeBeacon {
		}
	}
	
}
