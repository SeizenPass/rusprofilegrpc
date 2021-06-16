package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type SearchServiceInterface interface {
	Search(uin string) ([]*SearchResponse, error)
}

type SearchResponse struct {
	UIN  string
	KPP  string
	Name string
	Bio  string
}

var (
	regexUin  = regexp.MustCompile("clip_inn.*>(.*)<")
	regexKpp  = regexp.MustCompile("clip_kpp.*>(.*)<")
	regexName = regexp.MustCompile("legalName.?>([^<]*)<")
	regexBio  = regexp.MustCompile(`<meta name="keywords" content=.*, (.*), ИНН `)
	regexPre  = regexp.MustCompile("Руководитель")

	regexSearchResult = regexp.MustCompile("search-result")
	regexEmptyResult  = regexp.MustCompile("emptyresult")
	regexNestedSearch = regexp.MustCompile("id/([0-9]*)")
)

type SearchServiceImpl struct{}

func (s *SearchServiceImpl) Search(uin string) ([]*SearchResponse, error) {
	res, err := http.Get(fmt.Sprintf("https://www.rusprofile.ru/search?query=%v", uin))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}
	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	input := string(body)
	ers := regexSearchResult.FindAllStringSubmatch(input, 1)
	if ers != nil && len(ers[0]) > 0 {
		emp := regexEmptyResult.FindAllStringSubmatch(input, 1)
		if emp != nil && len(emp[0]) > 0 {
			return nil, nil
		} else {
			ers := regexNestedSearch.FindAllStringSubmatch(input, -1)
			ch := make(chan *SearchResponse, len(ers))
			wg := &sync.WaitGroup{}
			var output []*SearchResponse
			for _, id := range ers {
				wg.Add(1)
				go SearchPage(id[1], wg, ch)
			}
			wg.Wait()
			close(ch)
			for data := range ch {
				output = append(output, data)
			}
			fmt.Println(output)
			return output, nil
		}
	} else {
		input := string(body)
		searchRes := &SearchResponse{}
		wg := sync.WaitGroup{}
		wg.Add(4)
		go func() {
			defer wg.Done()
			els := regexUin.FindAllStringSubmatch(input, 1)
			if els != nil && len(els[0]) > 0 {
				searchRes.UIN = els[0][1]
			}
		}()
		go func() {
			defer wg.Done()
			els := regexKpp.FindAllStringSubmatch(input, 1)
			if els != nil && len(els[0]) > 0 {
				searchRes.KPP = els[0][1]
			}
		}()
		go func() {
			defer wg.Done()
			els := regexName.FindAllStringSubmatch(input, 1)
			if els != nil && len(els[0]) > 0 {
				searchRes.Name = els[0][1]
			}
			searchRes.Name = strings.Replace(searchRes.Name, "&quot;", `"`, -1)
		}()
		go func() {
			defer wg.Done()
			ps := regexPre.FindAllStringSubmatch(input, 1)
			if ps != nil && len(ps[0]) > 0 {
				els := regexBio.FindAllStringSubmatch(input, 1)
				if els != nil && len(els[0]) > 0 {
					searchRes.Bio = els[0][1]
				}
			}
		}()
		wg.Wait()
		return []*SearchResponse{searchRes}, nil
	}
}

func SearchPage(ID string, group *sync.WaitGroup, ch chan *SearchResponse) {
	defer group.Done()
	res, err := http.Get(fmt.Sprintf("https://www.rusprofile.ru/id/%v", ID))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	input := string(body)
	searchRes := &SearchResponse{}
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		els := regexUin.FindAllStringSubmatch(input, 1)
		if els != nil && len(els[0]) > 0 {
			searchRes.UIN = els[0][1]
		}
	}()
	go func() {
		defer wg.Done()
		els := regexKpp.FindAllStringSubmatch(input, 1)
		if els != nil && len(els[0]) > 0 {
			searchRes.KPP = els[0][1]
		}
	}()
	go func() {
		defer wg.Done()
		els := regexName.FindAllStringSubmatch(input, 1)
		if els != nil && len(els[0]) > 0 {
			searchRes.Name = els[0][1]
		}
		searchRes.Name = strings.Replace(searchRes.Name, "&quot;", `"`, -1)
	}()
	go func() {
		defer wg.Done()
		ps := regexPre.FindAllStringSubmatch(input, 1)
		if ps != nil && len(ps[0]) > 0 {
			els := regexBio.FindAllStringSubmatch(input, 1)
			if els != nil && len(els[0]) > 0 {
				searchRes.Bio = els[0][1]
			}
		}
	}()
	wg.Wait()
	ch <- searchRes
}
