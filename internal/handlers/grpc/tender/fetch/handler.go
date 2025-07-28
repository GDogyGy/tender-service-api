package fetch

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/tender/fetch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	fetch.UnimplementedTenderServiceFetchServer
	log                 log
	useCasesTenderFetch useCasesTenderFetch
}

type log interface {
	Error(msg string, args ...any)
}

type useCasesTenderFetch interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	FetchStatus(ctx context.Context, username string, tenderId string) (model.Tender, error)
}

func NewHandler(gRPC *grpc.Server, log log, useCasesTenderFetch useCasesTenderFetch) {
	fetch.RegisterTenderServiceFetchServer(gRPC, &Handler{log: log, useCasesTenderFetch: useCasesTenderFetch})
}

func (h *Handler) FetchList(
	ctx context.Context,
	req *fetch.RequestFetchListV1,
) (*fetch.ResponseTendersV1, error) {
	const op = "handlers.grpc.tender.fetch.FetchList"

	categories, err := h.useCasesTenderFetch.FetchList(ctx, req.ServiceType)
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	var tendersTransport []transport.Tender
	for _, v := range categories {
		tendersTransport = append(tendersTransport, convert.TenderModelToTransport(v))
	}

	if len(tendersTransport) == 0 {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusNotFound").Error())
		return nil, status.Error(codes.NotFound, "StatusNotFound")
	}

	response := &fetch.ResponseTendersV1{
		Tenders: make([]*fetch.TenderV1, 0, len(tendersTransport)),
	}
	for _, v := range tendersTransport {
		response.Tenders = append(response.Tenders,
			&fetch.TenderV1{
				Id:          v.Id,
				Name:        v.Name,
				Description: v.Description,
				ServiceType: v.ServiceType,
				Status:      string(v.Status),
				Version:     int32(v.Version),
				Responsible: v.Responsible,
			})
	}

	return response, nil
}

func (h *Handler) FetchListByUser(
	ctx context.Context,
	req *fetch.RequestFetchListByUserV1,
) (*fetch.ResponseTendersV1, error) {
	const op = "handlers.grpc.tender.fetch.FetchListByUser"

	if req.Username == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	categories, err := h.useCasesTenderFetch.FetchListByUser(ctx, req.Username)
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	var tendersTransport []transport.Tender
	for _, v := range categories {
		tendersTransport = append(tendersTransport, convert.TenderModelToTransport(v))
	}

	if len(tendersTransport) == 0 {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusNotFound").Error())
		return nil, status.Error(codes.NotFound, "StatusNotFound")
	}

	response := &fetch.ResponseTendersV1{
		Tenders: make([]*fetch.TenderV1, 0, len(tendersTransport)),
	}
	for _, v := range tendersTransport {
		response.Tenders = append(response.Tenders,
			&fetch.TenderV1{
				Id:          v.Id,
				Name:        v.Name,
				Description: v.Description,
				ServiceType: v.ServiceType,
				Status:      string(v.Status),
				Version:     int32(v.Version),
				Responsible: v.Responsible,
			})
	}

	return response, nil
}

func (h *Handler) FetchStatus(
	ctx context.Context,
	req *fetch.RequestFetchStatusV1,
) (*fetch.ResponseFetchStatusV1, error) {
	const op = "handlers.grpc.tender.fetch.FetchStatus"

	if req.GetUsername() == "" || req.GetTenderId() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "StatusBadRequest").Error())
		return nil, status.Error(codes.InvalidArgument, "username is invalid")
	}

	tender, err := h.useCasesTenderFetch.FetchStatus(ctx, req.GetUsername(), req.GetTenderId())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %s", op, model.NotFindResponsible).Error())
		return nil, status.Error(codes.PermissionDenied, model.NotFindResponsible.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %s", op, model.NotFound).Error())
		return nil, status.Error(codes.PermissionDenied, model.NotFound.Error())
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	resp := convert.TenderModelToTransport(tender)

	return &fetch.ResponseFetchStatusV1{Status: string(resp.Status)}, nil
}
