package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/micro/cli"
	"github.com/perthgophers/puddle/messagerouter"
	"github.com/pkg/errors"
)

type PortfolioCommand struct {
	*cli.App
}

func NewPortfolioCommand() *PortfolioCommand {
	return &PortfolioCommand{}
}

func (app *PortfolioCommand) Initialise(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) {

	cmd := cli.NewApp()
	cmd.Name = "Portfolio Management"
	cmd.Usage = "Add, track and manage your crypto portfolio!"
	cmd.Commands = []cli.Command{
		{
			Name:    "register",
			Aliases: []string{"r"},
			Usage:   "Register yourself for portfolio tracking",
			Action:  PortfolioRegister(cr, w),
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Get your current portoflio",
			Action:  PortfolioGet(cr, w),
		},
	}

	app.App = cmd
}

// PortfolioRegister registers the user for portfolio tracking
func PortfolioRegister(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) func(c *cli.Context) {
	return func(c *cli.Context) {
		fmt.Println("Register")
		return
	}
}

// PortfolioGet gets the current user's portfolio
func PortfolioGet(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) func(c *cli.Context) {
	return func(c *cli.Context) {
		if cr.Username != "nii236" {
			w.WriteError(errors.New("invalid username " + cr.Username).Error())
			return
		}

		ethAmt := os.Getenv("ETH")
		btcAmt := os.Getenv("BTC")
		bchAmt := os.Getenv("BCH")

		if ethAmt == "" || btcAmt == "" || bchAmt == "" {
			w.WriteError(errors.New("portfolio not provided in environment variables").Error())
			return
		}
		type HTTPResponse struct {
			Last string `json:"last"`
		}

		resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
		if err != nil {
			w.WriteError(err.Error())
		}
		result := &Tickers{}
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			w.WriteError(err.Error())
		}

		ETHticker := result.Symbol("ETH")
		BTCticker := result.Symbol("BTC")
		BCHticker := result.Symbol("BCH")

		if ETHticker == nil || BTCticker == nil || BCHticker == nil {
			w.WriteError(errors.New("could not find tickers from API response").Error())
			return
		}

		ethTotal, err := calcValue(ethAmt, ETHticker.PriceUsd)
		if err != nil {
			w.WriteError(errors.Wrap(err, "could not calculate total value").Error())

		}
		btcTotal, err := calcValue(btcAmt, BTCticker.PriceUsd)
		if err != nil {
			w.WriteError(errors.Wrap(err, "could not calculate total value").Error())

		}
		bchTotal, err := calcValue(bchAmt, BCHticker.PriceUsd)
		if err != nil {
			w.WriteError(errors.Wrap(err, "could not calculate total value").Error())

		}

		total := ethTotal + btcTotal + bchTotal
		aud, err := usdToAud(total)
		if err != nil {
			w.WriteError(errors.Wrap(err, "could not convert to AUD").Error())

		}
		w.Write(fmt.Sprintf("Your crypto net worth is: %.2f USD (%.2f AUD)", total, aud))
	}

}

func (cmd *PortfolioCommand) Run(cr *messagerouter.CommandRequest, rw messagerouter.ResponseWriter) error {
	msg := strings.Fields(cr.Text)

	err := cmd.App.Run(msg)
	if err != nil {
		return err
	}
	return nil

}

type CLIAdaptor interface {
	Initialise(cr *messagerouter.CommandRequest, rw messagerouter.ResponseWriter)
	Run(cr *messagerouter.CommandRequest, rw messagerouter.ResponseWriter) error
}

// PortfolioWrapper
func PortfolioWrapper(app CLIAdaptor) func(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	return func(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
		app.Initialise(cr, w)
		app.Run(cr, w)
		return nil
	}
}

func calcValue(ownAmt, tickerAmt string) (float64, error) {
	own, err := strconv.ParseFloat(ownAmt, 64)
	if err != nil {
		return 0, err
	}
	tick, err := strconv.ParseFloat(tickerAmt, 64)
	if err != nil {
		return 0, err
	}
	return own * tick, err
}

type forexTicker struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates struct {
		AUD float64 `json:"AUD"`
	} `json:"rates"`
}

func usdToAud(usd float64) (float64, error) {

	resp, err := http.Get("https://api.fixer.io/latest?symbols=USD,AUD")
	if err != nil {
		return 0, err
	}
	result := &forexTicker{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return 0, err
	}

	converted := result.Rates.AUD * usd
	return converted, nil
}

func init() {
	Handle("!portfolio", PortfolioWrapper(&PortfolioCommand{}))
}
