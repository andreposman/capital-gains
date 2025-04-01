package domain

import (
	"github.com/andreposman/capital-gains/pkg/helpers"
	"log"
)

type Portfolio struct {
	totalShares     int
	accumulatedLoss float64
	averageCost     float64
}

func (p *Portfolio) Buy(shareQuantity int, shareCost float64) {
	//TODO: validate/double check later
	totalCost := float64(p.totalShares)*p.averageCost + float64(shareQuantity)*shareCost
	p.totalShares += shareQuantity

	if p.totalShares > 0 {
		p.averageCost = helpers.ToFixedDecimal(totalCost/float64(p.totalShares), 2)
	}

	//log.Println("Total Cost is: R$", totalCost)
}

func (p *Portfolio) Sell(shareQuantity int, shareCost float64) float64 {
	if shareQuantity > p.totalShares {
		log.Fatal("error: insufficient amount of shares")
	}

	total := helpers.ToFixedDecimal(float64(shareQuantity)*shareCost, 2)
	cost := helpers.ToFixedDecimal(p.averageCost*float64(shareQuantity), 2)
	profit := total - cost

	p.totalShares -= shareQuantity

	tax := CalculateTax(p, total, profit)

	return tax
}
