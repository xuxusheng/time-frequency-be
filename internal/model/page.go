package model

import "github.com/xuxusheng/time-frequency-be/global"

type Page struct {
	pn    int
	ps    int
	total int
}

func New(pn, ps int) *Page {
	if pn < 0 {
		pn = 1
	}
	if ps < 0 || ps > global.Setting.App.MaxPs {
		ps = global.Setting.App.DefaultPs
	}
	return &Page{
		pn: pn,
		ps: ps,
	}
}

func (p *Page) Pn() int {
	return p.pn
}

func (p *Page) Ps() int {
	return p.ps
}

func (p *Page) Total() int {
	return p.total
}

func (p *Page) Offset() int {
	return (p.Pn() - 1) * p.Ps()
}

func (p *Page) Limit() int {
	return p.Ps()
}

func (p *Page) WithTotal(total int) *Page {
	n := p.clone()
	n.total = total
	return n
}

func (p *Page) clone() *Page {
	n := *p
	return &n
}
