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

const EventsPath = "/ufmRestV2/app/events"

// Get single alarm by id
func (u *UfmClient) EventsGet(id string) (ret string, err error) {
	path := EventsPath + "/" + id
	resp, err := u.Get(path, []string{})
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting events data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}

// Get all alarms, or all alarms for a device
func (u *UfmClient) EventsGetAll() (ret string, err error) {
	//fmt.Println("filtersarray: ", filtersArray)
	resp, err := u.Get(EventsPath, []string{})
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting events data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}
