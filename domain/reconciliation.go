package domain

type ReconciliationReq struct {
	AfterIndex BlockIndex `json:"after_index"`
}

type ReconciliationResp struct {
	metaSenderPeerID string

	AfterIndex BlockIndex `json:"after_index"`
	NextBlocks []*Block   `json:"next_blocks"`
}

func (resp ReconciliationResp) SetSender(peerID string) {
	resp.metaSenderPeerID = peerID
}

func (resp ReconciliationResp) GetSender() string {
	return resp.metaSenderPeerID
}
