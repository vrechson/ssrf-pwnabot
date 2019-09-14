package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/whoismath/ssrf-pwnabot/commands"
	"github.com/whoismath/ssrf-pwnabot/config"
	botFactory "github.com/whoismath/ssrf-pwnabot/factory/bot"
)

// SSRFPwnaBot aaaa
type SSRFPwnaBot struct {
	conf *config.Config
	bot  *tgbotapi.BotAPI
}

// CreateApp new bot
func CreateApp(conf *config.Config) *SSRFPwnaBot {

	// create telegram bot
	bot, err := botFactory.BotFactory(conf)
	if err != nil {
		log.Fatal("Can't create a telegram bot: ", err)
	}

	return &SSRFPwnaBot{conf, bot}
}

// Start aaa
func (SSRFPwnaBot *SSRFPwnaBot) Start() {
	updates, err := SSRFPwnaBot.getUpdates()
	if err != nil {
		log.Fatal("can't get updates")
	}
	fmt.Println("starting bot")

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			go func() {
				SSRFPwnaBot.handleCommand(&update)
			}()

		}

	}
}

func (SSRFPwnaBot *SSRFPwnaBot) getUpdates() (tgbotapi.UpdatesChannel, error) {
	if SSRFPwnaBot.conf.UseWebhook != true {
		return SSRFPwnaBot.setupPolling()
	}
	return SSRFPwnaBot.setupWebhook()
}

func (SSRFPwnaBot *SSRFPwnaBot) setupPolling() (tgbotapi.UpdatesChannel, error) {
	SSRFPwnaBot.bot.RemoveWebhook()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 5
	fmt.Println("Start polling.")
	return SSRFPwnaBot.bot.GetUpdatesChan(updateConfig)
}

func (SSRFPwnaBot *SSRFPwnaBot) setupWebhook() (tgbotapi.UpdatesChannel, error) {
	_, err := SSRFPwnaBot.bot.SetWebhook(tgbotapi.NewWebhook(SSRFPwnaBot.conf.WebhookURL + "/" + SSRFPwnaBot.bot.Token))
	if err != nil {
		log.Fatal("webhook problem: ", err)
		//return nil, err
	}

	updates := SSRFPwnaBot.bot.ListenForWebhook("/" + SSRFPwnaBot.bot.Token)
	go http.ListenAndServe(":"+SSRFPwnaBot.conf.Port, nil)

	fmt.Println("Listening for connection")

	return updates, nil

}

func (SSRFPwnaBot *SSRFPwnaBot) handleCommand(update *tgbotapi.Update) {
	command := update.Message.Command()
	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "some text")

	switch command {

	// Showing help menu
	case commands.Help:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "use /random to return a random image from a website, example:\n/random https://illuminat-us.tumblr.com/")
		SSRFPwnaBot.bot.Send(msg)

	case commands.Random:
		client := &http.Client{}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Selecting a random image from this source")
		SSRFPwnaBot.bot.Send(msg)

		//response, err := http.Get(update.Message.CommandArguments())
		req, err := http.NewRequest("GET", update.Message.CommandArguments(), nil)
		req.Header.Add("X-Services", "8080 keep our secrets")
		response, err := client.Do(req)

		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "This URL isn't available or doesn't exist\n(remember to use http:// or https://)")
			SSRFPwnaBot.bot.Send(msg)
		} else {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Cannot load this site source")
				SSRFPwnaBot.bot.Send(msg)
			}

			image := findRandomImage(string(body))

			if image != "none" {
				res, _ := regexp.MatchString("\\A^https?://\\S", image)
				if res {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, image)
					SSRFPwnaBot.bot.Send(msg)
				} else {
					re, _ := regexp.Compile("^(?:\\/\\/|[^\\/]+)*")
					domain := re.FindString(string(update.Message.CommandArguments())) + "/" + image
					msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
					msg.FileID = domain
					msg.UseExisting = true
					SSRFPwnaBot.bot.Send(msg)
				}
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Couldn't get an image from this url")
				SSRFPwnaBot.bot.Send(msg)
			}

		}

	}

}

func findRandomImage(htm string) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htm, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
	}

	if len(out) != 0 {
		return out[r1.Intn(len(out))]
	}

	return "none"
}
