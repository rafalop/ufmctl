package ufm

import (
	//"bytes"
	//"errors"
	"io"
	//	"os"
	//"encoding/json"
	"fmt"
)

func systemsdonnothing() {
	fmt.Println("nothing.")
}

const SystemsPath = "/ufmRestV2/resources/systems"

// get all systems, pass in list of filters such as [ip=1.1.1.1, name=mycomputer,...]
func (u *UfmClient) GetSystems(filters []string) (ret string, err error) {
	resp, err := u.Get(SystemsPath, filters)
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
