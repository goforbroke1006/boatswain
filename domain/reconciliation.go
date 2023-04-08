package domain

type ReconciliationReq struct {
	AfterIndex BlockIndex
}

type ReconciliationResp struct {
	MetaSenderPeerID string

	AfterIndex BlockIndex
	NextBlocks []*Block
}

func (resp ReconciliationResp) SetSender(peerID string) {
	resp.MetaSenderPeerID = peerID
}

func (resp ReconciliationResp) GetSender() string {
	return resp.MetaSenderPeerID
}
