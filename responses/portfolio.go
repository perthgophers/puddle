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
	err := pw.w.Write(string(p))
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
		return errors.New("portfolio not provided in environment variables")
	}
	type HTTPResponse struct {
		Last string `json:"last"`
	}

	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		return err
	}
	result := &Tickers{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}

	ETHticker := result.Symbol("ETH")
	BTCticker := result.Symbol("BTC")
	BCHticker := result.Symbol("BCH")

	if ETHticker == nil || BTCticker == nil || BCHticker == nil {
		return errors.New("could not find tickers from API response")
	}

	ethTotal, err := calcValue(ethAmt, ETHticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}
	btcTotal, err := calcValue(btcAmt, BTCticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}
	bchTotal, err := calcValue(bchAmt, BCHticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}

	total := ethTotal + btcTotal + bchTotal
	aud, err := usdToAud(total)
	if err != nil {
		return errors.Wrap(err, "could not convert to AUD")
	}
	w.Write(fmt.Sprintf("Your crypto net worth is: %.2f USD (%.2f AUD)", total, aud))
	return nil

}

// PortfolioRegister adds the current user to the DB (not implemented)
func (p *Portfolio) PortfolioRegister(c *cli.Context, r *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	fmt.Println("Register Called")
	return nil
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

	bw.Flush()
}

// Run will return a string containing the current rate in USD
func (p *Portfolio) Run(r *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	cmd := &cli.App{
		Name:      "Portfolio Management",
		UsageText: "Add, track and manage your crypto portfolio!",
		Commands: []*cli.Command{
			{
				Name:    "register",
				Aliases: []string{"r"},
				Usage:   "Register yourself for portfolio tracking",
				Action: func(c *cli.Context) error {
					p.PortfolioRegister(c, r, w)
					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get your current portoflio",
				Action: func(c *cli.Context) error {
					p.PortfolioGet(c, r, w)
					return nil
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
