package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/YWJSonic/ServerUtility/playerinfo"
	"github.com/YWJSonic/ServerUtility/user"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gitlab.fbk168.com/gamedevjp/blackjack/server/conf"
	"gitlab.fbk168.com/gamedevjp/blackjack/server/controller/protoc"
)

var version string = "v1"

// AuthUserURL ...
const AuthUserURL string = "%s/%s/users/%s"

// NewOrderURL ...
const NewOrderURL string = "%s/%s/orders"

// GetOrderURL ...
const GetOrderURL string = "%s/%s/orders/%s"

// NewSubOrderURL ...
const NewSubOrderURL string = "%s/%s/sub-orders"

// // GetSubOrderURL ...
// const GetSubOrderURL string = "%s/%s/orders/%s"

// GetUser ...
func GetUser(userToken string) (*user.Info, *protoc.Error, error) {
	if conf.Dev {
		return &user.Info{
			UserServerInfo: &playerinfo.AccountInfo{},
			UserGameInfo: &playerinfo.Info{
				IDStr:  "devtest",
				MoneyU: 10000000,
			},
		}, nil, nil
	}

	tokens := strings.Split(userToken, " ")
	if len(tokens) < 2 {
		return nil, nil, errors.New("token error")
	}

	url := fmt.Sprintf(AuthUserURL, conf.ULG168APIHost, version, tokens[1])
	fmt.Println("API URL: ", url)
	res, err := authUserAPI(url)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := proto.Unmarshal(res, errorProto); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	userProto := &protoc.User{}
	if jserr := proto.Unmarshal(res, userProto); jserr != nil {
		return nil, nil, jserr
	}

	return &user.Info{
		UserServerInfo: &playerinfo.AccountInfo{
			Token: userToken,
		},
		UserGameInfo: &playerinfo.Info{
			IDStr:  userProto.GetUserId(),
			MoneyU: userProto.GetBalance(),
		},
	}, nil, nil
}

// NewOrder ...
func NewOrder(token, userIDStr string, betMoney int64) (*protoc.Order, *protoc.Error, error) {
	if conf.Dev {
		return &protoc.Order{
			UserId:  userIDStr,
			GameId:  conf.GameID,
			Bet:     uint64(betMoney),
			OrderId: "testOrder",
		}, nil, nil
	}

	orderProto := &protoc.Order{
		UserId: userIDStr,
		GameId: conf.GameID,
		Bet:    uint64(betMoney),
	}
	payload, err := proto.Marshal(orderProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := newOrderAPI(fmt.Sprintf(NewOrderURL, conf.ULG168APIHost, version), token, payload)
	fmt.Println(string(res))
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := proto.Unmarshal(res, errorProto); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := proto.Unmarshal(res, orderProto); jserr != nil {
		return nil, nil, jserr
	}
	return orderProto, nil, nil
}

// NewSubOrder ...
func NewSubOrder(token string, orderProto *protoc.Order, betMoney int64) (*protoc.SubOrder, *protoc.Error, error) {
	if conf.Dev {
		return &protoc.SubOrder{
			State:      protoc.SubOrder_Pending,
			Bet:        uint64(betMoney),
			OrderId:    "testOrder",
			SubOrderId: "testSubOrder",
		}, nil, nil
	}

	orderSubProto := &protoc.SubOrder{
		Bet:     uint64(betMoney),
		OrderId: orderProto.OrderId,
	}
	payload, err := proto.Marshal(orderSubProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := newSubOrderAPI(fmt.Sprintf(NewSubOrderURL, conf.ULG168APIHost, version), token, payload)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := proto.Unmarshal(res, errorProto); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := proto.Unmarshal(res, orderSubProto); jserr != nil {
		return nil, nil, jserr
	}
	return orderSubProto, nil, nil
}

// UpdateOrder ...
func UpdateOrder(token string, orderProto *protoc.Order) (*protoc.Order, *protoc.Error, error) {
	if conf.Dev {
		return orderProto, nil, nil
	}

	payload, err := proto.Marshal(orderProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := updateOrderAPI(fmt.Sprintf(GetOrderURL, conf.ULG168APIHost, version, orderProto.GetOrderId()), token, payload)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := proto.Unmarshal(res, errorProto); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := proto.Unmarshal(res, orderProto); jserr != nil {
		return nil, nil, jserr
	}
	return orderProto, nil, nil
}

// EndOrder ...
func EndOrder(token string, orderProto *protoc.Order) (*protoc.Order, *protoc.Error, error) {
	orderProto.CompletedAt = ptypes.TimestampNow()
	if conf.Dev {
		return orderProto, nil, nil
	}

	payload, err := proto.Marshal(orderProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := updateOrderAPI(fmt.Sprintf(GetOrderURL, conf.ULG168APIHost, version, orderProto.GetOrderId()), token, payload)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := proto.Unmarshal(res, errorProto); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := proto.Unmarshal(res, orderProto); jserr != nil {
		return nil, nil, jserr
	}
	return orderProto, nil, nil
}

// ---------------------------------

// authUserAPI GET transation
func authUserAPI(url string) ([]byte, error) {
	fmt.Println(url)
	res, err := httpGET(url, map[string][]string{})
	if len(res) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(url + " return empty data.")
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

// newOrderAPI POST transation 下注
func newOrderAPI(url, token string, payload []byte) ([]byte, error) {
	header := map[string][]string{
		"Authorization": []string{token},
		"Content-Type":  []string{"application/protobuf"},
	}

	res, err := httpPOST(url, payload, header)
	if len(res) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(url + " return empty data.")
	}
	if err != nil {
		return res, err
	}

	return res, nil
}

// updateOrderAPI GET transation 	投降、
func updateOrderAPI(url, token string, payload []byte) ([]byte, error) {
	header := map[string][]string{
		"Authorization": []string{token},
		"Content-Type":  []string{"application/protobuf"},
	}

	res, err := httpPUT(url, payload, header)
	if len(res) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(url + " return empty data.")
	}

	if err != nil {
		return res, err
	}
	return res, nil
}

// newSubOrderAPI POST sub transation 分牌、加倍、保險
func newSubOrderAPI(url, token string, payload []byte) ([]byte, error) {
	header := map[string][]string{
		"Authorization": []string{token},
		"Content-Type":  []string{"application/protobuf"},
	}

	res, err := httpPOST(url, payload, header)
	fmt.Println("API URL: ", url)
	if len(res) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(url + " return empty data.")
	}

	if err != nil {
		return res, err
	}
	return res, nil
}

// ---------------------------------

func httpGET(url string, header map[string][]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header = header
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode == 200 || res.StatusCode == 201) {
		return body, errors.New(res.Status)
	}

	return body, nil
}
func httpPOST(url string, value []byte, header map[string][]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	if _, ok := header["Content-Type"]; !ok {
		header["Content-Type"] = []string{"application/json"}
	}
	req.Header = header
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode == 200 || res.StatusCode == 201) {
		return body, errors.New(res.Status)
	}

	return body, nil
}
func httpPUT(url string, value []byte, header map[string][]string) ([]byte, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(value))
	if _, ok := header["Content-Type"]; !ok {
		header["Content-Type"] = []string{"application/json"}
	}
	req.Header = header
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 || res.StatusCode == 201 {
		return body, errors.New(res.Status)
	}

	return body, nil
}
