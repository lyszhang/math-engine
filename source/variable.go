/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 11:25 AM
 */

package source

import (
	"bytes"
	"encoding/hex"
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
	Rs     []byte
	Offset int64
	Pub    paillier.PublicKey
	Err    error
}

func GetExternalGravity(key string) (*common.ArithmeticFactor, error) {
	if value, ok := common.Cache.Get(key); ok {
		if f, ok := value.(common.ArithmeticFactor); ok {
			return &f, nil
		}
	}
	return fetchExternalGravity(nil, key)
}

// 优先从本地cache中寻找，没有的话才会从远端请求
func fetchExternalGravity(pos *ExternalGravityPosition, key string) (*common.ArithmeticFactor, error) {
	fmt.Printf("\033[1;31;40m 尝试从数据端取数据: %s\033[0m\n", key)
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
	fmt.Printf("\033[1;31;40m 接收到键值: %s 的密文 %s, offset: %d\033[0m\n", key, hex.EncodeToString(resp.Rs), resp.Offset)
	return &common.ArithmeticFactor{
		Factor: common.TypePaillier,
		Cipher: common.NumberEncrypted{
			Data:      resp.Rs,
			PublicKey: &resp.Pub,
		},
		Offset: resp.Offset,
	}, nil
}

// UploadRequest collects the request parameters for the Query method.
type UploadRequest struct {
	R []byte
}

// UploadResponse collects the response values for the Query method.
type UploadResponse struct {
	Rs  float64
	Err error
}

func UploadResult(r []byte) (float64, error) {
	req := UploadRequest{R: r}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	fmt.Printf("\033[1;31;40m 计算完成，回传密文结果: %s\033[0m\n", hex.EncodeToString(r))
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
