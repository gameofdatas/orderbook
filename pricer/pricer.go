package pricer

import (
	"github.com/emirpasic/gods/maps/treemap"

	"pricer/models"
)

// Operations interface contains various operations to be performed on orders
type Operations interface {
	Add(user *models.Order)    // Add orders in a order book
	Reduce(user *models.Order) // Reduce orders from a order book
	Compute()                  // Compute calculate expense (if action is 'B') to buy target-size shares or income (if action is 'S') for selling target-size shares
}

// Pricer contains information of a order book
type Pricer struct {
	OrderMap       map[string]*models.Order // OrderMap is a map storing orderId as key and Order as value
	PriceMap       *treemap.Map             // PriceMap is a TreeMap to store the prices as key in sorted order and oder Size as value
	TotalSize      int64                    // TotalSize is current total standing size of a order book for a side
	TargetSize     int64                    // TargetSize is total size shares to be sold or bought for a side
	Side           models.Side              // Side of a order (B/S)
	PreviousSize   int64                    // PreviousSize stores previous size to to keep track if the size has changed for a side
	PreviousAmount float64                  // PreviousAmount store previous expense/income to keep track if the cost has changed for a side
}

// NewPricer initialize Pricer object
func NewPricer(side models.Side, targetSize int64) *Pricer {
	orderMap := make(map[string]*models.Order)
	var priceMap *treemap.Map
	if side == models.Buy {
		priceMap = treemap.NewWith(byPriceReverse) // Reverse treemap to store the price from high -> low for bids
	} else if side == models.Sell {
		priceMap = treemap.NewWith(byPrice) // treemap to store the price from low -> high for asks
	}
	return &Pricer{
		OrderMap:       orderMap,
		PriceMap:       priceMap,
		TotalSize:      0,
		TargetSize:     targetSize,
		Side:           side,
		PreviousAmount: float64(0),
		PreviousSize:   0,
	}
}

// Add adds an order to a OrderBook
func (p *Pricer) Add(order *models.Order) {
	if p.Side == order.Side {
		if size, found := p.PriceMap.Get(order.Price); found {
			p.PriceMap.Put(order.Price, size.(int64)+order.Size)
		} else {
			p.PriceMap.Put(order.Price, order.Size)
		}
		p.TotalSize += order.Size
		p.OrderMap[order.OrderID] = order
	}
}

// Reduce removes an order from a OrderBook
func (p *Pricer) Reduce(currentOrder *models.Order) {
	// check if the order already exist
	if previousOrder, found := p.OrderMap[currentOrder.OrderID]; found {
		// if found/exists
		reduceSize := currentOrder.Size
		price := previousOrder.Price
		previousOrder.Size -= reduceSize // Adjusting the size by reducing the size of the order for a order id
		currentPrice, _ := p.PriceMap.Get(price)
		p.PriceMap.Put(price, currentPrice.(int64)-reduceSize) // Reducing the size in the pricer Map too and update the new price.
		p.TotalSize -= reduceSize                              // Adjusting the current total standing size
		// If size is equal to or greater than the existing size of the order, the order is removed from the book.
		if previousOrder.Size == 0 {
			delete(p.OrderMap, currentOrder.OrderID)
		}
		// If price is equal to or greater than the existing price of the order, the price is removed from the price map.
		if currentPrice.(int64) == 0 {
			p.PriceMap.Remove(price)
		}
	}
}

// Compute calculate expense/income based on Bid/Ask order
func (p *Pricer) Compute() float64 {
	cost := float64(0)
	// Computing costs only when the current total standing size of a order book for a side has reached or crossed the target size
	if p.TotalSize >= p.TargetSize {
		remainingSize := p.TargetSize
		keys := p.PriceMap.Keys() // Get all keys of a price map in slice/array
		i := 0                    // counter to maintain index in the slice
		for remainingSize > 0 {   // the loop till the remainingSize is != 0
			currentPrice := keys[i]
			if currentSize, found := p.PriceMap.Get(currentPrice); found {
				if remainingSize >= currentSize.(int64) {
					// Calculates the total expense/income by multiplying size with price
					cost += currentPrice.(float64) * float64(currentSize.(int64))
					remainingSize -= currentSize.(int64)
				} else {
					// Calculates the total expense/income by multiplying size with price
					cost += currentPrice.(float64) * float64(remainingSize)
					remainingSize = 0
				}
			}
			i++
		}
		return cost
	}
	// return 0 if current total standing size is less than the target size
	return cost
}

// ProcessOrder process the order based on type mentioned in Order
func (p *Pricer) ProcessOrder(order *models.Order) {
	if order.Type == models.Add {
		p.Add(order)
	} else if order.Type == models.Remove {
		p.Reduce(order)
	}
}

// byPrice is a custom comparator for TreeMap which helps to store keys in ascending order
func byPrice(a, b interface{}) int {
	c1 := a.(float64)
	c2 := b.(float64)
	switch {
	case c1 > c2:
		return 1
	case c1 < c2:
		return -1
	default:
		return 0
	}
}

// byPriceReverse is a custom comparator for TreeMap which helps to store keys in descending order
func byPriceReverse(a, b interface{}) int {
	c1 := a.(float64)
	c2 := b.(float64)
	switch {
	case c1 < c2:
		return 1
	case c1 > c2:
		return -1
	default:
		return 0
	}
}
