package report_test

import (
	"context"
	"errors"
	"testing"

	"github.com/YAMI4YOU/balance-service/internal/models"
	"github.com/YAMI4YOU/balance-service/internal/service/report"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GenerateReport(t *testing.T) {

	tests := []struct {
		name         string
		prepare      func(repo *Mockrepository, file *Mockfile)
		model        models.Report
		expectations func(t assert.TestingT, got string, err error)
	}{
		{
			name: "failed to fetch Report",
			prepare: func(repo *Mockrepository, file *Mockfile) {
				repo.EXPECT().FetchReport(gomock.Any(), models.Report{
					Year:  2025,
					Month: 9,
				}).Return(nil, errors.New("failed to fetch Report"))
			},
			model: models.Report{
				Year:  2025,
				Month: 9,
			},
			expectations: func(t assert.TestingT, _ string, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "failed to save Report",
			prepare: func(repo *Mockrepository, file *Mockfile) {
				repo.EXPECT().FetchReport(gomock.Any(), models.Report{
					Year:  2025,
					Month: 10,
				}).Return([]models.ReportSummary{{ServiceID: 1, Amount: models.NewMoneyFromKopecks(100)}}, nil)

				file.EXPECT().Save(models.Report{
					Year:  2025,
					Month: 10,
				}, gomock.Any()).Return("", errors.New("save failed"))
			},
			model: models.Report{
				Year:  2025,
				Month: 10,
			},
			expectations: func(t assert.TestingT, got string, err error) {
				assert.Error(t, err)
				assert.Empty(t, got)
			},
		},
		{
			name: "successful report generation with empty data",
			prepare: func(repo *Mockrepository, file *Mockfile) {
				repo.EXPECT().FetchReport(gomock.Any(), models.Report{
					Year:  2025,
					Month: 11,
				}).Return([]models.ReportSummary{}, nil)

				file.EXPECT().Save(models.Report{
					Year:  2025,
					Month: 11,
				}, gomock.Any()).Return("file-name-2", nil)
			},
			model: models.Report{
				Year:  2025,
				Month: 11,
			},
			expectations: func(t assert.TestingT, got string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `file-name-2`, got)
			},
		},
		{
			name: "successful report generation with data",
			prepare: func(repo *Mockrepository, file *Mockfile) {
				repo.EXPECT().FetchReport(gomock.Any(), models.Report{
					Year:  2025,
					Month: 12,
				}).Return([]models.ReportSummary{
					{ServiceID: 1, Amount: models.NewMoneyFromKopecks(250)},
					{ServiceID: 2, Amount: models.NewMoneyFromKopecks(500)},
				}, nil)

				file.EXPECT().Save(models.Report{
					Year:  2025,
					Month: 12,
				}, gomock.Any()).Return("file-name-3", nil)
			},
			model: models.Report{
				Year:  2025,
				Month: 12,
			},
			expectations: func(t assert.TestingT, got string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "file-name-3", got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockrepository(ctrl)
			mockFile := NewMockfile(ctrl)

			tt.prepare(mockRepo, mockFile)

			s := report.NewService(mockRepo, mockFile)
			got, err := s.GenerateReport(context.Background(), tt.model)
			tt.expectations(t, got, err)
		})
	}
}
