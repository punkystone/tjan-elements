package twitch

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
)

type SubGift struct {
	Gifter    string
	Receivers []string
	Count     *int
}
type Sub struct {
	Name string
}
type Donation struct {
	Name   string
	Amount float64
}
type Bits struct {
	Name   string
	Amount int
}

type Chat struct {
	client        *twitch.Client
	SubGiftsCache map[string]SubGift
	Events        chan any
}

var donationPattern = regexp.MustCompile(`([^\s]+) spendet €([^\s]+) tjanL Danke!`)

const reconnectInterval = 10 * time.Second

const minimumBits = 200

func InitChat() *Chat {
	client := twitch.NewAnonymousClient()
	client.TLS = true
	client.OnConnect(func() {
		log.Info().Msg("Connected to Twitch!")
	})

	client.Join("tjan")
	return &Chat{
		client:        client,
		Events:        make(chan any),
		SubGiftsCache: make(map[string]SubGift),
	}
}

func (chat *Chat) StartChat() {
	chat.client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		chat.handleUserNotice(message)
	})
	chat.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		chat.handleMessage(message)
	})
	for {
		err := chat.client.Connect()
		if err != nil {
			log.Error().Err(err)
		}
		time.Sleep(reconnectInterval)
	}
}
func (chat *Chat) handleUserNotice(message twitch.UserNoticeMessage) {
	if message.MsgID == "sub" || message.MsgID == "resub" {
		chat.Events <- Sub{
			Name: message.User.DisplayName,
		}
		return
	}

	if message.MsgID == "submysterygift" {
		count, err := strconv.Atoi(message.MsgParams["msg-param-mass-gift-count"])
		if err != nil {
			log.Error().Err(err)
			return
		}
		giftID := message.MsgParams["msg-param-community-gift-id"]
		subgift, ok := chat.SubGiftsCache[giftID]
		if !ok {
			chat.SubGiftsCache[giftID] = SubGift{
				Gifter:    message.User.DisplayName,
				Count:     &count,
				Receivers: []string{},
			}
			return
		}
		subgift.Count = &count
		chat.SubGiftsCache[giftID] = subgift
		if len(subgift.Receivers) == *subgift.Count {
			chat.Events <- subgift
			delete(chat.SubGiftsCache, giftID)
		}
	}
	if message.MsgID == "subgift" {
		giftID, ok := message.MsgParams["msg-param-community-gift-id"]
		if !ok {
			count := 1
			chat.Events <- SubGift{
				Gifter:    message.User.DisplayName,
				Count:     &count,
				Receivers: []string{message.MsgParams["msg-param-recipient-display-name"]},
			}
			return
		}
		subgift, ok := chat.SubGiftsCache[giftID]
		if !ok {
			chat.SubGiftsCache[giftID] = SubGift{
				Gifter:    message.User.DisplayName,
				Count:     nil,
				Receivers: []string{message.MsgParams["msg-param-recipient-display-name"]},
			}
			return
		}
		subgift.Receivers = append(subgift.Receivers, message.MsgParams["msg-param-recipient-display-name"])
		chat.SubGiftsCache[giftID] = subgift
		if subgift.Count == nil {
			return
		}
		if len(subgift.Receivers) == *subgift.Count {
			chat.Events <- subgift
			delete(chat.SubGiftsCache, giftID)
		}
	}
}

func (chat *Chat) handleMessage(message twitch.PrivateMessage) {
	bits := message.Bits
	if bits >= minimumBits {
		chat.Events <- Bits{
			Name:   message.User.DisplayName,
			Amount: bits,
		}
		return
	}

	if message.User.Name != "streamelements" {
		return
	}
	matches := donationPattern.FindStringSubmatch(message.Message)
	const requiredGroups = 3
	if len(matches) != requiredGroups {
		return
	}
	amount, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		log.Error().Msgf("Error parsing donation amount: %s", err)
		return
	}
	chat.Events <- Donation{
		Name:   matches[1],
		Amount: amount,
	}
}

func (chat *Chat) Parse(message string) {
	parsedMessage := twitch.ParseMessage(message)
	switch message := parsedMessage.(type) {
	case *twitch.PrivateMessage:
		chat.handleMessage(*message)
	case *twitch.UserNoticeMessage:
		chat.handleUserNotice(*message)
	}
}
