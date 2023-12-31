package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/murtll/mcserver-rcon/pkg/entities"
	"github.com/murtll/mcserver-rcon/pkg/util"
)

type ItemRepository struct {
	ApiUrl *url.URL
	ApiKey string

	c *http.Client
}

func NewItemRepository(apiUrl string, apiKey string) (*ItemRepository, error) {
	parsed, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	ir := &ItemRepository{
		ApiUrl: parsed,
		ApiKey: apiKey,
		c: &http.Client{
			Jar: jar,
		},
	}

	err = ir.Authorize()
	if err != nil {
		return nil, err
	}

	return ir, nil
}

func (ir *ItemRepository) GetItem(id int) (*entities.Item, error) {
	requestUrl := ir.ApiUrl.JoinPath("admin").JoinPath("item")
	util.SetQueryParam(requestUrl, "id", strconv.Itoa(id))

	triedAuth := false

makerequest:
	res, err := ir.c.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		if res.StatusCode == 401 && !triedAuth {
			ir.Authorize()
			triedAuth = !triedAuth
			goto makerequest
		}
		return nil, fmt.Errorf("can't get item, status '%s' is not acceptable", res.Status)
	}

	item := &entities.Item{}
	err = json.NewDecoder(res.Body).Decode(item)
	if err != nil {
		log.Println("error unmarshalling item")
		return nil, err
	}

	return item, nil
}

func (ir *ItemRepository) Authorize() error {
	requestUrl := ir.ApiUrl.JoinPath("admin").JoinPath("login")
	req, err := http.NewRequest("POST", requestUrl.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", ir.ApiKey)
	res, err := ir.c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("can't authorize to api, status '%s' is not acceptable", res.Status)
	}

	return nil
}
