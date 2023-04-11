package domain

type ReconciliationReq struct {
	AfterIndex BlockIndex `json:"after_index"`

	metaSenderPeerID string
}

func (req *ReconciliationReq) SetSender(peerID string) {
	req.metaSenderPeerID = peerID
}

func (req *ReconciliationReq) GetSender() string {
	return req.metaSenderPeerID
}

type ReconciliationResp struct {
	AfterIndex BlockIndex `json:"after_index"`
	NextBlocks []*Block   `json:"next_blocks"`

	metaSenderPeerID string
}

func (resp *ReconciliationResp) SetSender(peerID string) {
	resp.metaSenderPeerID = peerID
}

func (resp *ReconciliationResp) GetSender() string {
	return resp.metaSenderPeerID
}
