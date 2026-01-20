package dynamodb

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/common/vo"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	leadPKPrefix = "LEAD"
	leadSKPrefix = "#LEAD"
)

type Lead struct {
	PK             string    `dynamodbav:"PK,omitempty"`
	SK             string    `dynamodbav:"SK,omitempty"`
	LeadUUID       string    `dynamodbav:"lead_uuid,omitempty"`
	Name           string    `dynamodbav:"name,omitempty"`
	Email          string    `dynamodbav:"email,omitempty"`
	Phone          string    `dynamodbav:"phone,omitempty"`
	DocumentNumber string    `dynamodbav:"document_number,omitempty"`
	MotherName     string    `dynamodbav:"mother_name,omitempty"`
	BirthDate      time.Time `dynamodbav:"birth_date,omitempty"`
}

func (m *Lead) ToDomain(items []map[string]types.AttributeValue) (*model.Lead, error) {
	for _, item := range items {
		skAttr, ok := item["SK"].(*types.AttributeValueMemberS)
		if !ok {
			continue
		}

		if strings.HasPrefix(skAttr.Value, leadSKPrefix) {
			if err := attributevalue.UnmarshalMap(item, m); err != nil {
				return nil, err
			}
		}
	}

	leadName, _ := vo.NewPersonName(m.Name)
	leadEmail, _ := vo.NewEmailAddress(m.Email)
	leadPhone, _ := vo.NewPhoneNumber(m.Phone)
	leadDocumentNumber, _ := vo.NewDocumentNumber(m.DocumentNumber)
	leadMotherName, _ := vo.NewPersonName(m.MotherName)

	l := &model.Lead{
		LeadUUID:       model.LeadUUID(m.LeadUUID),
		Name:           leadName,
		Email:          leadEmail,
		Phone:          leadPhone,
		DocumentNumber: leadDocumentNumber,
		MotherName:     leadMotherName,
		BirthDate:      m.BirthDate,
	}

	return l, nil
}

func (m *Lead) ToRepo(lead model.Lead) (map[string]types.AttributeValue, error) {
	m.PK = MakeLeadPK(string(lead.LeadUUID))
	m.SK = MakeLeadSK()

	m.LeadUUID = string(lead.LeadUUID)
	m.Name = lead.Name.String()
	m.Email = lead.Email.String()
	m.Phone = lead.Phone.String()
	m.DocumentNumber = lead.DocumentNumber.String()
	m.MotherName = lead.MotherName.String()
	m.BirthDate = lead.BirthDate

	item, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func MakeLeadPK(leadID string) (pk string) {
	return fmt.Sprintf("%s#%s", leadPKPrefix, leadID)
}
func MakeLeadSK() (sk string) {
	return leadSKPrefix
}
