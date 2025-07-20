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
	IqiyiAPIURL = "https://mesh.if.iqiyi.com/aid/ip/info"
)

type IqiyiResponse struct {
	Code interface{} `json:"code"`
	Data struct {
		CountryCN  string `json:"countryCN"`
		ProvinceCN string `json:"provinceCN"`
		CityCN     string `json:"cityCN"`
		ISPCN      string `json:"ispCN"`
		IP         string `json:"ip"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type IqiyiSource struct {
	name string
}

func NewIqiyiSource() *IqiyiSource {
	return &IqiyiSource{
		name: "iqiyi",
	}
}

func (i *IqiyiSource) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	url := fmt.Sprintf("%s?version=1.1.1&ip=%s", IqiyiAPIURL, query)
	client := common.GetHttpClient()
	
	// Create custom request with specific User-Agent for iqiyi
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
		return nil, fmt.Errorf("iqiyi API returned status %d", httpResp.StatusCode)
	}
	
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp IqiyiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal iqiyi response: %v", err)
	}

	// Handle both string and numeric response codes
	var codeStr string
	var codeInt int
	switch v := apiResp.Code.(type) {
	case string:
		codeStr = v
	case float64:
		codeInt = int(v)
	case int:
		codeInt = v
	default:
		return nil, fmt.Errorf("invalid response code type: %T", apiResp.Code)
	}
	
	if codeStr != "0" && codeInt != 0 {
		return nil, fmt.Errorf("iqiyi API error: %s", apiResp.Msg)
	}

	// Handle empty fields gracefully and skip * values
	location := ""
	if apiResp.Data.CountryCN != "" && apiResp.Data.CountryCN != "*" {
		location = apiResp.Data.CountryCN
	}
	if apiResp.Data.ProvinceCN != "" && apiResp.Data.ProvinceCN != "*" {
		if location != "" {
			location += " " + apiResp.Data.ProvinceCN
		} else {
			location = apiResp.Data.ProvinceCN
		}
	}
	if apiResp.Data.CityCN != "" && apiResp.Data.CityCN != "*" {
		if location != "" {
			location += " " + apiResp.Data.CityCN
		} else {
			location = apiResp.Data.CityCN
		}
	}
	if apiResp.Data.ISPCN != "" && apiResp.Data.ISPCN != "*" {
		if location != "" {
			location += " " + apiResp.Data.ISPCN
		} else {
			location = apiResp.Data.ISPCN
		}
	}

	// Fallback to IP if nothing else
	if location == "" {
		location = apiResp.Data.IP
	}

	result = &IqiyiResult{Location: location}
	return result, nil
}

func (i *IqiyiSource) Name() string {
	return i.name
}

type IqiyiResult struct {
	Location string
}

func (r *IqiyiResult) String() string {
	return r.Location
}