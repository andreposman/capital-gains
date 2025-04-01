package domain

import (
	"github.com/andreposman/capital-gains/pkg/helpers"
	"math"
)

type Tax struct {
	Tax float64
}

var MAX_SALE_VALUE = 20000.00

func CalculateTax(p *Portfolio, totalSale, profit float64) float64 {
	if totalSale <= MAX_SALE_VALUE {
		updateLoss(p, profit)
		return 0.00
	}

	netProfit := math.Max(profit-p.accumulatedLoss, 0)
	p.accumulatedLoss = math.Max(p.accumulatedLoss-profit, 0)
	tax := helpers.ToFixedDecimal(netProfit*0.2, 2)

	//log.Println("Tax is: R$", tax)

	return tax
}

func updateLoss(p *Portfolio, profit float64) {
	if profit > 0 {
		p.accumulatedLoss += profit
	} else {
		p.accumulatedLoss -= profit
	}
}
