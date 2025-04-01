package main

import (
	"context"
	"fmt"
	"log"

	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"

	"roulette/internal/config"
	tgService "roulette/internal/tg/service"
)

var configPath = "config/local.yaml"

func runGift() {
	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	service := tgService.NewService(conf)

	client, err := service.GetClient(context.Background())
	if err != nil {
		log.Fatalf("failed to load client: %v", err)
	}

	dp := client.Dispatcher

	dp.AddHandler(handlers.NewMessage(func(m *types.Message) bool {
		switch m.Action.(type) {
		case *tg.MessageActionStarGift:
			return true
		default:
			return false
		}
	}, processGift))

	err = client.Idle()
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}
}

func processGift(ctx *ext.Context, update *ext.Update) error {
	m := update.EffectiveMessage

	giftAction := m.Action.(*tg.MessageActionStarGift)
	if gift, ok := giftAction.Gift.(*tg.StarGift); ok == true {
		document := gift.Sticker.(*tg.Document)
		mediaDocument := &tg.MessageMediaDocument{Document: document}

		//var mediaDocument *tg.MessageMediaDocument
		fmt.Println(document.GetFileReference())

		res, err := ctx.DownloadMedia(
			mediaDocument,
			ext.DownloadOutputPath(""),
			nil,
		)
		if err != nil {
		}

		fmt.Println(res)
	}

	return nil
}
