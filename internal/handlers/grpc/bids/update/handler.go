package update

import (
	"context"
	"errors"
	"fmt"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/bids/update"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	update.UnimplementedBidsServiceFetchServer
	log         log
	bidsEdit    useCaseBidEdit
	bidDecision useCaseBidDecision
}

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidEdit interface {
	Edit(ctx context.Context, id string, username string, bidNew model.Bids) (model.Bids, error)
	Rollback(ctx context.Context, id string, username string, version string) (model.Bids, error)
	Status(ctx context.Context, username string, bidID string, status string) (model.Bids, error)
}

type useCaseBidDecision interface {
	SubmitDecision(ctx context.Context, username string, bidID string, decision string, organizationID string) (model.Bids, error)
}

func NewHandler(gRPC *grpc.Server, log log, bidsEdit useCaseBidEdit, bidDecision useCaseBidDecision) {
	update.RegisterBidsServiceFetchServer(gRPC, &Handler{log: log, bidsEdit: bidsEdit, bidDecision: bidDecision})
}

func (h *Handler) Edit(
	ctx context.Context,
	req *update.BidsRequestEditV1,
) (*update.BidEditV1, error) {
	const op = "handlers.grpc.bids.update.Edit"

	if req.GetUsername() == "" || req.GetBidID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	editReq := transport.BidEditRequest{
		Description: req.Description,
		Name:        req.Name,
	}

	resp, err := h.bidsEdit.Edit(ctx, req.GetBidID(), req.GetUsername(), convert.BidsReqEditTransportToModel(editReq))
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	bid := convert.BidsModelToTransport(resp)

	return &update.BidEditV1{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      string(bid.Status),
		TenderId:    bid.TenderId,
		Version:     int32(bid.Version),
		Responsible: bid.Responsible,
	}, nil
}

func (h *Handler) Rollback(
	ctx context.Context,
	req *update.BidsRequestRollbackV1,
) (*update.BidEditV1, error) {
	const op = "handlers.grpc.bids.update.Rollback"

	if req.GetUsername() == "" || req.GetVersion() == "" || req.GetBidID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	resp, err := h.bidsEdit.Rollback(ctx, req.GetBidID(), req.GetUsername(), req.GetVersion())
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	bid := convert.BidsModelToTransport(resp)

	return &update.BidEditV1{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      string(bid.Status),
		TenderId:    bid.TenderId,
		Version:     int32(bid.Version),
		Responsible: bid.Responsible,
	}, nil
}

func (h *Handler) Status(
	ctx context.Context,
	req *update.BidsRequestStatusV1,
) (*update.BidEditV1, error) {
	const op = "handlers.grpc.bids.update.Status"

	if req.GetUsername() == "" || req.GetStatus() == "" || req.GetBidID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "Arguments is invalid")
	}

	if transport.PutStatusBidParamsStatus(req.GetStatus()) != transport.PutStatusBidParamsStatusCREATED && transport.PutStatusBidParamsStatus(req.GetStatus()) != transport.PutStatusBidParamsStatusPUBLISHED {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "status is invalid")
	}

	resp, err := h.bidsEdit.Status(ctx, req.GetUsername(), req.GetBidID(), req.GetStatus())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, model.NotFindResponsible.Error())
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	bid := convert.BidsModelToTransport(resp)

	return &update.BidEditV1{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      string(bid.Status),
		TenderId:    bid.TenderId,
		Version:     int32(bid.Version),
		Responsible: bid.Responsible,
	}, nil
}

func (h *Handler) SubmitDecision(
	ctx context.Context,
	req *update.BidsRequestSubmitDecisionV1,
) (*update.BidEditV1, error) {
	const op = "handlers.grpc.bids.update.SubmitDecision"

	if req.GetUsername() == "" || req.GetBidID() == "" || req.GetDecision() == "" || req.GetOrganizationID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "Arguments is invalid")
	}

	if transport.SubmitDecisionBidParamsDecision(req.GetDecision()) != transport.APPROVED && transport.SubmitDecisionBidParamsDecision(req.GetDecision()) != transport.REJECTED {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "decision is invalid")
	}

	resp, err := h.bidDecision.SubmitDecision(ctx, req.GetUsername(), req.GetBidID(), req.GetDecision(), req.GetOrganizationID())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, model.NotFindResponsible.Error())
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	bid := convert.BidsModelToTransport(resp)

	return &update.BidEditV1{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      string(bid.Status),
		TenderId:    bid.TenderId,
		Version:     int32(bid.Version),
		Responsible: bid.Responsible,
	}, nil
}
