package ufm

import (
	"errors"
	"fmt"
	//"github.com/tidwall/gjson"
	//"github.com/tidwall/sjson"
	"io"
	//"strconv"
	//"strings"
	// "os"
)

func don_nothing() {
	fmt.Println("nothing.")
}

const VPortsPath = "/ufmRestV2/resources/vports"

//const ActionsPath = "/ufmRestV2/actions"

//func (u *UfmClient) VPortsGet(portName string) (ret string, err error) {
//	path := VPortsPath + "/" + portName
//	resp, err := u.Get(path, []string{})
//	if err != nil {
//		return
//	}
//	bodyBytes, _ := io.ReadAll(resp.Body)
//	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
//		err = errors.New("Error getting switchport data: " + resp.Status + " (" + string(bodyBytes) + ")")
//		return
//	}
//	ret = string(bodyBytes)
//	return
//
//}

//func (u *UfmClient) VPortsGetAllBrief(filters string) (ret string, err error) {
//	ret = "[]"
//	allPorts, err := u.PortsGetAll(filters)
//	if err != nil {
//		return
//	}
//	for i, port := range gjson.Parse(allPorts).Array() {
//		ret, err = sjson.Set(ret, strconv.Itoa(i)+".name", port.Get("name").String())
//		if err != nil {
//			return
//		}
//		ret, err = sjson.Set(ret, strconv.Itoa(i)+".path", port.Get("path").String())
//		if err != nil {
//			return
//		}
//		ret, err = sjson.Set(ret, strconv.Itoa(i)+".logical_state", port.Get("logical_state").String())
//		if err != nil {
//			return
//		}
//		ret, err = sjson.Set(ret, strconv.Itoa(i)+".physical_state", port.Get("physical_state").String())
//		if err != nil {
//			return
//		}
//	}
//	return
//
//}

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

// For some reason, there is a separate api path for actions where ports get reset
//func (u *UfmClient) PortsAction(guid string, action string) (ret string, err error) {
//	path := ActionsPath
//	//{
//	//  "params": {
//	//    "port_id": "port_name"
//	//  },
//	//  "action": "enable,disable,reset",
//	//  "object_ids": [
//	//    "system_guid"
//	//  ],
//	//  "object_type": "System",
//	//  "description": " description",
//	//  "identifier": "id"
//	//}
//	// assume name is guid_1 ?
//	payload := "{}"
//	payload,_ = sjson.Set(payload, ".params.port_id", guid+"_1")
//	payload,_ = sjson.Set(payload, ".action", action)
//	payload,_ = sjson.Set(payload, ".object_ids", []string{"system_guid"})
//	payload,_ = sjson.Set(payload, ".object_type", "System")
//	payload,_ = sjson.Set(payload, ".description", action+guid)
//	payload,_ = sjson.Set(payload, ".identifier", "dummy")
//
//	fmt.Println("Payload:", payload)
//	fmt.Println("Path:", path)
//	//resp, err := u.Put(path, bytes.NewReader([]byte(payload))
//	//if err != nil {
//	//	return
//	//}
//	//bodyBytes, _ := io.ReadAll(resp.Body)
//	//if resp.StatusCode < 200 || resp.StatusCode >= 300 {
//	//	err = errors.New("Error getting switchport data: " + resp.Status + " (" + string(bodyBytes) + ")")
//	//	return
//	//}
//	//ret = string(bodyBytes)
//	return
//
//}
