package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/assert"
	"github.com/zalora/benandjerry/httputil"
	"github.com/zalora/benandjerry/model"
)

type fakeIceCreamStore struct {
	Store          []*model.IceCream
	ListError      error
	GetByNameError error
	CreateError    error
	UpdateError    error
	DeleleError    error
}

func (f fakeIceCreamStore) List() ([]*model.IceCream, error) {
	if f.ListError != nil {
		return nil, f.ListError
	}
	return f.Store, nil
}

func (f fakeIceCreamStore) Get(name string) (*model.IceCream, error) {
	fmt.Println("=====s", name)

	if f.GetByNameError != nil {
		return nil, f.GetByNameError
	}

	for _, iceCreamInfo := range f.Store {
		if iceCreamInfo != nil && name == iceCreamInfo.Name {
			fmt.Println("s", name)
			return iceCreamInfo, nil
		}
	}
	return nil, &model.NotFoundError{}
}

func (f fakeIceCreamStore) Create(iceCream *model.IceCream) error {
	if f.CreateError != nil {
		return f.CreateError
	}

	f.Store = append(f.Store, iceCream)
	return nil
}

func (f fakeIceCreamStore) Update(newIceCream *model.IceCream) error {
	if f.UpdateError != nil {
		return f.UpdateError
	}

	for i, iceCreamInfo := range f.Store {
		if iceCreamInfo.Name == newIceCream.Name {
			f.Store[i] = newIceCream
			return nil
		}
	}
	return &model.NotFoundError{}
}

func (f fakeIceCreamStore) Delete(name string) error {
	if f.DeleleError != nil {
		return f.DeleleError
	}

	for i, iceCreamInfo := range f.Store {
		if name == iceCreamInfo.Name {
			f.Store[i] = nil
		}
	}
	return &model.NotFoundError{}
}

var sampleList = []*model.IceCream{
	{
		Name: "Sample1",
	},
	{
		Name: "Sample2",
	},
}
var getIceCreamsTestCases = []struct {
	description               string
	params                    url.Values
	fakeIceCreamStore         *fakeIceCreamStore
	expectedNoOfIceCreamsInfo int
	expectedIceCreamsInfo     []*model.IceCream
	expectedHandlerErr        *httputil.HandlerError
	expectedStatusCode        int
}{
	{
		"lists all iceCreams ",
		url.Values{},
		&fakeIceCreamStore{Store: sampleList},
		2,
		sampleList,
		nil,
		http.StatusOK,
	},
}

func fakeHttpRequest(method, path string, data []byte) *http.Request {
	return httptest.NewRequest(method, path, bytes.NewReader(data))
}

func Test_GetIceCreams_Returns_IceCreams_For_Given_Criteria(t *testing.T) {

	assert := assert.New(t)
	for _, testcase := range getIceCreamsTestCases {
		rr := httptest.NewRecorder()
		iceCreamHandler := &IceCreamHandler{IceCreamStore: testcase.fakeIceCreamStore}
		handlerErr := iceCreamHandler.List(rr, fakeHttpRequest("GET", "/iceCreams"+"?"+testcase.params.Encode(), nil), nil)
		t.Logf("%s", testcase.description)

		if testcase.expectedHandlerErr != nil {
			if assert.NotNil(handlerErr) {
				assert.Equal(testcase.expectedHandlerErr, handlerErr)
			}
		} else {
			rawResponse, _ := ioutil.ReadAll(rr.Body)
			response := struct {
				IceCreams []*model.IceCream `json:"IceCreams"`
			}{}
			assert.Equal(testcase.expectedStatusCode, rr.Code)
			json.Unmarshal(rawResponse, &response)
			assert.Len(response.IceCreams, testcase.expectedNoOfIceCreamsInfo)
			assert.Equal(testcase.expectedIceCreamsInfo, response.IceCreams)
		}

	}
}

var getIceCreamTestCases = []struct {
	description          string
	params               url.Values
	fakeIceCreamStore    *fakeIceCreamStore
	expectedIceCreamInfo *model.IceCream
	expectedHandlerErr   *httputil.HandlerError
	expectedStatusCode   int
}{
	{
		"Getting a specific info",
		url.Values{"name": []string{"Sample1"}},
		&fakeIceCreamStore{Store: sampleList},
		sampleList[0],
		nil,
		http.StatusOK,
	},
	{
		"Getting a specific info",
		url.Values{"name": []string{"random"}},
		&fakeIceCreamStore{Store: sampleList},
		nil,
		httputil.NewNotFoundError(fmt.Sprintf("'%s' not found", "random"), errors.New("not_found")),
		http.StatusNotFound,
	},
}

func Test_GetIceCream_Returns_IceCream_For_Given_Criteria(t *testing.T) {

	assert := assert.New(t)
	for _, testcase := range getIceCreamTestCases {
		rr := httptest.NewRecorder()
		iceCreamHandler := &IceCreamHandler{IceCreamStore: testcase.fakeIceCreamStore}
		path := fmt.Sprintf("/icecreams/show/%s", testcase.params)

		params := httprouter.Params{httprouter.Param{"name", testcase.params.Get("name")}}

		handlerErr := iceCreamHandler.Get(rr, fakeHttpRequest("GET", path, nil), params)
		t.Logf("%s", testcase.description)

		if testcase.expectedHandlerErr != nil {
			if assert.NotNil(handlerErr) {
				assert.Equal(testcase.expectedHandlerErr, handlerErr)
			}
		} else {
			rawResponse, _ := ioutil.ReadAll(rr.Body)

			var iceCream = &model.IceCream{}
			assert.Equal(testcase.expectedStatusCode, rr.Code)
			json.Unmarshal(rawResponse, iceCream)
			assert.Equal(testcase.expectedIceCreamInfo, iceCream)
		}

	}
}

var createNewInfoTestcases = []struct {
	description          string
	params               url.Values
	fakeIceCreamStore    *fakeIceCreamStore
	expectedIceCreamInfo *model.IceCream
	expectedHandlerErr   *httputil.HandlerError
	expectedStatusCode   int
}{
	{
		"create a new one",
		url.Values{"name": []string{"new"}},
		&fakeIceCreamStore{Store: sampleList},
		&model.IceCream{
			Name: "new",
		},
		nil,
		http.StatusOK,
	},
}

func newSampleIceCreamJSON(name string) []byte {
	return []byte(`{
  "name": "` + name + `",
  "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
  "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
  "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
  "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
  "sourcing_values": [
    "Non-GMO",
    "Cage-Free Eggs",
    "Fairtrade",
    "Responsibly Sourced Packaging",
    "Caring Dairy"
  ],
  "ingredients": [
    "cream",
    "skim milk",
    "liquid sugar",
    "water",
    "sugar",
    "coconut oil",
    "egg yolks",
    "butter",
    "vanilla extract",
    "almonds",
    "cocoa (processed with alkali)",
    "milk",
    "soy lecithin",
    "cocoa",
    "natural flavor",
    "salt",
    "vegetable oil",
    "guar gum",
    "carrageenan"
  ],
  "allergy_info": "may contain wheat, peanuts and other tree nuts",
  "dietary_certifications": "Kosher",
  "productId": "646"
}`)
}

func Test_createNewInfoTestcases_creates_a_new_IceCream_Info(t *testing.T) {

	assert := assert.New(t)
	for _, testcase := range createNewInfoTestcases {
		rr := httptest.NewRecorder()
		iceCreamHandler := &IceCreamHandler{IceCreamStore: testcase.fakeIceCreamStore}
		path := fmt.Sprintf("/icecreams/create")
		handlerErr := iceCreamHandler.Create(rr, fakeHttpRequest("POST", path, newSampleIceCreamJSON(testcase.params.Get("name"))), nil)
		t.Logf("%s", testcase.description)

		if testcase.expectedHandlerErr != nil {
			if assert.NotNil(handlerErr) {
				assert.Equal(testcase.expectedHandlerErr, handlerErr)
			}
		} else {
			rawResponse, _ := ioutil.ReadAll(rr.Body)

			response := struct {
				Name string `json:"name"`
			}{}
			assert.Equal(testcase.expectedStatusCode, rr.Code)
			json.Unmarshal(rawResponse, &response)
			assert.Equal(testcase.expectedIceCreamInfo.Name, response.Name)
		}

	}
}

var deleteTestcases = []struct {
	description        string
	params             url.Values
	fakeIceCreamStore  *fakeIceCreamStore
	expectedHandlerErr *httputil.HandlerError
	expectedStatusCode int
}{
	{
		"delete the existing one",
		url.Values{"name": []string{"new"}},
		&fakeIceCreamStore{Store: sampleList},
		nil,
		http.StatusOK,
	},
	{
		"delete the existing one",
		url.Values{"name": []string{"random"}},
		&fakeIceCreamStore{Store: sampleList},
		httputil.NewNotFoundError(fmt.Sprintf("'%s' not found", "random"), errors.New("not_found")),
		http.StatusOK,
	},
}

func Test_delete(t *testing.T) {

	assert := assert.New(t)
	for _, testcase := range deleteTestcases {
		rr := httptest.NewRecorder()
		iceCreamHandler := &IceCreamHandler{IceCreamStore: testcase.fakeIceCreamStore}
		path := fmt.Sprintf("/icecreams/delete")
		params := httprouter.Params{httprouter.Param{"name", testcase.params.Get("name")}}

		handlerErr := iceCreamHandler.Delete(rr, fakeHttpRequest("POST", path, newSampleIceCreamJSON(testcase.params.Get("name"))), params)
		t.Logf("%s", testcase.description)

		if testcase.expectedHandlerErr != nil {
			if assert.NotNil(handlerErr) {
				assert.Equal(testcase.expectedHandlerErr, handlerErr)
			}
		} else {
			assert.Equal(testcase.expectedStatusCode, rr.Code)
		}

	}
}
