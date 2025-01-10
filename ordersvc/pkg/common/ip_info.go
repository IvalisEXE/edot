package common

import (
	"encoding/json"
	"net/http"
)

type IpInfo struct {
	IP      string `json:"query"`
	Region  string `json:"regionName"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func GetIpInfo(ip string) (*IpInfo, error) {
	// Buat request ke API
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return &IpInfo{
			IP: ip,
		}, err
	}
	defer resp.Body.Close()

	var ipInfo IpInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, err
	}

	return &ipInfo, nil
}
