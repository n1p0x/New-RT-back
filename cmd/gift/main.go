package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"

	"roulette/internal/config"
	"roulette/internal/database"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	tgService "roulette/internal/tg/service"
)

const (
	configPath = "config/prod.yaml"
	envPath    = ".env"
	//downloadPath = "/var/www/rton/static/gift/"
	downloadPath = "/Users/gemini/dev/go/roulette/"
)

func main() {
	runGift()
}

func runGift() {
	cfg, err := config.Load(configPath, envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	service := tgService.NewService(cfg)

	client, err := service.GetClient(context.Background())
	if err != nil {
		log.Fatalf("failed to load client: %v", err)
	}

	dp := client.Dispatcher
	dp.AddHandler(handlers.NewMessage(func(m *types.Message) bool {
		switch m.Action.(type) {
		case *tg.MessageActionStarGiftUnique:
			return true
		default:
			return false
		}
	}, processGift))
	dp.AddHandler(handlers.NewMessage(func(m *types.Message) bool {
		channel, ok := m.PeerID.(*tg.PeerChannel)
		if !ok {
			return false
		}
		if channel.ChannelID != 2422226195 {
			return false
		}
		return true
	}, processGiftMonitor))

	log.Printf("started...")

	err = client.Idle()
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	log.Printf("completed...")
}

func processGift(ctx *ext.Context, update *ext.Update) error {
	m := update.EffectiveMessage

	log.Printf("processing...")

	giftAction := m.Action.(*tg.MessageActionStarGiftUnique)
	gift := giftAction.GetGift()
	if gift == nil {
		errGift := errors.New("failed to get gift")
		log.Printf(errGift.Error())
		return errGift
	}

	uniqueStarGift, ok := gift.(*tg.StarGiftUnique)
	if !ok {
		errUniqueStarGift := errors.New("failed to get unique star gift")
		log.Printf(errUniqueStarGift.Error())
		return errUniqueStarGift
	}

	cfg, err := config.Load(configPath, envPath)
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return err
	}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Printf("failed to init db: %v", err)
		return err
	}
	repo := giftRepo.NewRepo(db)
	service := giftService.NewService(repo)

	giftID := uniqueStarGift.GetID()
	if giftID == 0 {
		errGiftID := errors.New("failed to get gift id")
		log.Printf(errGiftID.Error())
		return errGiftID
	}
	msgID, ok := giftAction.GetSavedID()
	if !ok {
		errMsgID := errors.New("failed to get msg id")
		log.Printf(errMsgID.Error(), msgID)
		return errMsgID
	}
	slug := uniqueStarGift.GetSlug()
	if slug == "" {
		errSlug := errors.New("failed to get gift slug")
		log.Printf(errSlug.Error(), giftID)
		return errSlug
	}
	title := uniqueStarGift.GetTitle()
	if title == "" {
		errTitle := errors.New("failed to get gift title")
		log.Printf(errTitle.Error(), giftID)
		return errTitle
	}
	collectibleID := uniqueStarGift.GetNum()
	if collectibleID == 0 {
		errCollectibleID := errors.New("failed to get gift collectible id")
		log.Printf(errCollectibleID.Error(), giftID)
		return errCollectibleID
	}

	var document *tg.Document
	var senderID int64
	for _, attr := range uniqueStarGift.Attributes {
		switch v := attr.(type) {
		case *tg.StarGiftAttributeModel:
			document, ok = v.Document.(*tg.Document)
			if !ok {
				errDocument := errors.New("failed to get document")
				log.Printf(errDocument.Error())
				return errDocument
			}
		case *tg.StarGiftAttributeOriginalDetails:
			peer, ok := v.SenderID.(*tg.PeerUser)
			if ok {
				senderID = peer.UserID
			}
		}

		model, ok := attr.(*tg.StarGiftAttributeModel)
		if ok {
			document, ok = model.Document.(*tg.Document)
			if !ok {
				errDocument := errors.New("failed to get document")
				log.Printf(errDocument.Error())
				return errDocument
			}
		}
	}

	downloadOutputPath := filepath.Join(downloadPath, fmt.Sprintf("%s.tgs", strings.ToLower(slug)))
	mediaDocument := &tg.MessageMediaDocument{Document: document}
	_, err = ctx.DownloadMedia(
		mediaDocument,
		ext.DownloadOutputPath(downloadOutputPath),
		nil,
	)
	if err != nil {
		errDownload := fmt.Errorf("failed to download gift media: %v", err)
		log.Printf(errDownload.Error())
		return errDownload
	}

	collection, err := service.GetCollectionByName(ctx.Context, title)
	if err != nil {
		log.Printf("failed to get collection: %v", err)
		return err
	}

	lottieUrl := fmt.Sprintf("https://rouletton.ru/static/%s.tgs", strings.ToLower(slug))
	if err = service.AddUserGift(context.Background(), senderID, giftID, msgID, title, collectibleID, collection.ID, lottieUrl); err != nil {
		log.Printf("failed to add user gift <%d:%d>: %v", senderID, giftID, err)
	}

	return nil
}

func processGiftMonitor(ctx *ext.Context, update *ext.Update) error {
	m := update.EffectiveMessage
	fmt.Println(m.Text)

	return nil
}
