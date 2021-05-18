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
var user_zchien = "U696bcb700dfc9254b27605374b86968b"
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
		log.Println("ساعت بوقت تهران:  " + nowString)
		bot.ReplyMessage(replyToken, linebot.NewTextMessage("ساعت بوقت تهران: " + nowString)).Do()
	} else if silent != true {
		log.Println("ساعت بوقت تهران:  " + nowString)
		bot.PushMessage(replyToken, linebot.NewTextMessage("ساعت بوقت تهران:  " + nowString)).Do()
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
				
				if strings.Contains(message.Text, "test") {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("success")).Do()
				} else if "groupid"  == message.Text {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(string(source.GroupID))).Do()
				} else if "help" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("help\n・[image:画像url]=從圖片網址發送圖片\n・[speed]=測回話速度\n・[groupid]=發送GroupID\n・[roomid]=發送RoomID\n・[byebye]=取消訂閱\n・[about]=作者\n・[me]=發送發件人信息\n・[test]=test bowwow是否正常\n・[now]=現在時間\n・[mid]=mid\n・[sticker]=隨機圖片\n\n[其他機能]\n位置測試\n捉貼圖ID\n加入時發送消息")).Do()
				} else if "mid" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(source.UserID)).Do()
				} else if "mee6" == message.Text  {
					bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://www.itsfun.com.tw/cacheimg/41/ac/f9e580eb9cbd38241182198bcb1b.jpg","https://www.itsfun.com.tw/cacheimg/41/ac/f9e580eb9cbd38241182198bcb1b.jpg")).Do()
				} else if "mee5" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewImageMessage("https://www.itsfun.com.tw/cacheimg/41/ac/f9e580eb9cbd38241182198bcb1b.jpg","https://www.itsfun.com.tw/cacheimg/41/ac/f9e580eb9cbd38241182198bcb1b.jpg")).Do()
					bot.ReplyMessage(replyToken, linebot.NewImageMessage("https://www.itsfun.com.tw/cacheimg/1a/e7/e18b42a3703133301659f6ddd7a4.jpg","https://www.itsfun.com.tw/cacheimg/1a/e7/e18b42a3703133301659f6ddd7a4.jpg")).Do()
				} else if "roomid" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(source.RoomID)).Do()
				} else if "hidden" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("hidden")).Do()
				} else if "bowwow" == message.Text  {
					bot.ReplyMessage(replyToken, linebot.NewImageMessage("https://www.itsfun.com.tw/cacheimg/66/0c/eda9d251c3bd769ac820552b2ff1.jpg","https://www.itsfun.com.tw/cacheimg/66/0c/eda9d251c3bd769ac820552b2ff1.jpg")).Do()}
					if err != nil {
						log.Fatal(err)
				} else if "me" == message.Text  {
					mid := source.UserID
					p, err := bot.GetProfile(mid).Do()
					if err != nil {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("新增同意"))
					}

					bot.ReplyMessage(replyToken, linebot.NewTextMessage("mid:"+mid+"\nname:"+p.DisplayName+"\nstatusMessage:"+p.StatusMessage)).Do()
				} else if res := strings.Contains(message.Text, "hello"); res == true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("hello!"), linebot.NewTextMessage("my name is bowwow")).Do()
				} else if res := strings.Contains(message.Text, "image:"); res == true {
					image_url := strings.Replace(message.Text, "image:", "", -1)
					bot.ReplyMessage(replyToken, linebot.NewImageMessage(image_url, image_url)).Do()
				} else if "1" == message.Text {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀\n\n─═≡ϻఠ_ఠsɛɳ≡═─\n\n.1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J\n\n─═≡ϻఠ_ఠsɛɳ≡═─\n\n.1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9..0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.\n\n💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀"), linebot.NewTextMessage("7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7."), linebot.NewTextMessage("7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7."), linebot.NewTextMessage("7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7.7W0.G7.W0.G7W0.G7.W0.G7W0.G7.W0.G7W0.G7."), linebot.NewTextMessage("💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀\n\n─═≡ϻఠ_ఠsɛɳ≡═─\n\n.1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J\n\n─═≡ϻఠ_ఠsɛɳ≡═─\n\n.1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9..0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.\n\n💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀💀")).Do()
				} else if  "time" == message.Text {
					tellTime(replyToken, true)
				} else if  "4" == message.Text {
					silentMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H"),linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H"),linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H"), linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H"), linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H") ,linebot.NewTextMessage(".1.2.3.4.5.6.7.8.9.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H.7.I.8.J.9.K.0.A.1.B.2.D.3.E.4.F.5.G.6.H")).Do()
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
					imageURL := baseURL + "/static/buttons/1040.jpg"
					//log.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> "+imageURL)
					template := linebot.NewButtonsTemplate(
						imageURL, "My button sample", "Hello, my button",
						linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
						linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
						linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
						linebot.NewMessageTemplateAction("Say message", "Rice=米"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Buttons alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "confirm" == message.Text {
					template := linebot.NewConfirmTemplate(
						"Do it?",
						linebot.NewMessageTemplateAction("Yes", "Yes!"),
						linebot.NewMessageTemplateAction("No", "No!"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "carousel" == message.Text {
					imageURL := baseURL + "/static/buttons/1040.jpg"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "hoge", "fuga",
							linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
							linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
						),
						linebot.NewCarouselColumn(
							imageURL, "hoge", "fuga",
							linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
							linebot.NewMessageTemplateAction("Say message", "Rice=米"),
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
							baseURL + "/static/rich",
							"Imagemap alt text",
							linebot.ImagemapBaseSize{1040, 1040},
							linebot.NewURIImagemapAction("https://store.line.me/family/manga/en", linebot.ImagemapArea{0, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/music/en", linebot.ImagemapArea{520, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/play/en", linebot.ImagemapArea{0, 520, 520, 520}),
							linebot.NewMessageImagemapAction("URANAI!", linebot.ImagemapArea{520, 520, 520, 520}),
						),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "/bye" == message.Text {
					if rand.Intn(100) > 70 {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("請神容易送神難, 我偏不要, 嘿嘿")).Do()
					} else {
						switch source.Type {
						case linebot.EventSourceTypeUser:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我想走, 但是我走不了...")).Do()
						case linebot.EventSourceTypeGroup:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我揮一揮衣袖 不帶走一片雲彩")).Do()
							bot.LeaveGroup(source.GroupID).Do()
						case linebot.EventSourceTypeRoom:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我揮一揮衣袖 不帶走一片雲彩")).Do()
							bot.LeaveRoom(source.RoomID).Do()
						}
					}
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
