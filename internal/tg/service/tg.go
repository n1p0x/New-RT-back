package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"github.com/gotd/td/tg"
	"github.com/xssnick/tonutils-go/tlb"

	"roulette/internal/tg/model"
)

func (s *service) SendGift(ctx context.Context, userID int64, msgID int) error {
	client, err := s.GetClient(ctx)
	if err != nil {
		return err
	}

	api := client.API()

	gift := &tg.InputSavedStarGiftUser{MsgID: msgID}
	peer := &tg.InputPeerUser{UserID: userID}
	req := &tg.PaymentsTransferStarGiftRequest{
		Stargift: gift,
		ToID:     peer,
	}

	_, err = api.PaymentsTransferStarGift(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to transfer gift: %v", err)
	}

	return nil
}

func (s *service) GetClient(ctx context.Context) (*gotgproto.Client, error) {
	client, err := gotgproto.NewClient(
		s.ClientID,
		s.ClientHash,
		gotgproto.ClientTypePhone(s.ClientPhone),
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("session/gift.db")),
			Context: ctx,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	return client, nil
}

func (s *service) GetFloors(ctx context.Context, channelID int64, accessHash int64) ([]*model.CollectionFloor, error) {
	msg, err := s.getChannelMessages(ctx, channelID, accessHash)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\S+\s+(\w+\s+\w+)\s*–💎([\d.]+)`)
	matches := re.FindAllStringSubmatch(msg, -1)

	var floors []*model.CollectionFloor
	for _, match := range matches {
		if len(match) == 3 {
			floor, errNano := tlb.FromTON(strings.TrimSpace(match[2]))
			if errNano != nil {
				return nil, fmt.Errorf("failed to convert to nano: %v", errNano)
			}
			collFloor := &model.CollectionFloor{
				Name:  strings.TrimSpace(match[1]),
				Floor: floor.Nano(),
			}
			floors = append(floors, collFloor)
		}
	}

	return floors, nil
}

func (s *service) getChannelMessages(ctx context.Context, channelID int64, accessHash int64) (string, error) {
	client, err := s.GetClient(ctx)
	if err != nil {
		return "", err
	}

	api := client.API()

	req := &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  channelID,
			AccessHash: accessHash,
		},
		ID: []tg.InputMessageClass{&tg.InputMessageID{ID: 15}},
	}
	res, err := api.ChannelsGetMessages(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get channel msg: %v", err)
	}

	channelMsg, notEmpty := res.(*tg.MessagesChannelMessages).Messages[0].AsNotEmpty()
	if notEmpty == false {
		return "", fmt.Errorf("failed to parse msg")
	}

	msg := channelMsg.(*tg.Message).Message

	return msg, nil
}
