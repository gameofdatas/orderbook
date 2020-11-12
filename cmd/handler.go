package cmd

import (
	"fmt"

	"pricer/models"
	"pricer/pricer"
)

// handleOrderEvents process orders one by one as it comes
func handleOrderEvents(orders string, buyObject *pricer.Pricer, sellObject *pricer.Pricer) error {
	// Parse the incoming order string to Order
	order, err := models.NewOrder(orders)
	if err != nil {
		return err
	}
	// Perform Add/Reduce based on type for a Buy Order
	buyObject.ProcessOrder(order)
	// Perform Add/Reduce based on type for a Sell Order
	sellObject.ProcessOrder(order)

	// PrintOutput only if there is a change in size for buy object
	if buyObject.TotalSize != buyObject.PreviousSize {
		buyObject.PreviousSize = buyObject.TotalSize
		newCost := buyObject.Compute()
		printOutput(order, models.Sell, newCost, buyObject.PreviousAmount)
		buyObject.PreviousAmount = newCost
	}
	// PrintOutput only if there is a change in size for sell object
	if sellObject.TotalSize != sellObject.PreviousSize {
		sellObject.PreviousSize = sellObject.TotalSize
		newCost := sellObject.Compute()
		printOutput(order, models.Buy, newCost, sellObject.PreviousAmount)
		sellObject.PreviousAmount = newCost
	}
	return nil
}

// printOutput prints the new income/expense for buying and selling the target shares
// if the total size for a side crosses more than target size and if currentCost is not same as previous cost
func printOutput(order *models.Order, side models.Side, currentAmount float64, previousAmount float64) {
	var result string
	if currentAmount != previousAmount {
		if currentAmount != float64(0) {
			result = fmt.Sprintf("%d %s %.2f", order.Timestamp, side, currentAmount)
		} else {
			// If the book does not contain target-size shares in the appropriate type of order (asks for expense; bids for income), the total field contains the string 'NA'.
			result = fmt.Sprintf("%d %s NA", order.Timestamp, side)
		}

	}
	if result != "" {
		fmt.Println(result)
	}
}
