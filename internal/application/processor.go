package application

import (
	"github.com/andreposman/capital-gains/internal/domain"
	"github.com/andreposman/capital-gains/internal/infra/json"
	"log"
)

type OperationProcessor struct{}

func (op *OperationProcessor) ProcessOperations(operations []json.Operation) []domain.Tax {
	portfolio := domain.Portfolio{}
	results := make([]domain.Tax, len(operations))

	for i, operation := range operations {
		var currentTax = 0.0
		var err error

		switch operation.Operation {
		case "buy":
			portfolio.Buy(operation.Quantity, operation.UnitCost)
			currentTax = 0.0

		case "sell":
			currentTax, err = portfolio.Sell(operation.Quantity, operation.UnitCost)
			if err != nil {
				log.Fatalf("fatal error: Unexpected error during sell operation %d, assumption violated: %v", i+1, err)
			}

		default:
			log.Printf("Warning: Unknown operation type '%s' at index %d", operation.Operation, i)

		}

		results[i] = domain.Tax{Tax: currentTax}
	}
	return results
}
