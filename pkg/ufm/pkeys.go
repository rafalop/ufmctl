package ufm

import (
	"bytes"
	"errors"
	"io"
	"strings"
	//	"os"
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
)

func donnothing() {
	fmt.Println("nothing.")
}

const PkeysPath = "/ufmRestV2/resources/pkeys"

func (u *UfmClient) PkeyList() (ret string, err error) {
	// TODO: For some reason, you can't retrieve both guids_data and qos_conf together
	//queries := []string{"guids_data=true", "qos_conf=true", "port_info=true"}
	queries := []string{"guids_data=true", "port_info=true"}

	// get guid data
	resp, err := u.Get(PkeysPath, queries)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ret = string(bodyBytes)
	return

}

const (
	MembershipLimited = "limited"
	MembershipFull    = "full"
)

type CreatePkeyData struct {
	Pkey         string  `json:"pkey"`
	Index0       string  `json:"index0,omitempty"`
	IpOverIb     bool    `json:"ip_over_ib,omitempty"`
	MtuLimit     int     `json:"mtu_limit,omitempty"`
	ServiceLevel int     `json:"service_level,omitempty"`
	RateLimit    float64 `json:"rate_limit,omitempty"`
}

func (u *UfmClient) CreatePkey(data *CreatePkeyData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//fmt.Println("submitting: ", string(jsonData))
	resp, err := u.Post(PkeysPath+"/add", bytes.NewBuffer(jsonData))
	if resp.StatusCode != 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New("there was an error creating the pkey: " + resp.Status + " (" + string(bodyBytes) + ")")
	}
	return nil
}

//type PkeyAddGuidsData struct {
//    Guids       []string `json:"guids,omitempty"`
//    IpOverIB    bool     `json:"ip_over_ib,omitempty"`
//    Index0      bool     `json:"index0,omitempty"`
//    Membership  string   `json:"membership,omitempty"`
//    Pkey        string   `json:"pkey,omitempty"`
//}

//var pkeyAddGuidsData string = `
//{
//"guids": [],
//"ip_over_ib": false,
//"index0": true,
//"membership": "limited",
//"pkey": "0x0a12"
//}`

func (u *UfmClient) PkeyAddGuids(pkey string, index0 bool, ipoib bool, membership string, guids []string) (err error) {
	pkeyAddGuidsData := "{}"
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "pkey", pkey)
	if err != nil {
		return
	}
	for _, guid := range guids {
		pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "guids.-1", guid)
		if err != nil {
			return
		}
	}
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "index0", index0)
	if err != nil {
		return
	}
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "ipoib", ipoib)
	if err != nil {
		return
	}
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "membership", membership)
	if err != nil {
		return
	}
	fmt.Println("data:", pkeyAddGuidsData)
	//resp, err = u.Post(PkeysPath, bytes.NewReader([]bytes(pkeyAddGuidsData)))
	//if resp.StatusCode != 200 {
	//	bodyBytes, _ := io.ReadAll(resp.Body)
	//	err = errors.New("there was an error creating the pkey: "+resp.Status+" ("+string(bodyBytes)+")")
	//}
	return
}

func (u *UfmClient) PkeyRemoveGuids(pkey string, guids []string) (err error) {
	pkeyRemoveGuidsData := "{}"
	pkeyRemoveGuidsData, err = sjson.Set(pkeyRemoveGuidsData, "pkey", pkey)
	if err != nil {
		return
	}

	for _, guid := range guids {
		pkeyRemoveGuidsData, err = sjson.Set(pkeyRemoveGuidsData, "guids.-1", guid)
		if err != nil {
			return
		}
	}
	path := PkeysPath + "/" + pkey + "/guids/" + strings.Join(guids, ",")
	fmt.Println("data:", pkeyRemoveGuidsData)
	fmt.Println("path:", path)

	//resp, err = u.Remove(path, bytes.NewReader([]bytes(pkeyRemoveGuidsData)))
	//if resp.StatusCode != 200 {
	//	bodyBytes, _ := io.ReadAll(resp.Body)
	//	err = errors.New("there was an error creating the pkey: "+resp.Status+" ("+string(bodyBytes)+")")
	//}
	return
}
