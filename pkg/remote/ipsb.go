package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/abc1763613206/nabili/pkg/common"
)

const (
	IpsbAPIURL = "https://api.ip.sb/geoip"
)

type IpsbResponse struct {
	Organization    string `json:"organization"`
	CountryCode     string `json:"country_code"`
	ISP             string `json:"isp"`
	ASNOrganization string `json:"asn_organization"`
	ASN             int    `json:"asn"`
	Country         string `json:"country"`
	Region          string `json:"region"`
	City            string `json:"city"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
}

type IpsbSource struct {
	name string
}

func NewIpsbSource() *IpsbSource {
	return &IpsbSource{
		name: "ipsb",
	}
}

func (i *IpsbSource) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	url := fmt.Sprintf("%s/%s", IpsbAPIURL, query)
	client := common.GetHttpClient()
	
	// Create custom request with specific User-Agent for ipsb
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("ipsb API returned status %d", httpResp.StatusCode)
	}
	
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp IpsbResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ipsb response: %v, body: %s", err, string(body))
	}
	
	// Handle empty fields gracefully
	location := ""
	if apiResp.Country != "" {
		location = apiResp.Country
	}
	if apiResp.Region != "" {
		if location != "" {
			location += " " + apiResp.Region
		} else {
			location = apiResp.Region
		}
	}
	if apiResp.City != "" {
		if location != "" {
			location += " " + apiResp.City
		} else {
			location = apiResp.City
		}
	}
	if apiResp.ISP != "" {
		if location != "" {
			location += " " + apiResp.ISP
		} else {
			location = apiResp.ISP
		}
	}
	if apiResp.Organization != "" {
		if location != "" {
			location += " " + apiResp.Organization
		} else {
			location = apiResp.Organization
		}
	}

	result = &IpsbResult{Location: location}
	return result, nil
}

func (i *IpsbSource) Name() string {
	return i.name
}

type IpsbResult struct {
	Location string
}

func (r *IpsbResult) String() string {
	return r.Location
}