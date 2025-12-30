# Mudrex Go SDK

[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub](https://img.shields.io/github/stars/DecentralizedJM/mudrex-go-sdk?style=social)](https://github.com/DecentralizedJM/mudrex-go-sdk)

**Unofficial Go SDK for [Mudrex Futures Trading API](https://docs.trade.mudrex.com/docs/overview)** - High-performance trading client for Go developers.

**Built and maintained by [DecentralizedJM](https://github.com/DecentralizedJM)**

## 🚀 Features

- **Simple & Intuitive** - Go-native interface that feels natural
- **Type-Safe** - Full type definitions and models
- **High-Performance** - Built-in rate limiting and efficient request handling
- **Well-Documented** - Comprehensive examples and docstrings
- **Production-Ready** - Designed for real trading applications

## 📦 Installation

```bash
go get github.com/DecentralizedJM/mudrex-go-sdk
```

## ⚡ Quick Start

```go
package main

import (
	"fmt"
	"log"
	
	mudrex "github.com/DecentralizedJM/mudrex-go-sdk"
)

func main() {
	// Initialize the client
	client := mudrex.NewClient("your-api-secret")
	defer client.Close()
	
	// Check your balance
	balance, err := client.Wallet.GetFuturesBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %s %s\n", balance.Balance, balance.Currency)
	
	// List tradable assets
	assets, err := client.Assets.ListAll(1, 10, "", "")
	if err != nil {
		log.Fatal(err)
	}
	for _, asset := range assets {
		fmt.Printf("%s: up to %sx leverage\n", asset.Symbol, asset.MaxLeverage)
	}
	
	// Set leverage before trading
	leverage, err := client.Leverage.Set("BTCUSDT", "10", mudrex.MarginTypeIsolated)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Leverage set to %s\n", leverage.Leverage)
	
	// Place a market order
	order, err := client.Orders.CreateMarketOrder(
		"BTCUSDT",
		mudrex.OrderTypeLong,
		"0.001",
		"10",
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order placed: %s\n", order.OrderID)
	
	// Monitor positions
	positions, err := client.Positions.ListOpen()
	if err != nil {
		log.Fatal(err)
	}
	for _, pos := range positions {
		fmt.Printf("%s: %s PnL\n", pos.Symbol, pos.UnrealizedPnL)
	}
}
```

## 📚 Documentation

### API Modules

| Module | Description |
|--------|-------------|
| `client.Wallet` | Spot & futures wallet balances, fund transfers |
| `client.Assets` | Discover tradable instruments, get specifications |
| `client.Leverage` | Get/set leverage and margin type |
| `client.Orders` | Create, view, cancel, and amend orders |
| `client.Positions` | Manage positions, set SL/TP, close/reverse |
| `client.Fees` | View trading fee history |

### Complete Trading Workflow

```go
// 1. Check available balance
balance, _ := client.Wallet.GetFuturesBalance()

// 2. Discover available assets
assets, _ := client.Assets.ListAll(1, 50, "symbol", "asc")

// 3. Set leverage for the asset
client.Leverage.Set("BTCUSDT", "5", mudrex.MarginTypeIsolated)

// 4. Place an order with stop loss and take profit
order, _ := client.Orders.CreateMarketOrder(
	"BTCUSDT",
	mudrex.OrderTypeLong,
	"0.001",
	"5",
)

// 5. Monitor your position
positions, _ := client.Positions.ListOpen()
for _, pos := range positions {
	if pos.PositionID == order.OrderID {
		// Set stop loss at 5% below entry
		sl := "95000"
		client.Positions.SetStopLoss(pos.PositionID, sl)
	}
}

// 6. Close position when ready
client.Positions.Close(order.OrderID)
```

## 🔧 Configuration

### Custom Base URL and Timeout

```go
client := mudrex.NewClientWithConfig(
	"your-api-secret",
	"https://custom.mudrex.com/fapi/v1",
	30*time.Second,
)
```

## ⚠️ Error Handling

```go
import "github.com/DecentralizedJM/mudrex-go-sdk"

balance, err := client.Wallet.GetFuturesBalance()
if err != nil {
	switch e := err.(type) {
	case *mudrex.AuthenticationError:
		fmt.Println("Invalid API secret")
	case *mudrex.RateLimitError:
		fmt.Println("Rate limit exceeded, waiting...")
	case *mudrex.InsufficientBalanceError:
		fmt.Println("Insufficient balance for this trade")
	case *mudrex.ValidationError:
		fmt.Println("Invalid request parameters")
	default:
		fmt.Printf("Error: %v\n", e)
	}
}
```

## 🧪 Testing

```bash
go test ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 👥 Contributors

- [@DecentralizedJM](https://github.com/DecentralizedJM) - Creator & Maintainer

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🔗 Links

- [Mudrex Trading API Docs](https://docs.trade.mudrex.com/docs/overview)
- [API Quick Reference](https://docs.trade.mudrex.com/docs/overview)
- [Mudrex Platform](https://mudrex.com)

## ⚠️ Disclaimer

**This is an UNOFFICIAL SDK.** This SDK is for educational and informational purposes. Cryptocurrency trading involves significant risk. Always:
- Start with small amounts
- Use proper risk management (stop-losses)
- Never trade more than you can afford to lose
- Test thoroughly in a safe environment first

---

Built and maintained by [DecentralizedJM](https://github.com/DecentralizedJM) with ❤️
