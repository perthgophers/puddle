package responses

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/perthgophers/puddle/messagerouter"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

// Portfolio is the command handling portfolio commands
type Portfolio struct{}

// PuddleWriter implements io.Writer
type PuddleWriter struct {
	w messagerouter.ResponseWriter
	r *messagerouter.CommandRequest
}

// WriteError implements io.Writer for errors
func (pw *PuddleWriter) Write(p []byte) (int, error) {
	err := pw.w.Write("\n```\n" + string(p) + "\n```\n")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return len(p), nil
}

// PuddleErrWriter implements io.Writer
type PuddleErrWriter struct {
	w messagerouter.ResponseWriter
	r *messagerouter.CommandRequest
}

// WriteError implements io.Writer for errors
func (pw *PuddleErrWriter) Write(p []byte) (int, error) {
	err := pw.w.WriteError(string(p))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return len(p), nil
}

// PortfolioGet gets the current user's portfolio
func (p *Portfolio) PortfolioGet(c *cli.Context, r *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	if r.Username != "nii236" {
		return errors.New("invalid username " + r.Username)
	}

	ethAmt := os.Getenv("ETH")
	btcAmt := os.Getenv("BTC")
	bchAmt := os.Getenv("BCH")

	if ethAmt == "" || btcAmt == "" || bchAmt == "" {
		return errors.New("ETH, BTC or BCH not provided in environment variables")
	}
	type HTTPResponse struct {
		Last string `json:"last"`
	}

	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		return errors.Wrap(err, "could not get ticker")
	}
	result := &Tickers{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return errors.Wrap(err, "could not decode json response from ticker API")
	}

	ethticker := result.Symbol("ETH")
	btcticker := result.Symbol("BTC")
	bchticker := result.Symbol("BCH")

	if ethticker == nil || btcticker == nil || bchticker == nil {
		return errors.New("could not find tickers from API response")
	}

	ethTotal, err := calcValue(ethAmt, ethticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}

	btcTotal, err := calcValue(btcAmt, btcticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}

	bchTotal, err := calcValue(bchAmt, bchticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}

	total := ethTotal + btcTotal + bchTotal
	audTotal, err := usdToAud(total)
	if err != nil {
		return errors.Wrap(err, "could not convert to AUD")
	}

	ethAmtFloat, err := strconv.ParseFloat(ethAmt, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse ethAmt to float")
	}
	btcAmtFloat, err := strconv.ParseFloat(btcAmt, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse btcAmt to float")
	}
	bchAmtFloat, err := strconv.ParseFloat(bchAmt, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse bchAmt to float")
	}

	ethPriceUSDFloat, err := strconv.ParseFloat(ethticker.PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse ethticker.PriceUSD to float")
	}
	btcPriceUSDFloat, err := strconv.ParseFloat(btcticker.PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse btcticker.PriceUSD to float")
	}
	bchPriceUSDFloat, err := strconv.ParseFloat(bchticker.PriceUsd, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse bchticker.PriceUSD to float")
	}

	w.Write(fmt.Sprintf(`Your crypto net worth is: %.2f USD (%.2f AUD)
	
%.2f ETH @ %.2f USD = %.2f USD
%.2f BTC @ %.2f USD = %.2f USD
%.2f BCH @ %.2f USD = %.2f USD`,
		total, audTotal,
		ethAmtFloat, ethPriceUSDFloat, ethTotal,
		btcAmtFloat, btcPriceUSDFloat, btcTotal,
		bchAmtFloat, bchPriceUSDFloat, bchTotal,
	))
	return nil

}

// PortfolioRegister adds the current user to the DB (not implemented)
func (p *Portfolio) PortfolioRegister(c *cli.Context, r *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	return w.Write("Register called (not implemented)")
}

func helpPrinter(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	bw := bufio.NewWriter(out)
	w := tabwriter.NewWriter(bw, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bw.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Run will return a string containing the current rate in USD
func (p *Portfolio) Run(r *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	cmd := &cli.App{
		Name:    "Portfolio Management",
		Version: "0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "John Nguyen",
				Email: "jtnguyen236@gmail.com",
			},
		},
		Usage:     "Disruption through innovation, blockchain, AI and machine learning",
		UsageText: "Add, track and manage your crypto portfolio!",
		Commands: []*cli.Command{
			{
				Name:    "register",
				Aliases: []string{"r"},
				Usage:   "Register yourself for portfolio tracking",
				Action: func(c *cli.Context) error {
					return p.PortfolioRegister(c, r, w)
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get your current portoflio",
				Action: func(c *cli.Context) error {
					return p.PortfolioGet(c, r, w)
				},
			},
		},
		Writer:    &PuddleWriter{w, r},
		ErrWriter: &PuddleErrWriter{w, r},
	}

	return cmd.Run(strings.Fields(r.Text))
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
	cli.HelpPrinter = helpPrinter
	cmd := &Portfolio{}
	Handle("!portfolio", cmd.Run)
}
