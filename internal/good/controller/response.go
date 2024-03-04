package controller

import "goods-manager/internal/domain/entity"

type ListResponse struct {
	Meta  *Meta          `json:"meta"`
	Goods []*entity.Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

func MetaFromGoods(goods []*entity.Good, limit, offset int) *Meta {
	countRemoved := 0
	for _, good := range goods {
		if good.Removed {
			countRemoved++
		}
	}

	return &Meta{
		Total:   len(goods),
		Removed: countRemoved,
		Limit:   limit,
		Offset:  offset,
	}
}

type PrioritizeResponse struct {
	Priorities []UpratedPriority `json:"priorities"`
}

type UpratedPriority struct {
	Id       int `json:"id"`
	Priority int `json:"priority"`
}

func PrioritizeResponseFromMap(prioritize map[int]int) *PrioritizeResponse {
	priorities := make([]UpratedPriority, 0, len(prioritize))

	for id, priority := range prioritize {
		priorities = append(priorities, UpratedPriority{Id: id, Priority: priority})
	}

	return &PrioritizeResponse{Priorities: priorities}
}
