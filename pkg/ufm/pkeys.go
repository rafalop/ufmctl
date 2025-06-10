package ufm

import (
//	"bytes"
	"errors"
	"io"
	"strings"
	//	"os"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"path/filepath"
)

func donnothing() {
	fmt.Println("nothing.")
	a, _ := sjson.Set("[]", "", interface{}([]string{}))
	b := gjson.Get(a, "")
	fmt.Println("nothing.", a, b, errors.New("nothing"))
}

const PkeysPath = "/ufmRestV2/resources/pkeys"

func (u *UfmClient) PkeyList() (ret string, err error) {
	// TODO: For some reason, you can't retrieve both guids_data and qos_conf together
	// may have had something to do with SHARP aggregation manager was not running ?
	//queries := []string{"guids_data=true", "qos_conf=true", "port_info=true"}
	queries := []string{"guids_data=true", "port_info=true"}

	// get guid data
	resp, err := u.Get(PkeysPath, queries)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	//fmt.Println("resp:", resp.StatusCode, "body: ", string(bodyBytes))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("there was an error getting pkeys data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}

func (u *UfmClient) PkeyGet(pkey string, guids_data string) (ret string, err error) {
	// get guid data
	//fmt.Println(pkey)
	resp, err := u.Get(filepath.Join(PkeysPath, pkey), []string{"guids_data=" + guids_data})
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
	fmt.Println("Posting data: ", string(jsonData))

	// This broke other pkeys when testing adding a new pkey!
	//resp, err := u.Post(PkeysPath+"/add", bytes.NewBuffer(jsonData))
	//if resp.StatusCode != 201 {
	//	bodyBytes, _ := io.ReadAll(resp.Body)
	//	return errors.New("there was an error creating the pkey: " + resp.Status + " (" + string(bodyBytes) + ")")
	//}
	return nil
}

// Set (overwrite!) set of guids, index0, ipoib for a pkey
func (u *UfmClient) PkeySetGuids(pkey string, index0 bool, ipoib bool, membership string, guids []string) (err error) {
	var memberships []string
	for range guids {
		memberships = append(memberships, membership)
	}
	newGroup := map[string]interface{}{
		"pkey":        pkey,
		"index0":      index0,
		"ip_over_ib":  ipoib,
		"guids":       guids,
		"memberships": memberships,
	}
	resultJSON, _ := sjson.Set("{}", "-1", newGroup)
	fmt.Println(resultJSON)
	// insert code here to PUT
	return
}

func (u *UfmClient) PkeyAddGuids(pkey string, index0 bool, ipoib bool, membership string, guids []string) (err error) {
	//// START GARBAGE
	//// This hopefully won't have to be used... get current data, update it and push back to UFM.
	//// But this is how browser UI works to add/remove keys :(
	//// Get current pkey data for this pkey
	//currentPkeyData, err := u.PkeyGet(pkey, "true")
	//if err != nil {
	//	return
	//}
	////fmt.Println(currentPkeyData)
	//// Extract `guids` array
	//guidsArray := gjson.Get(currentPkeyData, "guids").Array()

	//// Prepare the result JSON string
	//resultJSON := `[]`

	//// Function to find if a group with the same properties exists
	//findGroup := func(pkey string, index0, ipOverIB bool) (int, bool) {
	//	for i, group := range gjson.Parse(resultJSON).Array() {
	//		if group.Get("pkey").String() == pkey && group.Get("index0").Bool() == index0 && group.Get("ip_over_ib").Bool() == ipOverIB {
	//			return i, true
	//		}
	//	}
	//	return -1, false
	//}

	//// Iterate over the `guids` array and group by `pkey`, `index0`, and `ip_over_ib`
	//for _, item := range guidsArray {
	//	//pkey := item.Get("pkey").String()
	//	index0 := item.Get("index0").Bool()
	//	ipOverIB := item.Get("ip_over_ib").Bool()
	//	guid := item.Get("guid").String()
	//	membership := item.Get("membership").String()

	//	// Find or create the group
	//	if idx, exists := findGroup(pkey, index0, ipOverIB); exists {
	//		// Append to the existing group
	//		resultJSON, _ = sjson.Set(resultJSON, fmt.Sprintf("%d.guids.-1", idx), guid)
	//		resultJSON, _ = sjson.Set(resultJSON, fmt.Sprintf("%d.memberships.-1", idx), membership)
	//	} else {
	//		// Create a new group
	//		newGroup := map[string]interface{}{
	//			"pkey":        pkey,
	//			"index0":      index0,
	//			"ip_over_ib":  ipOverIB,
	//			"guids":       []string{guid},
	//			"memberships": []string{membership},
	//		}
	//		resultJSON, _ = sjson.Set(resultJSON, "-1", newGroup)
	//	}
	//}

	//// add the guids passed to this func
	////fmt.Println(guids)
	//for _, g := range guids {
	//	sjson.Set(resultJSON, "-1", g)
	//}
	// END garbage
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
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "ip_over_ib", ipoib)
	if err != nil {
		return
	}
	pkeyAddGuidsData, err = sjson.Set(pkeyAddGuidsData, "membership", membership)
	if err != nil {
		return
	}
	//fmt.Println(pkeyAddGuidsData)
	resp, err := u.Post(PkeysPath, "", pkeyAddGuidsData)
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = errors.New("there was an error creating the pkey: " + resp.Status + " (" + string(bodyBytes) + ")")
	}
	return
}

func (u *UfmClient) PkeyRemoveGuids(pkey string, guids []string) (err error) {
	deleteGuidsPath := PkeysPath + "/" + pkey + "/guids/" + strings.Join(guids, ",")
	//fmt.Println("path:", deleteGuidsPath)
	resp, err := u.Delete(deleteGuidsPath)
	//bodyBytes, _ := io.ReadAll(resp.Body)
	//fmt.Println("resp:", resp.Status, "body:", string(bodyBytes))
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = errors.New("there was an error deleting GUIDs from the pkey: " + resp.Status + " (" + string(bodyBytes) + ")")
	}
	return

}
