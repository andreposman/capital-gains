package domain

import (
	"fmt"
	"github.com/andreposman/capital-gains/pkg/helpers"
)

type Portfolio struct {
	totalShares     int
	accumulatedLoss float64
	averageCost     float64
}

func (p *Portfolio) Buy(shareQuantity int, shareCost float64) {
	//calculo do valor total do ativo
	totalCost := float64(p.totalShares)*p.averageCost + float64(shareQuantity)*shareCost
	p.totalShares += shareQuantity

	if p.totalShares > 0 {
		p.averageCost = helpers.ToFixedDecimal(totalCost/float64(p.totalShares), 2)
	} else {
		p.averageCost = 0
	}

	//log.Println("Total Cost is: R$", totalCost)
}

// Sell updates the portfolio after a sell op and return the calculated tax and an error
func (p *Portfolio) Sell(shareQuantity int, shareCost float64) (float64, error) {
	if shareQuantity > p.totalShares {
		return 0.00, fmt.Errorf("insufficient shares: attempt to sell %d, but only %d have", shareQuantity, p.totalShares)
	}

	//calc o valor total e custo baseado no pm
	totalSellValue := helpers.ToFixedDecimal(float64(shareQuantity)*shareCost, 2)
	costSoldShares := helpers.ToFixedDecimal(p.averageCost*float64(shareQuantity), 2)
	profit := totalSellValue - costSoldShares

	//update qtd de acoes
	p.totalShares -= shareQuantity

	//vender tudo, reseta o custo e nao carrega pro futuro
	if p.totalShares == 0 {
		p.averageCost = 0.0
	}

	//calculo final de taxa
	tax := CalculateTax(p, totalSellValue, profit)

	return tax, nil
}
