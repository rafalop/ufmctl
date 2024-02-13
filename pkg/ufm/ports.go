package ufm

import (
	"errors"
	"io"
//	"os"
	"encoding/json"
	"fmt"
)

func don_nothing(){
	fmt.Println("nothing.")
}


// Port represents the structure for the Switch IB Port data.
type Port struct {
	Description         string   `json:"description"`
	Number              int      `json:"number"`
	ExternalNumber      int      `json:"external_number"`
	PhysicalState       string   `json:"physical_state"`
	Path                string   `json:"path"`
	Tier                int      `json:"tier"`
	LID                 int      `json:"lid"`
	Mirror              string   `json:"mirror"`
	LogicalState        string   `json:"logical_state"`
	Capabilities        []string `json:"capabilities"`
	MTU                 int      `json:"mtu"`
	PeerPortDName       string   `json:"peer_port_dname"`
	Severity            string   `json:"severity"`
	ActiveSpeed         string   `json:"active_speed"`
	EnabledSpeed        []string `json:"enabled_speed"`
	SupportedSpeed      []string `json:"supported_speed"`
	ActiveWidth         string   `json:"active_width"`
	EnabledWidth        []string `json:"enabled_width"`
	SupportedWidth      []string `json:"supported_width"`
	DName               string   `json:"dname"`
	PeerNodeName        string   `json:"peer_node_name"`
	Peer                string   `json:"peer"`
	PeerNodeGUID        string   `json:"peer_node_guid"`
	SystemID            string   `json:"systemID"`
	NodeDescription     string   `json:"node_description"`
	Name                string   `json:"name"`
	Module              string   `json:"module"`
	PeerLID             int      `json:"peer_lid"`
	PeerGUID            string   `json:"peer_guid"`
	PeerNodeDescription string   `json:"peer_node_description"`
	GUID                string   `json:"guid"`
}

type PortBrief struct {
	Name                string   `json:"name"`
	Path                string   `json:"path"`
	LogicalState        string   `json:"logical_state"`
	PhysicalState       string   `json:"physical_state"`
}

func (u *UfmClient) GetPort(portName string) (*Port, error) {
	path := "/ufmRestV2/resources/ports/"+portName
	resp, err := u.Get(path, []string{})
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		return nil, errors.New("Error getting switchport data: "+resp.Status+" ("+string(bodyBytes)+")")
	}
	ports := []Port{}
	json.Unmarshal(bodyBytes, &ports)	
	return &ports[0], nil
	
}


func (u *UfmClient) GetPortsBrief(queries ...string) ([]PortBrief, error) {
	path := "/ufmRestV2/resources/ports"
	resp, err := u.Get(path, queries)
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		return nil, errors.New("Error getting switchport data: "+resp.Status+" ("+string(bodyBytes)+")")
	}
	allPorts := []PortBrief{}
	json.Unmarshal(bodyBytes, &allPorts)	
	return allPorts, nil
	
}

func (u *UfmClient) GetPortsFull(queries ...string) ([]Port, error) {
	path := "/ufmRestV2/resources/ports"
	resp, err := u.Get(path, queries)
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		return nil, errors.New("Error getting switchport data: "+resp.Status+" ("+string(bodyBytes)+")")
	}
	allPorts := []Port{}
	json.Unmarshal(bodyBytes, &allPorts)	
	return allPorts, nil
	
}
