package ufm

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"strconv"
	"strings"
	// "os"
)

func don_nothing() {
	fmt.Println("nothing.")
}

const PortsPath = "/ufmRestV2/resources/ports"
const ActionsPath = "/ufmRestV2/actions"

func (u *UfmClient) PortsGet(portName string) (ret string, err error) {
	path := PortsPath + "/" + portName
	resp, err := u.Get(path, []string{})
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

func (u *UfmClient) PortsGetAllBrief(filters string) (ret string, err error) {
	ret = "[]"
	if filters == "" {
		filters = `show_disabled=true`
	} else {
		filters = `show_disabled=true,` + filters
	}
	allPorts, err := u.PortsGetAll(filters)
	if err != nil {
		return
	}
	for i, port := range gjson.Parse(allPorts).Array() {
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".name", port.Get("name").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".guid", port.Get("guid").String())
		if err != nil {
			return
		}
		//path := strings.TrimRight(strings.TrimLeft(strings.Split(port.Get("path").String(), "/")[1], " "), " ")
		//ret, err = sjson.Set(ret, strconv.Itoa(i)+".path", path)
		//if err != nil {
		//	return
		//}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".system_name", port.Get("system_name").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".node_description", port.Get("node_description").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".logical_state", port.Get("logical_state").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".physical_state", port.Get("physical_state").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".enabled_speed", port.Get("enabled_speed").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".peer_node_description", port.Get("peer_node_description").String())
		if err != nil {
			return
		}
	}
	return

}

func (u *UfmClient) PortsGetAll(filters string) (ret string, err error) {
	var filtersArray []string
	if filters == "" {
		filtersArray = []string{}
	} else {
		filtersArray = strings.Split(filters, ",")
	}
	//fmt.Println("filtersarray: ", filtersArray)
	resp, err := u.Get(PortsPath, filtersArray)
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

// For some reason, there is a separate api path for actions where ports get reset
//func (u *UfmClient) PortsAction(guid string, action string) (ret string, err error) {
//	path := ActionsPath
//	resp, err := u.Put(path, []string{})
//	if err != nil {
//		return
//	}
//	// assume name is guid_1 ?
//	payload := "{}"
//	jsonData := sjson.Set(payload, ".params.port_id", guid+"_1")
//	jsonData := sjson.Set(payload, ".action", action)
//	jsonData := sjson.Set(payload, ".object_ids", []string{"system_guid"})
//	jsonData := sjson.Set(payload, ".object_type", "System")
//	jsonData := sjson.Set(payload, ".description", action+guid)
//	jsonData := sjson.Set(payload, ".identifier", ??)
//
//	bodyBytes, _ := io.ReadAll(resp.Body)
//	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
//		err = errors.New("Error getting switchport data: " + resp.Status + " (" + string(bodyBytes) + ")")
//		return
//	}
//	ret = string(bodyBytes)
//	return
//
//}
