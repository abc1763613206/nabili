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
	Code int `json:"code"`
	Data struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		ISP      string `json:"isp"`
		IP       string `json:"ip"`
	} `json:"data"`
	Message string `json:"message"`
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

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("iqiyi API error: %s", apiResp.Message)
	}

	// Handle empty fields gracefully
	location := ""
	if apiResp.Data.Country != "" {
		location = apiResp.Data.Country
	}
	if apiResp.Data.Province != "" {
		if location != "" {
			location += " " + apiResp.Data.Province
		} else {
			location = apiResp.Data.Province
		}
	}
	if apiResp.Data.City != "" {
		if location != "" {
			location += " " + apiResp.Data.City
		} else {
			location = apiResp.Data.City
		}
	}
	if apiResp.Data.ISP != "" {
		if location != "" {
			location += " " + apiResp.Data.ISP
		} else {
			location = apiResp.Data.ISP
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