package ufm

import (
	"errors"
	"io"
	"github.com/tidwall/gjson"
	"fmt"
	"github.com/tidwall/sjson"
	"strconv"
	"strings"
//	"os"
)

func don_nothing(){
	fmt.Println("nothing.")
}

const PortsPath = "/ufmRestV2/resources/ports"


func (u *UfmClient) PortsGet(portName string) (ret string, err error) {
	path := PortsPath+"/"+portName
	resp, err := u.Get(path, []string{})
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		err = errors.New("Error getting switchport data: "+resp.Status+" ("+string(bodyBytes)+")")
		return
	}
	ret = string(bodyBytes)
	return
	
}


func (u *UfmClient) PortsGetAllBrief(filters string) (ret string, err error) {
	ret = "[]"
	allPorts, err := u.PortsGetAll(filters)
	if err != nil {
		return
	}
	for i, port := range gjson.Parse(allPorts).Array() {
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".name", port.Get("name").String())
		if err != nil {
			return
		}
		ret, err = sjson.Set(ret, strconv.Itoa(i)+".path", port.Get("path").String())
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
	resp, err := u.Get(PortsPath, filtersArray)
	if err != nil {
		return 
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		err = errors.New("Error getting switchport data: "+resp.Status+" ("+string(bodyBytes)+")")
		return
	}
	ret = string(bodyBytes)
	return
	
}
