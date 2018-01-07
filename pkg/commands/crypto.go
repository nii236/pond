package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"gopkg.in/urfave/cli.v2"
)

// Tickers contain a slice of Ticker
type Tickers []*Ticker

// Ticker contains all the data for a single crypto
type Ticker struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

// Get gets the ticker from a slice of tickers
func (ts Tickers) Get(sym string) *Ticker {
	for _, v := range ts {
		if v.Symbol == strings.ToUpper(sym) {
			return v
		}
	}
	return nil
}
func newCryptoCmd() *cli.Command {
	cmd := &cli.Command{
		Name:    "crypto",
		Usage:   "Toolkit for crypto currencies",
		Aliases: []string{"c"},
		Subcommands: []*cli.Command{
			{
				Name:   "all",
				Action: cryptoAll,
			},
			{
				Name:   "get",
				Action: get,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "symbol",
						Aliases: []string{"s"},
					},
				},
			},
		},
	}

	return cmd
}
func getTickers() (*Tickers, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/?limit=200")
	if err != nil {
		return nil, err
	}
	result := &Tickers{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func get(c *cli.Context) error {
	if c.String("symbol") == "" {
		return errors.New("please provide a symbol")
	}

	tickers, err := getTickers()
	if err != nil {
		return errors.Wrap(err, "could not fetch tickers")
	}

	sym := c.String("symbol")
	ticker := tickers.Get(sym)
	if ticker == nil {
		return errors.New("Ticker not found: " + sym)
	}

	fmt.Fprintf(c.App.Writer, "%s USD", ticker.PriceUsd)

	return nil
}

func cryptoAll(c *cli.Context) error {
	tickers, err := getTickers()
	if err != nil {
		return errors.Wrap(err, "could not fetch tickers")
	}

	bch, err := strconv.ParseFloat(tickers.Get("bch").PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse float")
	}
	btc, err := strconv.ParseFloat(tickers.Get("btc").PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse float")
	}
	eth, err := strconv.ParseFloat(tickers.Get("eth").PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse float")
	}

	fmt.Fprintf(c.App.Writer, "```\n"+
		`BTC %.2f USD
BCH %.2f USD
ETH %.2f USD
`+"```", btc, bch, eth)

	return nil
}

func init() {
	crypto := newCryptoCmd()
	Commands = append(Commands, crypto)
}
