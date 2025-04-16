package application

import (
	"github.com/andreposman/capital-gains/internal/domain"
	"github.com/andreposman/capital-gains/internal/infra/json"
	"reflect"
	"testing"
)

func op(opType string, cost float64, qty int) json.Operation {
	return json.Operation{Operation: opType, UnitCost: cost, Quantity: qty}
}

func taxResult(taxAmount float64) domain.Tax {
	return domain.Tax{Tax: taxAmount}
}

func TestOperationProcessor_ProcessOperations_Case1(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 100},
	// {"operation":"sell", "unit-cost":15.00, "quantity": 50},
	// {"operation":"sell", "unit-cost":15.00, "quantity": 50}]
	// valor venda <= 20k
	operations := []json.Operation{
		op("buy", 10.00, 100),
		op("sell", 15.00, 50), // Total = 750 <= 20k -> Exempt, Profit = 50*(15-10)=250. Loss becomes max(0, 0-250)=0
		op("sell", 15.00, 50), // Total = 750 <= 20k -> Exempt, Profit = 50*(15-10)=250. Loss becomes max(0, 0-250)=0
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(0.0),
		taxResult(0.0),
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 1 failed: Expected %v, got %v", expected, result)
	}
}

func TestOperationProcessor_ProcessOperations_Case2(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
	// {"operation":"sell", "unit-cost":20.00, "quantity": 5000}, -> Taxable Profit 50k -> Tax 10k
	// {"operation":"sell", "unit-cost":5.00, "quantity": 5000}] -> Taxable Loss 25k -> Acc Loss 25k
	operations := []json.Operation{
		op("buy", 10.00, 10000),
		op("sell", 20.00, 5000),
		op("sell", 5.00, 5000),
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(10000.0),
		taxResult(0.0),
		// Final state should have accLoss = 25000
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 2 failed: Expected %v, got %v", expected, result)
	}
}

func TestOperationProcessor_ProcessOperations_Case3(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
	// {"operation":"sell", "unit-cost":5.00, "quantity": 5000}, -> Taxable Loss 25k -> Acc Loss 25k
	// {"operation":"sell", "unit-cost":20.00, "quantity": 3000}] -> Taxable Profit 3k*(20-10)=30k. Net=30k-25k=5k. Tax=1k
	operations := []json.Operation{
		op("buy", 10.00, 10000),
		op("sell", 5.00, 5000),
		op("sell", 20.00, 3000),
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(0.0),
		taxResult(1000.0),
		// Final state should have accLoss = 0
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 3 failed: Expected %v, got %v", expected, result)
	}
}

func TestOperationProcessor_ProcessOperations_Case4(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 10000}, -> WAC 10
	// {"operation":"buy", "unit-cost":25.00, "quantity": 5000}, -> WAC (100k+125k)/15k = 225k/15k = 15
	// {"operation":"sell", "unit-cost":15.00, "quantity": 10000}] -> Taxable Breakeven -> Tax 0
	operations := []json.Operation{
		op("buy", 10.00, 10000),
		op("buy", 25.00, 5000),
		op("sell", 15.00, 10000),
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(0.0),
		taxResult(0.0),
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 4 failed: Expected %v, got %v", expected, result)
	}
}

func TestOperationProcessor_ProcessOperations_Case5(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 10000}, -> WAC 10
	// {"operation":"buy", "unit-cost":25.00, "quantity": 5000}, -> WAC 15
	// {"operation":"sell", "unit-cost":15.00, "quantity": 10000}, -> Taxable Break even -> Tax 0
	// {"operation":"sell", "unit-cost":25.00, "quantity": 5000}] -> Taxable Profit 5k*(25-15)=50k -> Tax 10k
	operations := []json.Operation{
		op("buy", 10.00, 10000),
		op("buy", 25.00, 5000),
		op("sell", 15.00, 10000),
		op("sell", 25.00, 5000),
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(0.0),
		taxResult(0.0),
		taxResult(10000.0),
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 5 failed: Expected %v, got %v", expected, result)
	}
}

// Case 6 involves exempt loss followed by profit reduced by that loss
func TestOperationProcessor_ProcessOperations_Case6(t *testing.T) {
	// [{"operation":"buy", "unit-cost":10.00, "quantity": 10000}, -> WAC 10
	// {"operation":"sell", "unit-cost":2.00, "quantity": 5000}, -> Exempt Loss 5k*(2-10)=-40k. Total=10k<=20k. Acc Loss=40k
	// {"operation":"sell", "unit-cost":20.00, "quantity": 2000}, -> Taxable Profit 2k*(20-10)=20k. Net=max(0,20k-40k)=0. Tax=0. Rem Loss=20k
	// {"operation":"sell", "unit-cost":20.00, "quantity": 2000}, -> Taxable Profit 2k*(20-10)=20k. Net=max(0,20k-20k)=0. Tax=0. Rem Loss=0
	// {"operation":"sell", "unit-cost":25.00, "quantity": 1000}] -> Taxable Profit 1k*(25-10)=15k. Net=max(0,15k-0)=15k. Tax=3k
	operations := []json.Operation{
		op("buy", 10.00, 10000),
		op("sell", 2.00, 5000),
		op("sell", 20.00, 2000),
		op("sell", 20.00, 2000),
		op("sell", 25.00, 1000),
	}
	expected := []domain.Tax{
		taxResult(0.0),
		taxResult(0.0),
		taxResult(0.0),
		taxResult(0.0),
		taxResult(3000.0),
	}

	processor := OperationProcessor{}
	result := processor.ProcessOperations(operations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case 6 failed: Expected %v, got %v", expected, result)
	}
}
