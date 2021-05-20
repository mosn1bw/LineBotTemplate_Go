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

type SelfIntro struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

func NewSelfIntro(channelSecret, channelToken string) (*SelfIntro, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &SelfIntro{
		bot:         bot,
		appBaseURL:  "nil",
		downloadDir: "nil",
	}, nil
}

func (s *SelfIntro) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := s.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		fmt.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := s.handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Println(err)
				}
			case *linebot.StickerMessage:
				if err := s.handleSticker(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		case linebot.EventTypeFollow:
			if err := s.handleJoin(event.ReplyToken, event.Source); err != nil {
				log.Println(err)
			}
		case linebot.EventTypeJoin:
			if err := s.handleJoin(event.ReplyToken, event.Source); err != nil {
				log.Println(err)
			}
		case linebot.EventTypeLeave:
			log.Printf("Left: %v", event)
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (s *SelfIntro) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	domain, keyword := ParseMessage(message.Text)
	switch domain {
	case "m2":
		works, err := readJSON("static/message/works.json")
		if err != nil {
			return err
		}
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(works))
		if err != nil {
			return err
		}
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 偵測到關鍵字 '%s'!\n 推斷你想要知道我的 '%s'！", keyword, domain)).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "098")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "098")),
			linebot.NewFlexMessage("作品集介紹", contents),
		).Do(); err != nil {
			return err
		}
	case "m1":
		experience, err := readJSON("static/message/experience.json")
		if err != nil {
			return err
		}
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(experience))
		if err != nil {
			return err
		}
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 偵測到關鍵字 '%s'!\n 推斷你想要知道我的 '%s'！", keyword, domain)).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "098")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "098")),
			linebot.NewFlexMessage("經歷介紹", contents),
		).Do(); err != nil {
			return err
		}
	case "m3":
		education, err := readJSON("static/message/education.json")
		if err != nil {
			return err
		}
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(education))
		if err != nil {
			return err
		}
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 偵測到關鍵字 '%s'!\n 推斷你想要知道我的 '%s'！", keyword, domain)).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "098")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "098")),
			linebot.NewFlexMessage("學歷介紹", contents),
		).Do(); err != nil {
			return err
		}
	case "m4":
		profile, _ := s.bot.GetProfile(source.UserID).Do()
		intro, err := readJSON("static/message/intro.json")
		if err != nil {
			return err
		}
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(intro))
		if err != nil {
			return err
		}
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 歡迎 %s!!\n 按下方的按鈕來認識我吧！", profile.DisplayName)).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "098")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "098")),
			linebot.NewFlexMessage("自我介紹", contents),
		).Do(); err != nil {
			return err
		}
	case "m5":
		skills, err := readJSON("static/message/skills.json")
		if err != nil {
			return err
		}
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(skills))
		if err != nil {
			return err
		}
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 偵測到關鍵字 '%s'!\n 推斷你想要知道我的 '%s'！", keyword, domain)).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "098")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "098")),
			linebot.NewFlexMessage("技能介紹", contents),
		).Do(); err != nil {
			return err
		}
	default:
		log.Printf("Echo message to %s: %s", replyToken, message.Text)
		if _, err := s.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(fmt.Sprintf("$$ 謝謝您傳訊息給James!\n 可以按下方主選單問我問題或輸入'介紹自己'喔！")).AddEmoji(
				linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "094")).AddEmoji(
				linebot.NewEmoji(1, "5ac1bfd5040ab15980c9b435", "094")),
		).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (s *SelfIntro) handleSticker(message *linebot.StickerMessage, replyToken string) error {
	if _, err := s.bot.ReplyMessage(
		replyToken,
		linebot.NewStickerMessage(message.PackageID, message.StickerID),
	).Do(); err != nil {
		return err
	}
	return nil
}

func (s *SelfIntro) replyText(replyToken, text string) error {
	if _, err := s.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

func readJSON(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	return byteValue, nil
}
