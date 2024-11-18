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

func dont_nothing() {
	fmt.Println("nothing.")
}

const LinksPath = "/ufmRestV2/resources/links"

func (u *UfmClient) LinksGetAll() (ret string, err error) {
	path := LinksPath
	resp, err := u.Get(path, []string{})
	fmt.Println("getting path:", path)
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("Error getting links data: " + resp.Status + " (" + string(bodyBytes) + ")")
		return
	}
	ret = string(bodyBytes)
	return

}
