package application

import (
	"github.com/andreposman/capital-gains/internal/domain"
	"github.com/andreposman/capital-gains/internal/infra/json"
)

type OperationProcessor struct{}

func (op *OperationProcessor) ProcessOperations(operations []json.Operation) []domain.Tax {
	portfolio := domain.Portfolio{}
	results := make([]domain.Tax, len(operations))

	for i, op := range operations {
		switch op.Operation {
		case "buy":
			portfolio.Buy(op.Quantity, op.UnitCost)
			results[i] = domain.Tax{0.00}
		case "sell":
			tax := portfolio.Sell(op.Quantity, op.UnitCost)
			results[i] = domain.Tax{tax}

		}
	}
	return results
}
