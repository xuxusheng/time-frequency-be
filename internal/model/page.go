package model

type Page struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

func (p *Page) WithTotal(total int) *Page {
	np := *p
	np.Total = total
	return &np
}

func (p *Page) Offset() int {
	offset := (p.Pn - 1) * p.Ps
	if offset < 0 {
		offset = 0
	}
	return offset
}

func (p *Page) Limit() int {
	if p.Ps < 0 {
		p.Ps = 0
	}
	return p.Ps
}
