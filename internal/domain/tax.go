package domain

import (
	"github.com/andreposman/capital-gains/pkg/helpers"
	"math"
)

type Tax struct {
	Tax float64 `json:"tax"`
}

const (
	MAX_SALE_VALUE = 20000.00
	TAX_RATE       = 0.20
)

// CalculateTax determines the tax for a sell op
func CalculateTax(p *Portfolio, totalSale, profit float64) float64 {
	if totalSale <= MAX_SALE_VALUE {
		// mesmo isento, afeta o valor acumulado
		updateLoss(p, profit)
		return 0.00
	}

	// venda potencialmente taxavel, calculando o netProfit considerando o loss
	netProfit := math.Max(profit-p.accumulatedLoss, 0)

	/**
		atualizando accumlatedLoss:
			if profit > 0, loss é consumido
	 		if profit <= 0, loss é aumentado
	**/
	p.accumulatedLoss = math.Max(p.accumulatedLoss-profit, 0)

	//calculando a taxa no netProfit
	tax := helpers.ToFixedDecimal(netProfit*TAX_RATE, 2)

	//log.Println("Tax is: R$", tax)
	return tax
}

// updateLoss adjusts the accumulated loss for exempt sales
func updateLoss(p *Portfolio, profit float64) {
	if profit > 0 {
		p.accumulatedLoss = math.Max(p.accumulatedLoss-profit, 0)
	} else {
		p.accumulatedLoss -= profit
	}
}
