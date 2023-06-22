package impl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/component/node/api/spec"
)

func NewHandlers(
	p2pHost host.Host,
	txCache domain.TransactionCache,
	txSpreadSvc domain.TransactionSpreadInfoService,
) spec.ServerInterface {
	return &handlers{
		p2pHost: p2pHost,

		txCache:     txCache,
		txSpreadSvc: txSpreadSvc,
	}
}

var _ spec.ServerInterface = (*handlers)(nil)

type handlers struct {
	p2pHost host.Host

	txCache     domain.TransactionCache
	txSpreadSvc domain.TransactionSpreadInfoService
}

func (h handlers) GetPing(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}

func (h handlers) GetAddresses(ctx echo.Context) error {
	addrs := h.p2pHost.Addrs()

	respArr := make([]spec.PeerMultiAddress, 0, len(addrs))

	for _, addr := range addrs {
		ma := fmt.Sprintf("%s/p2p/%s", addr, h.p2pHost.ID().String())
		respArr = append(respArr, ma)
	}

	return ctx.JSON(http.StatusOK, respArr)
}

func (h handlers) GetAddressesPeerId(ctx echo.Context, peerId spec.PeerHostID) error {
	id, _ := peer.Decode(peerId)
	peerInfo := h.p2pHost.Peerstore().PeerInfo(id)

	addrs := peerInfo.Addrs

	respArr := make([]spec.PeerMultiAddress, 0, len(addrs))

	for _, addr := range addrs {
		ma := fmt.Sprintf("%s/p2p/%s", addr, peerInfo.ID.String())
		respArr = append(respArr, ma)
	}

	return ctx.JSON(http.StatusOK, respArr)
}

func (h handlers) GetPeers(ctx echo.Context) error {
	peers := h.p2pHost.Peerstore().Peers()

	respArr := make([]spec.PeerHostID, 0, len(peers))

	for _, p := range peers {
		respArr = append(respArr, p.String())
	}

	return ctx.JSON(http.StatusOK, respArr)
}

func (h handlers) PostTransaction(ctx echo.Context) error {
	var tx domain.Transaction // TODO: parse data

	h.txCache.Append(tx)

	if spreadErr := h.txSpreadSvc.Spread(tx); spreadErr != nil {
		return spreadErr
	}

	return ctx.JSON(http.StatusCreated, nil)
}
