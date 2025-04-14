package fetch

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/bids/fetch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	fetch.UnimplementedBidsServiceFetchServer
	log              log
	bidsFetch        useCaseBidsFetch
	bidFeedbackFetch useCaseBidFeedbackFetch
}

type log interface {
	Error(msg string, args ...any)
}

type useCaseBidsFetch interface {
	FetchListByTender(ctx context.Context, username string, tenderId string) ([]model.Bids, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Bids, error)
	FetchStatus(ctx context.Context, username string, bidsId string) (model.Bids, error)
}

type useCaseBidFeedbackFetch interface {
	FetchReviews(ctx context.Context, username string, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error)
}

func NewHandler(gRPC *grpc.Server, log log, bidsFetch useCaseBidsFetch, bidFeedbackFetch useCaseBidFeedbackFetch) {
	fetch.RegisterBidsServiceFetchServer(gRPC, &Handler{log: log, bidsFetch: bidsFetch, bidFeedbackFetch: bidFeedbackFetch})
}

func (h *Handler) FetchListByTender(
	ctx context.Context,
	req *fetch.BidsRequestFetchListV1,
) (*fetch.ResponseBidsV1, error) {
	const op = "handlers.grpc.bids.fetch.FetchListByTender"

	if req.GetUsername() == "" || req.GetTenderID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	resp, err := h.bidsFetch.FetchListByTender(ctx, req.GetUsername(), req.GetTenderID())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, model.NotFindResponsible.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	if len(resp) == 0 {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusNoContent")
	}

	var bidsTransport []transport.Bid
	for _, v := range resp {
		bidsTransport = append(bidsTransport, convert.BidsModelToTransport(v))
	}

	response := &fetch.ResponseBidsV1{
		Bids: make([]*fetch.BidV1, 0, len(bidsTransport)),
	}
	for _, v := range bidsTransport {
		response.Bids = append(response.Bids,
			&fetch.BidV1{
				Id:          v.Id,
				Name:        v.Name,
				Description: v.Description,
				Status:      string(v.Status),
				TenderId:    v.TenderId,
				Version:     int32(v.Version),
				Responsible: v.Responsible,
			})
	}

	return response, nil
}

func (h *Handler) FetchListByUser(
	ctx context.Context,
	req *fetch.BidsRequestFetchListByUserV1,
) (*fetch.ResponseBidsV1, error) {
	const op = "handlers.grpc.bids.fetch.FetchListByUser"

	if req.GetUsername() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	resp, err := h.bidsFetch.FetchListByUser(ctx, req.GetUsername())
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	if len(resp) == 0 {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusNoContent")
	}

	var bidsTransport []transport.Bid
	for _, v := range resp {
		bidsTransport = append(bidsTransport, convert.BidsModelToTransport(v))
	}

	response := &fetch.ResponseBidsV1{
		Bids: make([]*fetch.BidV1, 0, len(bidsTransport)),
	}
	for _, v := range bidsTransport {
		response.Bids = append(response.Bids,
			&fetch.BidV1{
				Id:          v.Id,
				Name:        v.Name,
				Description: v.Description,
				Status:      string(v.Status),
				TenderId:    v.TenderId,
				Version:     int32(v.Version),
				Responsible: v.Responsible,
			})
	}

	return response, nil
}

func (h *Handler) FetchStatus(
	ctx context.Context,
	req *fetch.BidsRequestFetchStatusV1,
) (*fetch.BidsResponseFetchStatusV1, error) {
	const op = "handlers.grpc.bids.fetch.FetchStatus"

	if req.GetUsername() == "" || req.GetBidID() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	bids, err := h.bidsFetch.FetchStatus(ctx, req.GetUsername(), req.GetBidID())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, model.NotFindResponsible.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	resp := convert.BidsModelToTransport(bids)
	return &fetch.BidsResponseFetchStatusV1{Status: string(resp.Status)}, nil
}

func (h *Handler) FetchReviews(
	ctx context.Context,
	req *fetch.BidsRequestFetchReviewsV1,
) (*fetch.BidsResponseFetchReviewsV1, error) {
	const op = "handlers.grpc.bids.fetch.FetchReviews"

	if req.GetUsername() == "" || req.GetTenderId() == "" || req.GetOrganizationId() == "" || req.GetAuthorUsername() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	feedback, err := h.bidFeedbackFetch.FetchReviews(ctx, req.GetUsername(), req.GetTenderId(), req.GetAuthorUsername(), req.OrganizationId)
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, model.NotFindResponsible.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	var feedbackTransport []transport.Feedback
	for _, v := range feedback {
		feedbackTransport = append(feedbackTransport, convert.BidFeedbackModelToTransport(v))
	}

	response := &fetch.BidsResponseFetchReviewsV1{
		Feedback: make([]*fetch.FeedbackV1, 0, len(feedbackTransport)),
	}
	for _, v := range feedbackTransport {
		response.Feedback = append(response.Feedback,
			&fetch.FeedbackV1{
				Id:          v.Id,
				BidId:       v.BidId,
				Description: v.Description,
				Responsible: v.Responsible,
				CreatedAt:   v.CreatedAt,
			})
	}

	return response, nil
}
