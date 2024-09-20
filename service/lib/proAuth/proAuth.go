package proAuth

import (
	"errors"
	"fmt"
	"math"
	"sun-panel/lib/proAuth/request"
	"time"

	"github.com/google/uuid"
)

type ConfigType struct {
	Url            string
	AESKey         string
	IsDebug        bool
	TimeDifference int64
}

var Config = ConfigType{
	Url:            "http://127.0.0.1:3005",
	AESKey:         "",
	IsDebug:        false,
	TimeDifference: 120,
}

func Register(req RegisterReq) (registerResp RegisterResp, err error) {
	req.Timestamp = time.Now().Unix()
	req.Code = fmt.Sprintf("%s%d", uuid.NewString(), time.Now().Unix())
	if _, err := requestCommon(Config.Url+"/api/register", req, &registerResp); err != nil {
		return registerResp, err
	}
	if err := CheckTimestampDifference(registerResp.Timestamp); err != nil {
		return registerResp, err
	}
	return
}

func Login(req LoginReq) (loginResp LoginResp, err error) {
	req.Timestamp = time.Now().Unix()
	if code, err := requestCommon(Config.Url+"/api/login", req, &loginResp); err != nil {
		if code != 0 {
			return loginResp, GetErrorAndDefault(code, err)
		}
		return loginResp, err
	}
	if err := CheckTimestampDifference(loginResp.Timestamp); err != nil {
		return loginResp, err
	}
	return
}

func AutoLogin(req LoginReq) (loginResp LoginResp, err error) {
	req.Timestamp = time.Now().Unix()
	if code, err := requestCommon(Config.Url+"/api/autoLogin", req, &loginResp); err != nil {
		if code != 0 {
			return loginResp, GetErrorAndDefault(code, err)
		}
		return loginResp, err
	}
	if err := CheckTimestampDifference(loginResp.Timestamp); err != nil {
		return loginResp, err
	}
	return
}

func RenewTempAuth(req RenewTempAuthReq) (RenewTempAuthResp RenewTempAuthResp, err error) {
	req.Timestamp = time.Now().Unix()
	if _, err := requestCommon(Config.Url+"/api/renewTempAuth", req, &RenewTempAuthResp); err != nil {
		return RenewTempAuthResp, err
	}
	if err := CheckTimestampDifference(RenewTempAuthResp.Timestamp); err != nil {
		return RenewTempAuthResp, err
	}
	return
}

func Ping(req PingReq) (pingResp PingResp, err error) {
	req.Timestamp = time.Now().Unix()
	if code, err := requestCommon(Config.Url+"/api/ping", req, &pingResp); err != nil {
		if code != 0 {
			return pingResp, GetErrorAndDefault(code, err)
		}
		return pingResp, err
	}
	if err := CheckTimestampDifference(pingResp.Timestamp); err != nil {
		return pingResp, err
	}
	return
}

func RefreshInfo(req RefreshInfoReq) (pingResp RefreshInfoResp, err error) {
	req.Timestamp = time.Now().Unix()
	if code, err := requestCommon(Config.Url+"/api/refreshInfo", req, &pingResp); err != nil {
		if code != 0 {
			return pingResp, GetErrorAndDefault(code, err)
		}
		return pingResp, err
	}
	if err := CheckTimestampDifference(pingResp.Timestamp); err != nil {
		return pingResp, err
	}
	return
}

func CheckVersion(req CheckVersionReq) (checkVersionResp CheckVersionResp, err error) {
	req.Timestamp = time.Now().Unix()
	if code, err := requestCommon(Config.Url+"/api/checkVersion", req, &checkVersionResp); err != nil {
		if code != 0 {
			return checkVersionResp, GetErrorAndDefault(code, err)
		}
		return checkVersionResp, err
	}
	if err := CheckTimestampDifference(checkVersionResp.Timestamp); err != nil {
		return checkVersionResp, err
	}
	return
}

func CheckTimestampDifference(timestamp int64) error {
	currentTimestamp := time.Now().Unix()
	TimeDifference := math.Abs(float64(currentTimestamp) - float64(timestamp))
	if TimeDifference < float64(Config.TimeDifference) {
		return nil
	}
	return errors.New("request expire")
}

func requestCommon(url string, requestData interface{}, responseData interface{}) (code int, err error) {
	if Config.IsDebug {
		return request.SendPostJsonRequest(url, requestData, responseData)
	}
	return request.SendPostJsonAESKeyRequest(url, Config.AESKey, requestData, responseData)
}
