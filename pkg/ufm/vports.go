package ufm

import (
	"errors"
	"io"
)

const VPortsPath = "/ufmRestV2/resources/vports"

func (u *UfmClient) VPortsGetAll(physPort string) (ret string, err error) {
	var filtersArray []string
	if physPort == "" {
		filtersArray = []string{}
	} else {
		filtersArray = []string{"port=" + physPort}
	}
	resp, err := u.Get(VPortsPath, filtersArray)
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting switchport data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}
