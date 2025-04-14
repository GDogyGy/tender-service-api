package create

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/tender/create"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	create.UnimplementedTenderServiceCreateServer
	log                  log
	useCasesTenderCreate useCasesTenderCreate
}

type log interface {
	Error(msg string, args ...any)
}

type useCasesTenderCreate interface {
	Create(ctx context.Context, creatorUsername string, organizationId string, saveModel model.Tender) (model.Tender, error)
}

func NewHandler(gRPC *grpc.Server, log log, useCasesTenderCreate useCasesTenderCreate) {
	create.RegisterTenderServiceCreateServer(gRPC, &Handler{log: log, useCasesTenderCreate: useCasesTenderCreate})
}

func (h *Handler) Create(
	ctx context.Context,
	req *create.RequestCreateV1,
) (*create.TenderV1, error) {
	const op = "handlers.grpc.tender.create.Create"

	tenderRequest := transport.TenderCreateRequest{
		CreatorUsername: req.CreatorUsername,
		Description:     req.Description,
		Name:            req.Name,
		OrganizationId:  req.OrganizationId,
		ServiceType:     req.ServiceType,
		Status:          transport.TenderCreateRequestStatus(req.Status),
	}

	if tenderRequest.Status != transport.TenderCreateRequestStatusCREATED {
		h.log.Error(model.BadStatus.Error())
		return nil, status.Error(codes.InvalidArgument, model.BadStatus.Error())
	}

	resp, err := h.useCasesTenderCreate.Create(ctx, tenderRequest.CreatorUsername, tenderRequest.OrganizationId, convert.TenderReqCreateTransportToModel(tenderRequest))
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.PermissionDenied, "StatusForbidden")
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	tender := convert.TenderModelToTransport(resp)

	return &create.TenderV1{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      string(tender.Status),
		Version:     int32(tender.Version),
		Responsible: tender.Responsible,
	}, nil
}
