package ufm

import (
	"errors"
	//"fmt"
	//"github.com/tidwall/gjson"
	//"github.com/tidwall/sjson"
	"io"
	//"strconv"
	//"strings"
	// "os"
)

const AlarmsPath = "/ufmRestV2/app/alarms"

// Get single alarm by id
func (u *UfmClient) AlarmsGet(id string) (ret string, err error) {
	path := AlarmsPath + "/" + id
	resp, err := u.Get(path, []string{})
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting alarms data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}

// Get all alarms, or all alarms for a device
func (u *UfmClient) AlarmsGetAll(device_id string) (ret string, err error) {
	var filtersArray []string
	if device_id == "" {
		filtersArray = []string{}
	} else {
		filtersArray = []string{device_id}
	}
	//fmt.Println("filtersarray: ", filtersArray)
	resp, err := u.Get(AlarmsPath, filtersArray)
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting alarms data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}
