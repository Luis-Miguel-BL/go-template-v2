package usecase

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/common/vo"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
)

type CreateLead struct {
	authService service.AuthService
	leadRepo    lead.LeadRepository
}

func NewCreateLead(authService service.AuthService, leadRepo lead.LeadRepository) *CreateLead {
	return &CreateLead{
		authService: authService,
		leadRepo:    leadRepo,
	}
}

func (u *CreateLead) Execute(ctx context.Context, input dto.CreateLeadInput) (output dto.CreateLeadOutput, err error) {
	var newLeadParams model.NewLeadParams

	if newLeadParams.Name, err = vo.NewPersonName(input.Name); err != nil {
		return output, err
	}

	if newLeadParams.Email, err = vo.NewEmailAddress(input.Email); err != nil {
		return output, err
	}

	if newLeadParams.Phone, err = vo.NewPhoneNumber(input.Phone); err != nil {
		return output, err
	}

	if newLeadParams.DocumentNumber, err = vo.NewDocumentNumber(input.DocumentNumber); err != nil {
		return output, err
	}

	lead := model.NewLead(newLeadParams)

	if err = u.leadRepo.Save(ctx, lead); err != nil {
		return output, err
	}
	leadID := string(lead.LeadUUID)

	accessToken, err := u.authService.RefreshToken(ctx, &service.TokenClaims{
		LeadID: &leadID,
	})

	if err != nil {
		return output, err
	}

	return dto.CreateLeadOutput{
		LeadID:      string(lead.LeadUUID),
		AccessToken: accessToken,
	}, nil
}
