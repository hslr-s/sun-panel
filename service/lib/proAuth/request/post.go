package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sun-panel/lib/proAuth/AES"
)

type RequestParam struct {
	P string
}

type RequestResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

type RequestRespDebug struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// 发送POST请求
func SendPostRequest(url string, contentType string, requestData []byte) ([]byte, *http.Response, error) {
	// 发起 POST 请求
	resp, err := http.Post(url, contentType, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, resp, fmt.Errorf("POST:Request error: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("POST:Response status code error: %d", resp.StatusCode)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("POST:Error reading response body: %v", err)
	}

	return responseBody, resp, nil
}

// 加密请求
func EncryptRequestData(AESKey string, param interface{}) ([]byte, error) {
	if byt, err := json.Marshal(param); err != nil {
		return nil, errors.New("json error:" + err.Error())
	} else {
		// fmt.Println("请求数据", string(byt))
		if encrypt, err := AES.Encrypt(AESKey, string(byt)); err != nil {
			return nil, errors.New("encrypt error:" + err.Error())
		} else {
			return []byte(encrypt), nil
		}
	}
}

// 解密响应
func DecryptResponseData(AESKey string, response string, value interface{}) error {
	if decryptStr, err := AES.Decrypt(AESKey, response); err != nil {
		return errors.New("decrypt error:" + err.Error())
	} else {
		// fmt.Println("接口返回Data数据", decryptStr)
		if err := json.Unmarshal([]byte(decryptStr), value); err != nil {
			return errors.New("JSON encrypted string error:" + err.Error())
		}
	}
	return nil
}

func SendPostJsonAESKeyRequest(url, AESKey string, requestData interface{}, responseData interface{}) (int, error) {

	// reqByte, _ := json.Marshal(requestData)
	// fmt.Println("接口返回请求数据", string(reqByte))

	reqBytes, err := EncryptRequestData(AESKey, requestData)
	if err != nil {
		return 0, errors.New("Encrypt request fail:" + err.Error())
	}

	req := RequestParam{
		P: string(reqBytes),
	}

	reqString, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	respContent, _, err := SendPostRequest(url, "application/json", []byte(reqString))
	if err != nil {
		return 0, err
	}
	// fmt.Println("接口返回总数据", string(respContent))
	requestResp := RequestResp{}
	if err := json.Unmarshal(respContent, &requestResp); err != nil {
		return requestResp.Code, errors.New("requestResp error :" + err.Error())
	}

	// fmt.Println("接口返回Data数据", requestResp.Data)
	if requestResp.Code != 0 {
		return requestResp.Code, errors.New(requestResp.Msg)
	}

	return requestResp.Code, DecryptResponseData(AESKey, requestResp.Data, responseData)
}

func SendPostJsonRequest(url string, requestData interface{}, responseData interface{}) (int, error) {

	reqByte, err := json.Marshal(requestData)
	// fmt.Println("请求数据", string(reqByte))
	if err != nil {
		return 0, err
	}

	respContent, _, err := SendPostRequest(url, "application/json", reqByte)
	if err != nil {
		return 0, err
	}
	// fmt.Println("接口返回总数据", string(respContent))

	// 第一轮解析
	requestResp := RequestRespDebug{}
	if err := json.Unmarshal(respContent, &requestResp); err != nil {
		return requestResp.Code, errors.New("requestResp error :" + err.Error())
	}

	if requestResp.Code != 0 {
		return requestResp.Code, errors.New(requestResp.Msg)
	}

	requestRespData, _ := json.Marshal(requestResp.Data)
	// fmt.Println("data数据", string(requestRespData))
	// 第二轮解析
	if err := json.Unmarshal(requestRespData, &responseData); err != nil {
		return requestResp.Code, errors.New("requestResp data error :" + err.Error())
	}

	return 0, nil
}
