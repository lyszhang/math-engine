/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 11:25 AM
 */

package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dengsgo/math-engine/common"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"io/ioutil"
	"net/http"
)

const FuncPrefix = "query"

type ExternalGravityPosition struct {
	AppID string
	Url   string
}

// QueryResponse collects the response values for the Query method.
type QueryResponse struct {
	Rs  []byte
	Pub paillier.PublicKey
	Err error
}

func FetchExternalGravity(pos *ExternalGravityPosition, key string) ([]byte, *paillier.PublicKey, error) {
	body := fmt.Sprintf("{\"S\": \"%s\"}", key)
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

// UploadResponse collects the response values for the Query method.
type UploadResponse struct {
	Rs  int64
	Err error
}

func UploadResult(r []byte) (int64, error) {
	req := UploadRequest{R: r}
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

	var resp UploadResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	return resp.Rs, nil
}

// CompareRequest collects the request parameters for the Query method.
type CompareRequest struct {
	A []byte
	B []byte
}

// CompareResponse collects the response values for the Query method.
type CompareResponse struct {
	Rs  int64
	Err error
}

func CompareResult(a, b []byte) (int64, error) {
	req := CompareRequest{A: a, B: b}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	res, err := http.Post("http://127.0.0.1:5666/compare", "application/json;charset=utf-8",
		bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	var resp CompareResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	return resp.Rs, nil
}

// TransformRequest collects the request parameters for the Query method.
type TransformRequest struct {
	Data common.CipherCompression
}

// TransformResponse collects the response values for the Query method.
type TransformResponse struct {
	Data common.CipherCompression
	Err  error
}

func TransformExternal(d *common.CipherCompression) (*common.CipherCompression, error) {
	req := TransformRequest{Data: *d}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	res, err := http.Post("http://127.0.0.1:5666/transform", "application/json;charset=utf-8",
		bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	var resp TransformResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	return &resp.Data, nil
}
