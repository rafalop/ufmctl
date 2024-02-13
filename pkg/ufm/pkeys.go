package ufm

import (
	"bytes"
	"errors"
	"io"
//	"os"
	"encoding/json"
	"fmt"
)

func donnothing(){
	fmt.Println("nothing.")
}

const PkeysPath = "/ufmRestV2/resources/pkeys"


func (u *UfmClient) GetPkeys() (ret string, err error) {
		// For some reason, you can't retrieve both guids_data and qos_conf together
		//queries := []string{"guids_data=true", "qos_conf=true", "port_info=true"}
		queries := []string{"guids_data=true", "port_info=true"} 
		
		// get keys alone
		var keys []string
		resp, err := u.Get(PkeysPath, []string{}) 
		if err != nil {
			return
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}	
		json.Unmarshal(bodyBytes, &keys)

		// get guid data
		resp, err = u.Get(PkeysPath, queries)
		if err != nil {
			return
		}
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		ret = string(bodyBytes)
		return
		
}
const (
	MembershipLimited = "limited"
	MembershipFull = "full"
)


type CreatePkeyData struct {
	Pkey string `json:"pkey"`
	Index0 string `json:"index0,omitempty"`
	IpOverIb bool	`json:"ip_over_ib,omitempty"`
	MtuLimit int	`json:"mtu_limit,omitempty"`
	ServiceLevel int	`json:"service_level,omitempty"`
	RateLimit float64	`json:"rate_limit,omitempty"`
}


func (u *UfmClient) CreatePkey(data *CreatePkeyData) (error){
	jsonData, err := json.Marshal(data)	
	if err != nil {
		return err
	}
	//fmt.Println("submitting: ", string(jsonData))
	resp, err := u.Post(PkeysPath+"/add", bytes.NewBuffer(jsonData))
	if resp.StatusCode != 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New("there was an error creating the pkey: "+resp.Status+" ("+string(bodyBytes)+")")
	}
	return nil
}
