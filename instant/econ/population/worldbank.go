package population

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// WorldBank holds settings for The World Bank API
type WorldBank struct {
	HTTPClient *http.Client
}

// TheWorldBankProvider indicates the data came from The World Bank Group
const TheWorldBankProvider Provider = "The World Bank"

// Fetch retrieves population data from The World Bank
func (w *WorldBank) Fetch(country string, from, to time.Time) (*Response, error) {
	u, err := w.buildURL(country, from, to)
	if err != nil {
		return nil, err
	}

	resp, err := w.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	wr := &WorldBankResponse{}

	if err = xml.Unmarshal(bdy, &wr); err != nil {
		return nil, err
	}

	r := &Response{
		Provider: TheWorldBankProvider,
	}

	for _, pp := range wr.Data {
		i := Instant{
			Date:  time.Date(pp.Date.Date, 12, 31, 0, 0, 0, 0, time.UTC),
			Value: pp.Value.Value,
		}
		r.History = append(r.History, i)
	}

	return r, err
}

func (w *WorldBank) buildURL(country string, from, to time.Time) (*url.URL, error) {
	// http://api.worldbank.org/v2/countries/it/indicators/SP.POP.TOTL?date=1960:2018
	u, err := url.Parse(fmt.Sprintf("http://api.worldbank.org/v2/countries/%v/indicators/SP.POP.TOTL", country))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("date", fmt.Sprintf("%d:%d", from.Year(), to.Year()))
	u.RawQuery = q.Encode()

	return u, err
}

// WorldBankResponse is the raw response received from The World Bank
type WorldBankResponse struct {
	//XMLName         xml.Name `xml:"data,omitempty" json:"data,omitempty"`
	AttrXmlnswb     string `xml:"xmlns wb,attr"  json:",omitempty"`
	Attrpage        string `xml:"page,attr"  json:",omitempty"`
	Attrpages       string `xml:"pages,attr"  json:",omitempty"`
	AttrPerPage     string `xml:"per_page,attr"  json:",omitempty"`
	Attrtotal       string `xml:"total,attr"  json:",omitempty"`
	Attrlastupdated string `xml:"lastupdated,attr"  json:",omitempty"`
	Data            []struct {
		Indicator struct {
			XMLName   xml.Name `xml:"indicator,omitempty" json:"indicator,omitempty"`
			ID        string   `xml:"id,attr"  json:",omitempty"`
			Indicator string   `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org indicator,omitempty" json:"indicator,omitempty"`
		Country struct {
			XMLName xml.Name `xml:"country,omitempty" json:"country,omitempty"`
			ID      string   `xml:"id,attr"  json:",omitempty"`
			Country string   `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org country,omitempty" json:"country,omitempty"`
		CountryISO3Code struct {
			XMLName         xml.Name `xml:"countryiso3code,omitempty" json:"countryiso3code,omitempty"`
			CountryISO3Code string   `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org countryiso3code,omitempty" json:"countryiso3code,omitempty"`
		Date struct {
			XMLName xml.Name `xml:"date,omitempty" json:"date,omitempty"`
			Date    int      `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org date,omitempty" json:"date,omitempty"`
		Value struct {
			XMLName xml.Name `xml:"value,omitempty" json:"value,omitempty"`
			Value   int      `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org value,omitempty" json:"value,omitempty"`
		Unit struct {
			XMLName xml.Name `xml:"unit,omitempty" json:"unit,omitempty"`
		} `xml:"http://www.worldbank.org unit,omitempty" json:"unit,omitempty"`
		OBSStatus struct {
			XMLName xml.Name `xml:"obs_status,omitempty" json:"obs_status,omitempty"`
		} `xml:"http://www.worldbank.org obs_status,omitempty" json:"obs_status,omitempty"`
		Decimal struct {
			XMLName xml.Name `xml:"decimal,omitempty" json:"decimal,omitempty"`
			Decimal string   `xml:",chardata" json:",omitempty"`
		} `xml:"http://www.worldbank.org decimal,omitempty" json:"decimal,omitempty"`
	} `xml:"http://www.worldbank.org data,omitempty" json:",omitempty"`
}
