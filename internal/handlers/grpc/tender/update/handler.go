package update

import (
	"TenderServiceApi/internal/handlers/rest/types/convert"
	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/protos/gen/tender/update"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	update.UnimplementedTenderServiceUpdateServer
	log        log
	tenderEdit useCaseTenderEdit
}

type log interface {
	Error(msg string, args ...any)
}

type useCaseTenderEdit interface {
	Edit(ctx context.Context, id string, username string, tenderNew model.Tender) (model.Tender, error)
	Rollback(ctx context.Context, id string, username string, version string) (model.Tender, error)
	Status(ctx context.Context, username string, tenderId string, status string) (model.Tender, error)
}

func NewHandler(gRPC *grpc.Server, log log, tenderEdit useCaseTenderEdit) {
	update.RegisterTenderServiceUpdateServer(gRPC, &Handler{log: log, tenderEdit: tenderEdit})
}

func (h *Handler) Edit(
	ctx context.Context,
	req *update.RequestEditV1,
) (*update.TenderV1, error) {
	const op = "handlers.grpc.tender.update.Edit"

	if req.GetTenderId() == "" || req.GetUsername() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "invalid Request").Error())
		return nil, status.Error(codes.InvalidArgument, "invalid Request")
	}

	if req.GetName() == "" && req.GetDescription() == "" && req.GetServiceType() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "invalid Request").Error())
		return nil, status.Error(codes.InvalidArgument, "invalid Request")
	}

	tenderNew := model.Tender{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ServiceType: req.GetServiceType(),
	}
	resp, err := h.tenderEdit.Edit(ctx, req.GetTenderId(), req.GetUsername(), tenderNew)
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	tender := convert.TenderModelToTransport(resp)

	return &update.TenderV1{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      string(tender.Status),
		Version:     int32(tender.Version),
		Responsible: tender.Responsible,
	}, nil
}

func (h *Handler) Rollback(
	ctx context.Context,
	req *update.RequestRollbackV1,
) (*update.TenderV1, error) {
	const op = "handlers.grpc.tender.update.Rollback"

	if req.GetTenderId() == "" && req.GetVersion() == "" && req.GetUsername() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "invalid Request").Error())
		return nil, status.Error(codes.InvalidArgument, "invalid Request")
	}

	resp, err := h.tenderEdit.Rollback(ctx, req.GetTenderId(), req.GetUsername(), req.GetVersion())
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	tender := convert.TenderModelToTransport(resp)

	return &update.TenderV1{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      string(tender.Status),
		Version:     int32(tender.Version),
		Responsible: tender.Responsible,
	}, nil
}

func (h *Handler) Status(
	ctx context.Context,
	req *update.RequestStatusV1,
) (*update.ResponseStatusV1, error) {
	const op = "handlers.grpc.tender.update.Status"

	if req.GetTenderId() == "" && req.GetUsername() == "" && req.GetStatus() == "" {
		h.log.Error(fmt.Errorf("%s: %s", op, "invalid Request").Error())
		return nil, status.Error(codes.InvalidArgument, "invalid Request")
	}

	resp, err := h.tenderEdit.Status(ctx, req.GetUsername(), req.GetTenderId(), req.GetStatus())
	if errors.Is(err, model.NotFindResponsible) {
		h.log.Error(fmt.Errorf("%s: %w", op, model.NotFindResponsible).Error())
		return nil, status.Error(codes.Internal, model.NotFindResponsible.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Error(fmt.Errorf("%s: %w", op, model.NotFound).Error())
		return nil, status.Error(codes.Internal, model.NotFound.Error())
	}
	if err != nil {
		h.log.Error(fmt.Errorf("%s: %w", op, err).Error())
		return nil, status.Error(codes.Internal, "StatusInternalServerError")
	}

	tender := convert.TenderModelToTransport(resp)

	return &update.ResponseStatusV1{
		Status: string(tender.Status),
	}, nil
}
