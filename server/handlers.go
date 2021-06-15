package main

import (
	proto "github.com/SeizenPass/rusprofilegrpc/proto"
)

func (a *app) Search(req *proto.SearchRequest, stream proto.SearchService_SearchServer) error  {
	data, err := a.service.Search(req.GetUin())
	if err != nil {
		return err
	}
	if data != nil {
		for _, resp := range data {
			res := &proto.SearchResponse{
				Uin:  resp.UIN,
				Kpp:  resp.KPP,
				Name: resp.Name,
				Bio:  resp.Bio,
			}
			if err := stream.Send(res); err != nil {
				a.errorLog.Printf("error while sending responses: %v", err.Error())
			}
		}
	}
	return nil
}