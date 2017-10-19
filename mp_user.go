package wechat

import (
	"fmt"

	"github.com/esap/wechat/util"
)

// MP 公众号接口
const (
	MP_GetUserList = WXAPI + "user/get?access_token=%s&next_openid=%s"
	MP_BatchGet    = WXAPI + "user/info/batchget?access_token="
)

type (
	MpUserInfoList struct {
		WxErr
		MpUserInfoList []MpUserInfo `json:"user_info_list"`
	}

	//MpUserInfo
	MpUserInfo struct {
		Subscribe     int
		OpenId        string
		NickName      string
		Sex           int
		Language      string
		City          string
		Province      string
		Country       string
		HeadImgUrl    string
		SubscribeTime int `json:"subscribe_time"`
		UnionId       string
		Remark        string
		GroupId       int
		TagIdList     []int `json:"tagid_list"`
	}

	// MpUser
	MpUser struct {
		WxErr
		Total int
		Count int
		Data  struct {
			OpenId []string
		}
		NextOpenId string
	}

	MpUserListReq struct {
		UserList interface{} `json:"user_list"`
	}
)

func (s *Server) BatchGetAll() (ui []MpUserInfo, err error) {
	var ul []string
	ul, err = s.GetAllMpUserList()
	if err != nil {
		return
	}
	leng := len(ul)
	if leng <= 100 {
		return s.BatchGet(ul)
	}
	for i := 0; i < leng/100+1; i++ {
		end := (i + 1) * 100
		if end > leng {
			end = leng
		}

		ui2, err2 := s.BatchGet(ul[i*100 : end])
		if err != nil {
			err = err2
			return
		}
		ui = append(ui, ui2...)
	}
	return
}

func (s *Server) BatchGet(ul []string) (ui []MpUserInfo, err error) {
	m := make([]map[string]interface{}, len(ul))

	for k, v := range ul {
		m[k] = make(map[string]interface{})
		m[k]["openid"] = v
	}
	ml := new(MpUserInfoList)
	err = util.PostJsonPtr(MP_BatchGet+s.GetAccessToken(), MpUserListReq{m}, ml)
	return ml.MpUserInfoList, ml.Error()
}

// GetAllMpUserList
func (s *Server) GetAllMpUserList() (ul []string, err error) {
	ul = make([]string, 0)
	mul, err := s.GetMpUserList()
	if err != nil {
		return
	}
	if mul.Error() == nil {
		ul = append(ul, mul.Data.OpenId...)
	}
	for mul.Count == 10000 {
		mul, err = s.GetMpUserList(mul.NextOpenId)
		if err != nil {
			return
		}
		if mul.Error() == nil {
			ul = append(ul, mul.Data.OpenId...)
		}
	}
	return
}

// GetMpUserList
func (s *Server) GetMpUserList(openid ...string) (ul *MpUser, err error) {
	if len(openid) == 0 {
		openid = append(openid, "")
	}
	mpuser := new(MpUser)
	url := fmt.Sprintf(MP_GetUserList, s.GetAccessToken(), openid[0])
	if err = util.GetJson(url, &mpuser); err != nil {
		return
	}
	return mpuser, mpuser.Error()
}
