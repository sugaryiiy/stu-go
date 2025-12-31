package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
	"stu-go/common"
	"time"
)

type service struct {
	repo repository
}

func (s *service) getSign(util *SignUtil) error {

	timeStamp := time.Now().Unix()
	if util.Status != 1 {
		timeStamp = time.Now().Unix() * 100
	}
	str := fmt.Sprintf("%sappid=%s&time_stamp=%d%s", util.AppSecret, util.AppId, timeStamp, util.AppSecret)
	hash := md5.Sum([]byte(str))
	pF := "&"
	if !strings.Contains(util.Url, "?") {
		pF = "?"
	} else if strings.HasSuffix(util.Url, "?") {
		pF = ""
	} else {
		pF = "&"
	}
	result := common.Result{
		Code: 0,
		Msg:  "",
		Data: nil,
	}
	result.Data = nil
	util.Url = util.Url + pF
	util.Result = fmt.Sprintf("%sappid=%s&time_stamp=%d&sign=%x", util.Url, util.AppId, timeStamp, hash)
	result.Data = util
	return nil
}
