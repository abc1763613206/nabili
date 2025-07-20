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
	BaiduAPIURL = "https://qifu-api.baidubce.com/ip/geo/v1/district"
)

type BaiduResponse struct {
	Code interface{} `json:"code"`
	Data struct {
		Continent  string `json:"continent"`
		Country    string `json:"country"`
		Prov       string `json:"prov"`
		City       string `json:"city"`
		District   string `json:"district"`
		ISP        string `json:"isp"`
		IP         string `json:"ip"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type BaiduSource struct {
	name string
}

func NewBaiduSource() *BaiduSource {
	return &BaiduSource{
		name: "baidu",
	}
}

func (b *BaiduSource) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	url := fmt.Sprintf("%s?ip=%s", BaiduAPIURL, query)
	client := common.GetHttpClient()
	
	// Create custom request with specific User-Agent for baidu
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
		return nil, fmt.Errorf("baidu API returned status %d", httpResp.StatusCode)
	}
	
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp BaiduResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal baidu response: %v", err)
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
	
	if codeStr != "Success" && codeInt != 0 {
		return nil, fmt.Errorf("baidu API error: %s", apiResp.Msg)
	}

	// Handle empty fields gracefully
	location := ""
	if apiResp.Data.Country != "" {
		location = apiResp.Data.Country
	}
	if apiResp.Data.Prov != "" {
		if location != "" {
			location += " " + apiResp.Data.Prov
		} else {
			location = apiResp.Data.Prov
		}
	}
	if apiResp.Data.City != "" {
		if location != "" {
			location += " " + apiResp.Data.City
		} else {
			location = apiResp.Data.City
		}
	}
	if apiResp.Data.District != "" {
		if location != "" {
			location += " " + apiResp.Data.District
		} else {
			location = apiResp.Data.District
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

	result = &BaiduResult{Location: location}
	return result, nil
}

func (b *BaiduSource) Name() string {
	return b.name
}

type BaiduResult struct {
	Location string
}

func (r *BaiduResult) String() string {
	return r.Location
}