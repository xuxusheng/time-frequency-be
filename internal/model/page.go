package model

type Page struct {
	Pn    int   `json:"pn"`
	Ps    int   `json:"ps"`
	Total int64 `json:"total"`
}

func (p *Page) WithTotal(total int64) *Page {
	np := *p
	np.Total = total
	return &np
}

func (p *Page) GetOffset() int {
	offset := (p.Pn - 1) * p.Ps
	if offset < 0 {
		offset = 0
	}
	return offset
}
