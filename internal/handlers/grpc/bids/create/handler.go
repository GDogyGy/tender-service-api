package create

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/bids/create"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	create.UnimplementedBidsServiceCreateServer
	log               log
	bidsCreate        useCaseBidsCreate
	bidFeedbackCreate useCaseBidFeedbackCreate
}

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidsCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Bids) (model.Bids, error)
}

type useCaseBidFeedbackCreate interface {
	Create(ctx context.Context, creatorUsername string, tenderID string, saveModel model.BidFeedback) (model.BidFeedback, error)
}

func NewHandler(gRPC *grpc.Server, log log, bidsCreate useCaseBidsCreate, bidFeedbackCreate useCaseBidFeedbackCreate) {
	create.RegisterBidsServiceCreateServer(gRPC, &Handler{log: log, bidsCreate: bidsCreate, bidFeedbackCreate: bidFeedbackCreate})
}

func (h *Handler) Create(
	ctx context.Context,
	req *create.BidRequestCreateV1,
) (*create.BidV1, error) {
	const op = "handlers.grpc.bids.create.Create"

	if req.GetCreatorUsername() == "" || req.GetOrganizationId() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	bidRequest := transport.BidCreateRequest{
		CreatorUsername: req.CreatorUsername,
		Description:     req.Description,
		Name:            req.Name,
		OrganizationId:  req.OrganizationId,
		Status:          transport.BidCreateRequestStatus(req.Status),
		TenderId:        req.TenderId,
	}

	if bidRequest.Status != transport.BidCreateRequestStatusCREATED {
		h.log.Error(model.BadStatus.Error())
		return nil, status.Error(codes.InvalidArgument, model.BadStatus.Error())
	}

	resp, err := h.bidsCreate.Create(ctx, bidRequest.CreatorUsername, bidRequest.OrganizationId, convert.BidsReqCreateTransportToModel(bidRequest))
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	bid := convert.BidsModelToTransport(resp)

	return &create.BidV1{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      string(bid.Status),
		TenderId:    bid.TenderId,
		Version:     int32(bid.Version),
		Responsible: bid.Responsible,
	}, nil
}

func (h *Handler) Feedback(
	ctx context.Context,
	req *create.BidRequestFeedbackV1,
) (*create.BidFeedbackV1, error) {
	const op = "handlers.grpc.bids.create.Feedback"

	if req.GetUsername() == "" || req.GetBidID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	feedback := transport.Feedback{
		Description: req.GetDescription(),
	}

	resp, err := h.bidFeedbackCreate.Create(ctx, req.GetUsername(), req.GetBidID(), convert.BidFeedbackTransportToModel(feedback))
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	feedback = convert.BidFeedbackModelToTransport(resp)
	return &create.BidFeedbackV1{
		Id:          feedback.Id,
		BidID:       feedback.BidId,
		Description: feedback.Description,
		Responsible: feedback.Responsible,
		CreatedAt:   feedback.CreatedAt,
	}, nil
}
