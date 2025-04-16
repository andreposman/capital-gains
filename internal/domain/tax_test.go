package domain

import (
	"math"
	"testing"
)

const floatTolerance = 1e-9

func floatsAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

// Test updateLoss directly (Example - may be redundant if covered by Portfolio tests)
func TestUpdateLoss_ExemptProfitReducesLoss(t *testing.T) {
	p := Portfolio{accumulatedLoss: 100.0}
	profit := 50.0
	updateLoss(&p, profit)
	expectedLoss := 50.0
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss to be %f after exempt profit, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestUpdateLoss_ExemptProfitExceedsLoss(t *testing.T) {
	p := Portfolio{accumulatedLoss: 100.0}
	profit := 150.0
	updateLoss(&p, profit)
	expectedLoss := 0.0
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss to be %f after exempt profit exceeded loss, got %f", expectedLoss, p.accumulatedLoss)
	}
}

func TestUpdateLoss_ExemptLossIncreasesLoss(t *testing.T) {
	p := Portfolio{accumulatedLoss: 100.0}
	profit := -50.0 // Represents a loss
	updateLoss(&p, profit)
	expectedLoss := 150.0
	if !floatsAlmostEqual(p.accumulatedLoss, expectedLoss) {
		t.Errorf("Expected accumulated loss to be %f after exempt loss, got %f", expectedLoss, p.accumulatedLoss)
	}
}
