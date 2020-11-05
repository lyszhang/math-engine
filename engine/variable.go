/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 11:25 AM
 */

package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"io/ioutil"
	"net/http"
)

const FuncPrefix = "query"

type ExternalGravityPosition struct {
	AppID string
	Url string
}

// QueryResponse collects the response values for the Query method.
type QueryResponse struct {
	Rs  []byte
	Pub paillier.PublicKey
	Err error
}


func FetchExternalGravity(pos *ExternalGravityPosition) ([]byte, *paillier.PublicKey, error){
	body := "{\"S\": \"newuser\"}"
	res, err := http.Post("http://127.0.0.1:5666/query", "application/json;charset=utf-8",
		bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	var resp QueryResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	return resp.Rs, &resp.Pub, nil
}

// UploadRequest collects the request parameters for the Query method.
type UploadRequest struct {
	R []byte
}
func UploadResult(r []byte) ([]byte, error) {
	req := UploadRequest{R:r}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	res, err := http.Post("http://127.0.0.1:5666/upload", "application/json;charset=utf-8",
		bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	fmt.Println(string(content))
	return content, nil
}
