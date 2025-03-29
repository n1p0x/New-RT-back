package model

type Collection struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Floor   *uint  `json:"floor"`
	ImgUrl  string `json:"imgUrl"`
}

type Message struct {
	Hash        string `json:"hash"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Value       string `json:"value"`
	MsgContent  struct {
		Body    string `json:"body"`
		Decoded struct {
			Type    string  `json:"type"`
			Comment *string `json:"comment"`
		} `json:"decoded"`
		Hash string `json:"hash"`
	} `json:"message_content"`
}

type NftTransfer struct {
	Sender         string  `json:"old_owner"`
	NftAddress     string  `json:"nft_address"`
	NftCollection  string  `json:"nft_collection"`
	ForwardPayload *string `json:"forward_payload"`
	TraceID        string  `json:"trace_id"`
}

type NftItem struct {
	Address string `json:"address"`
	Content struct {
		Uri string `json:"uri"`
	} `json:"content"`
}
