package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/abc1763613206/nabili/pkg/common"
)

const (
	BiliAPIURL = "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr"
)

type BiliResponse struct {
	Code int `json:"code"`
	Data struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		ISP      string `json:"isp"`
		Addr     string `json:"addr"`
	} `json:"data"`
	Message string `json:"message"`
}

type BiliSource struct {
	name string
}

func NewBiliSource() *BiliSource {
	return &BiliSource{
		name: "bili",
	}
}

func (b *BiliSource) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	url := fmt.Sprintf("%s?ip=%s", BiliAPIURL, query)
	client := common.GetHttpClient()
	
	// Create custom request with specific User-Agent for bili
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
		return nil, fmt.Errorf("bili API returned status %d", httpResp.StatusCode)
	}
	
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp BiliResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	result = &BiliResult{
		Country:  resp.Data.Country,
		Province: resp.Data.Province,
		City:     resp.Data.City,
		ISP:      resp.Data.ISP,
		Addr:     resp.Data.Addr,
	}

	return result, nil
}

func (b *BiliSource) Name() string {
	return b.name
}

type BiliResult struct {
	Country  string
	Province string
	City     string
	ISP      string
	Addr     string
}

func (r *BiliResult) String() string {
	parts := []string{}
	if r.Country != "" && r.Country != "N/A" {
		parts = append(parts, r.Country)
	}
	if r.Province != "" && r.Province != "N/A" {
		parts = append(parts, r.Province)
	}
	if r.City != "" && r.City != "N/A" {
		parts = append(parts, r.City)
	}
	if r.ISP != "" && r.ISP != "N/A" {
		parts = append(parts, r.ISP)
	}
	
	if len(parts) == 0 && r.Addr != "" {
		return r.Addr
	}
	
	return strings.Join(parts, " ")
}