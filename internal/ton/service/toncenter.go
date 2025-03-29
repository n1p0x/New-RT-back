package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"

	giftModel "roulette/internal/gift/model"
	"roulette/internal/models"
	"roulette/internal/ton/model"
)

func (s *service) GetTonTransfers(ctx context.Context, start *int64) ([]*model.Message, error) {
	baseURL := "https://toncenter.com/api/v3"
	if s.IsTestnet {
		baseURL = "https://testnet.toncenter.com/api/v3"
	}
	url := baseURL + "/messages"

	apiKey := s.TonCenterApiKey
	if s.IsTestnet {
		apiKey = s.TonCenterApiKeyTestnet
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for ton transfer: %v", err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	q := req.URL.Query()
	q.Add("destination", s.AdminWallet)
	q.Add("exclude_externals", "true")
	if start != nil {
		q.Add("start_utime", strconv.FormatUint(uint64(*start), 10))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for ton transfers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get ton transfers: %d", resp.StatusCode)
	}

	type Response struct {
		Messages []*model.Message `json:"messages"`
	}
	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for ton transfer: %v", err)
	}

	return result.Messages, nil
}

func (s *service) GetNftTransfer(ctx context.Context, itemAddress string) (*model.NftTransfer, error) {
	baseURL := "https://toncenter.com/api/v3"
	if s.IsTestnet {
		baseURL = "https://testnet.toncenter.com/api/v3"
	}
	url := baseURL + "/nft/transfers"

	apiKey := s.TonCenterApiKey
	if s.IsTestnet {
		apiKey = s.TonCenterApiKeyTestnet
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for gift transfer: %v", err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	q := req.URL.Query()
	q.Add("owner_address", s.AdminWallet)
	q.Add("direction", "in")
	q.Add("item_address", itemAddress)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for gift transfers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gift transfers response code: %d", resp.StatusCode)
	}

	type Response struct {
		NftTransfers []*model.NftTransfer `json:"nft_transfers"`
	}
	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for gift transfer: %v", err)
	}

	//if len(result.NftTransfers) != 1 {
	//	return nil, fmt.Errorf("fetched more than one transfer")
	//}

	return result.NftTransfers[0], nil
}

func (s *service) GetNftTransfers(ctx context.Context, start *int64) ([]*model.NftTransfer, error) {
	baseURL := "https://toncenter.com/api/v3"
	if s.IsTestnet {
		baseURL = "https://testnet.toncenter.com/api/v3"
	}
	url := baseURL + "/nft/transfers"

	apiKey := s.TonCenterApiKey
	if s.IsTestnet {
		apiKey = s.TonCenterApiKeyTestnet
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for nft transfer: %v", err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	q := req.URL.Query()
	q.Add("owner_address", s.AdminWallet)
	q.Add("direction", "in")
	if start != nil {
		q.Add("start_utime", strconv.FormatUint(uint64(*start), 10))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for nft transfers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nft transfers response code: %d", resp.StatusCode)
	}

	type Response struct {
		NftTransfers []*model.NftTransfer `json:"nft_transfers"`
	}
	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for nft transfer: %v", err)
	}

	return result.NftTransfers, nil
}

func (s *service) GetWalletNfts(ctx context.Context, wallet string, collections []*giftModel.Collection) ([]*models.Nft, error) {
	baseURL := "https://toncenter.com/api/v3"
	if s.IsTestnet {
		baseURL = "https://testnet.toncenter.com/api/v3"
	}
	url := baseURL + "/nft/items"

	apiKey := s.TonCenterApiKey
	if s.IsTestnet {
		apiKey = s.TonCenterApiKeyTestnet
	}

	var (
		wg       sync.WaitGroup
		nftsCh   = make(chan []*models.Nft, len(collections))
		errCh    = make(chan error, len(collections))
	)

	for collectionID, collection := range collections {
		if collection.Address == nil {
			continue
		}

		wg.Add(1)
		go func(collectionAddr string) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				errCh <- fmt.Errorf("failed to create request for %s: %v", collectionAddr, err)
				nftsCh <- nil
				return
			}

			req.Header.Set("X-Api-Key", apiKey)

			q := req.URL.Query()
			q.Add("owner_address", wallet)
			q.Add("collection_address", collectionAddr)
			req.URL.RawQuery = q.Encode()

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errCh <- fmt.Errorf("failed to execute request for %s: %v", collectionAddr, err)
				nftsCh <- nil
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				nftsCh <- nil
				return
			}

			type Response struct {
				NftItems []*model.NftItem `json:"nft_items"`
				Metadata map[string]struct {
					IsIndexed bool `json:"is_indexed"`
					TokenInfo []struct {
						Type  string `json:"type"`
						Name  string `json:"name"`
						Extra struct {
							Lottie string `json:"lottie"`
						} `json:"extra"`
					} `json:"token_info"`
				} `json:"metadata"`
			}
			var result Response
			if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
				errCh <- fmt.Errorf("failed to decode response for %s: %v", collectionAddr, err)
				nftsCh <- nil
				return
			}

			if len(result.NftItems) == 0 {
				nftsCh <- nil
				return
			}

			nfts, err := s.fetchMetadata(ctx, collectionID, result.NftItems...)

			//var nfts []*models.Nft
			//for _, item := range result.NftItems {
			//	meta, ok := result.Metadata[item.Address]
			//	if ok != true {
			//		continue
			//	}
			//
			//	if !meta.IsIndexed {
			//
			//	}
			//}

			//for addr, item := range result.Metadata {
			//	if item.IsIndexed && len(item.TokenInfo) == 1 && item.TokenInfo[0].Type != "nft_collections" {
			//		res := strings.Split(item.TokenInfo[0].Name, "#")
			//		if len(res) == 2 {
			//			collectibleID, _ := strconv.ParseUint(res[1], 10, 64)
			//			nft := &models.Nft{
			//				Title:         res[0],
			//				CollectibleID: collectibleID,
			//				Address:       addr,
			//				LottieUrl:     item.TokenInfo[0].Extra.Lottie,
			//				CollectionID:  collectionId,
			//			}
			//			nfts = append(nfts, nft)
			//		}
			//	}
			//}

			nftsCh <- nfts
		}(*collection.Address)
	}

	go func() {
		wg.Wait()
		close(nftsCh)
		close(errCh)
	}()

	var allNfts []*models.Nft
	var errs []error

	for nfts := range nftsCh {
		if nfts != nil {
			allNfts = slices.Concat(allNfts, nfts)
		}
	}

	if len(errs) > 0 {
		return allNfts, fmt.Errorf("errors: %v", errs)
	}
	if len(allNfts) == 0 {
		return nil, ErrNftsNotFound
	}

	return allNfts, nil
}

func (s *service) GetNft(ctx context.Context, address string) (*models.Nft, error) {
	baseURL := "https://toncenter.com/api/v3"
	if s.IsTestnet {
		baseURL = "https://testnet.toncenter.com/api/v3"
	}
	url := baseURL + "/nft/items"

	apiKey := s.TonCenterApiKey
	if s.IsTestnet {
		apiKey = s.TonCenterApiKeyTestnet
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for nft item: %v", err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	q := req.URL.Query()
	q.Add("owner_address", s.AdminWallet)
	q.Add("direction", "in")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for nft item: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nft item response code: %d", resp.StatusCode)
	}

	type Response struct {
		NftItems []struct {
			Address string `json:"address"`
		} `json:"nft_items"`
		Metadata map[string]struct {
			TokenInfo []struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Extra struct {
					Lottie string `json:"lottie"`
				} `json:"extra"`
			} `json:"token_info"`
		} `json:"metadata"`
	}
	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for nft item: %v", err)
	}

	if len(result.NftItems) == 0 {
		return nil, nil
	}

	item, ok := result.Metadata[address]
	if ok != true {
		return nil, nil
	}
	res := strings.Split(item.TokenInfo[0].Name, " #")
	collectibleID, _ := strconv.Atoi(res[1])

	nft := &models.Nft{
		Name:          res[0],
		CollectibleID: uint64(collectibleID),
		Address:       address,
		LottieUrl:     item.TokenInfo[0].Extra.Lottie,
	}

	return nft, nil
}

func (s *service) fetchMetadata(ctx context.Context, collectionID int, items ...*model.NftItem) ([]*models.Nft, error) {
	var wg sync.WaitGroup
	metaCh := make(chan *models.Nft, len(items))
	errCh := make(chan error, len(items))

	for _, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, item.Content.Uri, nil)
			if err != nil {
				metaCh <- nil
				errCh <- fmt.Errorf("failed to create request for meta %s: %v", item.Content.Uri, err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				metaCh <- nil
				errCh <- fmt.Errorf("failed to fetch meta %s: %v", item.Content.Uri, err)
				return
			}
			defer resp.Body.Close()

			type Response struct {
				Name   string `json:"name"`
				Lottie string `json:"lottie"`
			}
			var result Response

			if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
				metaCh <- nil
				errCh <- fmt.Errorf("failed to decode response for meta %s: %v", item.Content.Uri, err)
				return
			}

			res := strings.Split(result.Name, "#")
			collectibleID, _ := strconv.ParseUint(res[1], 10, 64)
			nft := &models.Nft{
				Name:          res[0],
				CollectibleID: collectibleID,
				Address:       item.Address,
				LottieUrl:     result.Lottie,
				CollectionID:  collectionID,
			}

			metaCh <- nft
		}()
	}

	go func() {
		wg.Wait()
		close(metaCh)
		close(errCh)
	}()

	var allMetas []*models.Nft

	for meta := range metaCh {
		if meta != nil {
			allMetas = append(allMetas, meta)
		}
	}

	if len(allMetas) == 0 {
		return nil, nil
	}

	return allMetas, nil
}
