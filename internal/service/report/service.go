package report

import (
	"context"
	"fmt"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Service struct {
	repo repository
	file file
}

func NewService(repo repository, file file) *Service {
	return &Service{
		repo: repo,
		file: file,
	}
}

func (s *Service) GenerateReport(ctx context.Context, model models.Report) (string, error) {
	summaries, err := s.repo.FetchReport(ctx, model)
	if err != nil {
		return "", fmt.Errorf("generate monthly report error: %w", err)
	}

	return s.file.Save(model, summaries)
}
