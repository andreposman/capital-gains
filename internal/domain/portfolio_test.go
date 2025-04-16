package domain

import (
	"math"
	"testing"
)

// --- Buy Method Tests ---
func TestPortfolio_Buy_FirstPurchase(t *testing.T) {
	p := Portfolio{}
	p.Buy(100, 10.00)

	if p.totalShares != 100 {
		t.Errorf("Expected totalShares to be 100, got %d", p.totalShares)
	}
	if !floatsAlmostEqual(p.averageCost, 10.00) {
		t.Errorf("Expected averageCost to be 10.00, got %f", p.averageCost)
	}
}

func TestPortfolio_Buy_MultiplePurchasesWAC(t *testing.T) {
	p := Portfolio{}
	p.Buy(100, 10.00) // Total cost = 1000
	p.Buy(50, 20.00)  // Total cost = 1000

	// Expected: (100 * 10 + 50 * 20) / (100 + 50) = (1000 + 1000) / 150 = 2000 / 150 = 13.333...
	expectedShares := 150
	// Need to re-calculate expected WAC using the *exact same* rounding logic as the code
	calculatedExpectedWAC := math.Round((100.0*10.00+50.0*20.00)/150.0*100) / 100

	if p.totalShares != expectedShares {
		t.Errorf("Expected totalShares to be %d, got %d", expectedShares, p.totalShares)
	}
	if !floatsAlmostEqual(p.averageCost, calculatedExpectedWAC) {
		t.Errorf("Expected averageCost to be %f (calculated), got %f", calculatedExpectedWAC, p.averageCost)
	}
}

// --- Sell Method Tests (includes tax logic implicitly) ---

func TestPortfolio_Sell_Profit_Taxable_NoLoss(t *testing.T) {
	p := Portfolio{}
	p.Buy(10000, 10.00) // WAC = 10.00

	// Sell 5000 @ 20.00. Total Sale = 100,000 (> 20k). Profit = 5000 * (20 - 10) = 50,000
	// Tax = 50,000 * 0.20 = 10,000
	sellQuantity := 5000
	sellPrice := 20.00
	expectedTax := 10000.00
	expectedSharesLeft := 5000
	expectedLoss := 0.0

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_Profit_Exempt_NoLoss(t *testing.T) {
	p := Portfolio{}
	p.Buy(100, 10.00) // WAC = 10.00

	// Sell 50 @ 15.00. Total Sale = 750 (<= 20k). Profit = 50 * (15 - 10) = 250
	// Tax = 0 (exempt)
	sellQuantity := 50
	sellPrice := 15.00
	expectedTax := 0.00
	expectedSharesLeft := 50
	expectedLoss := 0.0

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_Loss_Taxable_NoLoss(t *testing.T) {
	p := Portfolio{}
	p.Buy(10000, 20.00) // WAC = 20.00

	// Sell 5000 @ 10.00. Total Sale = 50,000 (> 20k). Loss = 5000 * (10 - 20) = -50,000
	// Tax = 0. Accumulated Loss = 50,000
	sellQuantity := 5000
	sellPrice := 10.00
	expectedTax := 0.00
	expectedSharesLeft := 5000
	expectedLoss := 50000.00 // Absolute value of the loss

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_Loss_Exempt_NoLoss(t *testing.T) {
	p := Portfolio{}
	p.Buy(100, 20.00) // WAC = 20.00

	// Sell 50 @ 10.00. Total Sale = 500 (<= 20k). Loss = 50 * (10 - 20) = -500
	// Tax = 0. Accumulated Loss = 500
	sellQuantity := 50
	sellPrice := 10.00
	expectedTax := 0.00
	expectedSharesLeft := 50
	expectedLoss := 500.00 // Absolute value of the loss

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_Profit_Taxable_WithLoss_Consumed(t *testing.T) {
	p := Portfolio{accumulatedLoss: 25000.00} // Start with loss
	p.Buy(10000, 10.00)                       // WAC = 10.00

	// Sell 5000 @ 20.00. Total Sale = 100,000 (> 20k). Gross Profit = 50,000
	// Net Profit = 50,000 - 25,000 = 25,000
	// Tax = 25,000 * 0.20 = 5,000
	sellQuantity := 5000
	sellPrice := 20.00
	expectedTax := 5000.00
	expectedSharesLeft := 5000
	expectedLossAfter := 0.0 // Loss fully consumed

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLossAfter) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLossAfter, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_Profit_Exempt_WithLoss_Consumed(t *testing.T) {
	p := Portfolio{accumulatedLoss: 500.00} // Start with loss
	p.Buy(100, 10.00)                       // WAC = 10.00

	// Sell 50 @ 25.00. Total Sale = 1250 (<= 20k). Gross Profit = 50 * (25 - 10) = 750
	// Tax = 0 (exempt)
	// Remaining Loss = max(0, 500 - 750) = 0
	sellQuantity := 50
	sellPrice := 25.00
	expectedTax := 0.00
	expectedSharesLeft := 50
	expectedLossAfter := 0.00

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err != nil {
		t.Fatalf("Sell returned unexpected error: %v", err)
	}
	if !floatsAlmostEqual(tax, expectedTax) {
		t.Errorf("Expected tax %f, got %f", expectedTax, tax)
	}
	if p.totalShares != expectedSharesLeft {
		t.Errorf("Expected %d shares left, got %d", expectedSharesLeft, p.totalShares)
	}
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLossAfter) {
		t.Errorf("Expected accumulated loss %f, got %f", expectedLossAfter, p.accumulatedLoss)
	}
}

func TestPortfolio_Sell_InsufficientShares(t *testing.T) {
	p := Portfolio{}
	p.Buy(50, 10.00) // Only 50 shares

	sellQuantity := 100 // Try to sell more
	sellPrice := 15.00

	tax, err := p.Sell(sellQuantity, sellPrice)

	if err == nil {
		t.Fatalf("Expected an error for insufficient shares, but got nil")
	}

	if !floatsAlmostEqual(tax, 0.0) {
		t.Errorf("Expected tax to be 0 on error, got %f", tax)
	}
	if p.totalShares != 50 {
		t.Errorf("Expected totalShares to remain 50 on error, got %d", p.totalShares)
	}
	if !floatsAlmostEqual(p.averageCost, 10.00) {
		t.Errorf("Expected averageCost to remain 10.00 on error, got %f", p.averageCost)
	}
}
