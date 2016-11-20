package responses

import (
	// "encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dutchcoders/goftp"
	"github.com/perthgophers/puddle/messagerouter"
	"io"
)

// Weather prints out the current temperature & forecast.
// Usage: !weather
func Weather(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	var temperature string
	var forecast string
	var wdata W_product
	var wtemp WTEMP_product
	var ftp *goftp.FTP
	var err error

	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	if ftp, err = goftp.Connect("ftp.bom.gov.au:21"); err != nil {
		panic(err)
	}
	defer ftp.Close()
	fmt.Println("Successfully connected to bom")

	// TLS client authentication
	// Username / password authentication
	if err = ftp.Login("anonymous", ""); err != nil {
		panic(err)
	}
	_, err = ftp.Retr("/anon/gen/fwo/IDW60920.xml", func(r io.Reader) error {
		decoder := xml.NewDecoder(r)

		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == "product" && se.Name.Space == "" {
					fmt.Println("DECODING")
					decoder.DecodeElement(&wtemp, &se)
				}
			}
		}
		temperature = wtemp.WTEMP_observations.WTEMP_station[0].WTEMP_period.WTEMP_level.WTEMP_element[0].Text

		return nil
	})

	_, err = ftp.Retr("/anon/gen/fwo/IDW12400.xml", func(r io.Reader) error {
		decoder := xml.NewDecoder(r)

		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				fmt.Println("Decode")
				if se.Name.Local == "product" && se.Name.Space == "" {
					decoder.DecodeElement(&wdata, &se)
				}
			}
		}

		forecast = wdata.W_forecast.W_area[1].W_forecast_period[0].W_text[0].Text
		return err
	})

	w.Write(fmt.Sprintf("Current Temperature: %sC°\n%s", temperature, forecast))

	return nil
}

func init() {
	Handle("!weather", Weather)
}

///////////////////////////
/// structs
///////////////////////////

type W_amoc struct {
	W_expiry_time                   *W_expiry_time                   `xml:" expiry-time,omitempty" json:"expiry-time,omitempty"`
	W_identifier                    *W_identifier                    `xml:" identifier,omitempty" json:"identifier,omitempty"`
	W_issue_time_local              *W_issue_time_local              `xml:" issue-time-local,omitempty" json:"issue-time-local,omitempty"`
	W_issue_time_utc                *W_issue_time_utc                `xml:" issue-time-utc,omitempty" json:"issue-time-utc,omitempty"`
	W_next_routine_issue_time_local *W_next_routine_issue_time_local `xml:" next-routine-issue-time-local,omitempty" json:"next-routine-issue-time-local,omitempty"`
	W_next_routine_issue_time_utc   *W_next_routine_issue_time_utc   `xml:" next-routine-issue-time-utc,omitempty" json:"next-routine-issue-time-utc,omitempty"`
	W_phase                         *W_phase                         `xml:" phase,omitempty" json:"phase,omitempty"`
	W_product_type                  *W_product_type                  `xml:" product-type,omitempty" json:"product-type,omitempty"`
	W_sent_time                     *W_sent_time                     `xml:" sent-time,omitempty" json:"sent-time,omitempty"`
	W_service                       *W_service                       `xml:" service,omitempty" json:"service,omitempty"`
	W_source                        *W_source                        `xml:" source,omitempty" json:"source,omitempty"`
	W_status                        *W_status                        `xml:" status,omitempty" json:"status,omitempty"`
	W_sub_service                   *W_sub_service                   `xml:" sub-service,omitempty" json:"sub-service,omitempty"`
	W_validity_bgn_time_local       *W_validity_bgn_time_local       `xml:" validity-bgn-time-local,omitempty" json:"validity-bgn-time-local,omitempty"`
	W_validity_end_time_local       *W_validity_end_time_local       `xml:" validity-end-time-local,omitempty" json:"validity-end-time-local,omitempty"`
}

type W_area struct {
	Attr_aac          string               `xml:" aac,attr"  json:",omitempty"`
	Attr_description  string               `xml:" description,attr"  json:",omitempty"`
	Attr_parent_aac   string               `xml:" parent-aac,attr"  json:",omitempty"`
	Attr_type         string               `xml:" type,attr"  json:",omitempty"`
	W_forecast_period []*W_forecast_period `xml:" forecast-period,omitempty" json:"forecast-period,omitempty"`
}

type W_copyright struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_disclaimer struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_expiry_time struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_forecast struct {
	W_area []*W_area `xml:" area,omitempty" json:"area,omitempty"`
}

type W_forecast_period struct {
	Attr_end_time_local   string    `xml:" end-time-local,attr"  json:",omitempty"`
	Attr_end_time_utc     string    `xml:" end-time-utc,attr"  json:",omitempty"`
	Attr_index            string    `xml:" index,attr"  json:",omitempty"`
	Attr_start_time_local string    `xml:" start-time-local,attr"  json:",omitempty"`
	Attr_start_time_utc   string    `xml:" start-time-utc,attr"  json:",omitempty"`
	W_text                []*W_text `xml:" text,omitempty" json:"text,omitempty"`
}

type W_identifier struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_issue_time_local struct {
	Attr_tz string `xml:" tz,attr"  json:",omitempty"`
	Text    string `xml:",chardata" json:",omitempty"`
}

type W_issue_time_utc struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_next_routine_issue_time_local struct {
	Attr_tz string `xml:" tz,attr"  json:",omitempty"`
	Text    string `xml:",chardata" json:",omitempty"`
}

type W_next_routine_issue_time_utc struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_office struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_p struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_phase struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_product struct {
	Attr_xsi_noNamespaceSchemaLocation string      `xml:"http://www.w3.org/2001/XMLSchema-instance noNamespaceSchemaLocation,attr"  json:",omitempty"`
	Attr_version                       string      `xml:" version,attr"  json:",omitempty"`
	Attr_xsi                           string      `xml:"xmlns xsi,attr"  json:",omitempty"`
	W_amoc                             *W_amoc     `xml:" amoc,omitempty" json:"amoc,omitempty"`
	W_forecast                         *W_forecast `xml:" forecast,omitempty" json:"forecast,omitempty"`
}

type W_product_type struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_region struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_root struct {
	W_product *W_product `xml:" product,omitempty" json:"product,omitempty"`
}

type W_sender struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_sent_time struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_service struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_source struct {
	W_copyright  *W_copyright  `xml:" copyright,omitempty" json:"copyright,omitempty"`
	W_disclaimer *W_disclaimer `xml:" disclaimer,omitempty" json:"disclaimer,omitempty"`
	W_office     *W_office     `xml:" office,omitempty" json:"office,omitempty"`
	W_region     *W_region     `xml:" region,omitempty" json:"region,omitempty"`
	W_sender     *W_sender     `xml:" sender,omitempty" json:"sender,omitempty"`
}

type W_status struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_sub_service struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type W_text struct {
	Attr_type string `xml:" type,attr"  json:",omitempty"`
	Text      string `xml:",chardata" json:",omitempty"`
	W_p       []*W_p `xml:" p,omitempty" json:"p,omitempty"`
}

type W_validity_bgn_time_local struct {
	Attr_tz string `xml:" tz,attr"  json:",omitempty"`
	Text    string `xml:",chardata" json:",omitempty"`
}

type W_validity_end_time_local struct {
	Attr_tz string `xml:" tz,attr"  json:",omitempty"`
	Text    string `xml:",chardata" json:",omitempty"`
}

///////////////////////////

///////////////////////////
/// structs
///////////////////////////

type WTEMP_amoc struct {
	WTEMP_identifier       *WTEMP_identifier       `xml:" identifier,omitempty" json:"identifier,omitempty"`
	WTEMP_issue_time_local *WTEMP_issue_time_local `xml:" issue-time-local,omitempty" json:"issue-time-local,omitempty"`
	WTEMP_issue_time_utc   *WTEMP_issue_time_utc   `xml:" issue-time-utc,omitempty" json:"issue-time-utc,omitempty"`
	WTEMP_phase            *WTEMP_phase            `xml:" phase,omitempty" json:"phase,omitempty"`
	WTEMP_product_type     *WTEMP_product_type     `xml:" product-type,omitempty" json:"product-type,omitempty"`
	WTEMP_sent_time        *WTEMP_sent_time        `xml:" sent-time,omitempty" json:"sent-time,omitempty"`
	WTEMP_service          *WTEMP_service          `xml:" service,omitempty" json:"service,omitempty"`
	WTEMP_source           *WTEMP_source           `xml:" source,omitempty" json:"source,omitempty"`
	WTEMP_status           *WTEMP_status           `xml:" status,omitempty" json:"status,omitempty"`
}

type WTEMP_copyright struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_disclaimer struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_element struct {
	Attr_duration         string `xml:" duration,attr"  json:",omitempty"`
	Attr_end_time_local   string `xml:" end-time-local,attr"  json:",omitempty"`
	Attr_end_time_utc     string `xml:" end-time-utc,attr"  json:",omitempty"`
	Attr_instance         string `xml:" instance,attr"  json:",omitempty"`
	Attr_start_time_local string `xml:" start-time-local,attr"  json:",omitempty"`
	Attr_start_time_utc   string `xml:" start-time-utc,attr"  json:",omitempty"`
	Attr_time_local       string `xml:" time-local,attr"  json:",omitempty"`
	Attr_time_utc         string `xml:" time-utc,attr"  json:",omitempty"`
	Attr_type             string `xml:" type,attr"  json:",omitempty"`
	Attr_units            string `xml:" units,attr"  json:",omitempty"`
	Text                  string `xml:",chardata" json:",omitempty"`
}

type WTEMP_identifier struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_issue_time_local struct {
	Attr_tz string `xml:" tz,attr"  json:",omitempty"`
	Text    string `xml:",chardata" json:",omitempty"`
}

type WTEMP_issue_time_utc struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_level struct {
	Attr_index    string           `xml:" index,attr"  json:",omitempty"`
	Attr_type     string           `xml:" type,attr"  json:",omitempty"`
	WTEMP_element []*WTEMP_element `xml:" element,omitempty" json:"element,omitempty"`
}

type WTEMP_observations struct {
	WTEMP_station []*WTEMP_station `xml:" station,omitempty" json:"station,omitempty"`
}

type WTEMP_office struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_period struct {
	Attr_index      string       `xml:" index,attr"  json:",omitempty"`
	Attr_time_local string       `xml:" time-local,attr"  json:",omitempty"`
	Attr_time_utc   string       `xml:" time-utc,attr"  json:",omitempty"`
	Attr_wind_src   string       `xml:" wind-src,attr"  json:",omitempty"`
	WTEMP_level     *WTEMP_level `xml:" level,omitempty" json:"level,omitempty"`
}

type WTEMP_phase struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_product struct {
	Attr_xsi_noNamespaceSchemaLocation string              `xml:"http://www.w3.org/2001/XMLSchema-instance noNamespaceSchemaLocation,attr"  json:",omitempty"`
	Attr_version                       string              `xml:" version,attr"  json:",omitempty"`
	Attr_xsi                           string              `xml:"xmlns xsi,attr"  json:",omitempty"`
	WTEMP_amoc                         *WTEMP_amoc         `xml:" amoc,omitempty" json:"amoc,omitempty"`
	WTEMP_observations                 *WTEMP_observations `xml:" observations,omitempty" json:"observations,omitempty"`
}

type WTEMP_product_type struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_region struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_root struct {
	WTEMP_product *WTEMP_product `xml:" product,omitempty" json:"product,omitempty"`
}

type WTEMP_sender struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_sent_time struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_service struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type WTEMP_source struct {
	WTEMP_copyright  *WTEMP_copyright  `xml:" copyright,omitempty" json:"copyright,omitempty"`
	WTEMP_disclaimer *WTEMP_disclaimer `xml:" disclaimer,omitempty" json:"disclaimer,omitempty"`
	WTEMP_office     *WTEMP_office     `xml:" office,omitempty" json:"office,omitempty"`
	WTEMP_region     *WTEMP_region     `xml:" region,omitempty" json:"region,omitempty"`
	WTEMP_sender     *WTEMP_sender     `xml:" sender,omitempty" json:"sender,omitempty"`
}

type WTEMP_station struct {
	Attr_bom_id               string        `xml:" bom-id,attr"  json:",omitempty"`
	Attr_description          string        `xml:" description,attr"  json:",omitempty"`
	Attr_forecast_district_id string        `xml:" forecast-district-id,attr"  json:",omitempty"`
	Attr_lat                  string        `xml:" lat,attr"  json:",omitempty"`
	Attr_lon                  string        `xml:" lon,attr"  json:",omitempty"`
	Attr_stn_height           string        `xml:" stn-height,attr"  json:",omitempty"`
	Attr_stn_name             string        `xml:" stn-name,attr"  json:",omitempty"`
	Attr_type                 string        `xml:" type,attr"  json:",omitempty"`
	Attr_tz                   string        `xml:" tz,attr"  json:",omitempty"`
	Attr_wmo_id               string        `xml:" wmo-id,attr"  json:",omitempty"`
	WTEMP_period              *WTEMP_period `xml:" period,omitempty" json:"period,omitempty"`
}

type WTEMP_status struct {
	Text string `xml:",chardata" json:",omitempty"`
}

///////////////////////////
