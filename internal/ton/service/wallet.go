package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"roulette/internal/utils"
)

func (s *service) SendTon(ctx context.Context, dst string, amount uint) (string, error) {
	w, _, err := s.getWallet(ctx)
	if err != nil {
		return "", err
	}

	receivers := map[string]string{
		dst: strconv.FormatUint(uint64(amount), 10),
	}

	var messages []*wallet.Message
	for addrStr, amountStr := range receivers {
		addr, errAddr := utils.GetAddress(addrStr)
		if errAddr != nil {
			continue
		}

		messages = append(messages, &wallet.Message{
			Mode: wallet.PayGasSeparately + wallet.IgnoreErrors,
			InternalMessage: &tlb.InternalMessage{
				IHRDisabled: true,
				Bounce:      addr.IsBounceable(),
				DstAddr:     addr,
				Amount:      tlb.MustFromTON(amountStr),
			},
		})
	}

	txHash, err := w.SendManyWaitTxHash(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to send tx: %v", err)
	}

	txHashBase := base64.StdEncoding.EncodeToString(txHash)
	return txHashBase, nil
}

func (s *service) SendNft(ctx context.Context, dst string, nftAddress string) (string, error) {
	nftAddr, err := utils.GetAddress(nftAddress)
	if err != nil {
		return "", err
	}
	dstAddr, err := utils.GetAddress(dst)
	if err != nil {
		return "", err
	}
	amountForward, _ := tlb.FromTON("0.15")

	w, api, err := s.getWallet(ctx)
	if err != nil {
		return "", err
	}

	nftItem := nft.NewItemClient(api, nftAddr)
	transferPayload, err := nftItem.BuildTransferPayload(dstAddr, amountForward, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create transfer payload")
	}

	msg := wallet.SimpleMessage(nftAddr, tlb.MustFromTON("0.05"), transferPayload)
	tx, _, err := w.SendWaitTransaction(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send tx: %v", err)
	}

	txHashBase := base64.StdEncoding.EncodeToString(tx.Hash)
	return txHashBase, nil
}

func (s *service) getWallet(ctx context.Context) (*wallet.Wallet, ton.APIClientWrapped, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/global.config.json"
	if s.IsTestnet {
		configUrl = "https://ton-blockchain.github.io/testnet-global.config.json"
	}

	err := client.AddConnectionsFromConfigUrl(ctx, configUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to lite: %v", err)
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	mnemonic := strings.Split(s.Mnemonic, " ")

	var w *wallet.Wallet
	if s.IsTestnet {
		w, err = wallet.FromSeed(api, mnemonic, wallet.ConfigV5R1Final{
			NetworkGlobalID: wallet.TestnetGlobalID,
		})
	} else {
		w, err = wallet.FromSeed(api, mnemonic, wallet.ConfigHighloadV3{
			MessageTTL: 60 * 5,
			MessageBuilder: func(ctx context.Context, subWalletId uint32) (id uint32, createdAt int64, err error) {
				createdAt = time.Now().Unix() - 30
				return uint32(createdAt % (1 << 23)), createdAt, nil
			},
		})
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize wallet: %v", err)
	}

	return w, api, nil
}
