package models

import (
	"time"
)

type Doct struct {
	Id         string    `json:"id" dt:"keyword"`
	Name       string    `json:"name,omitempty" dt:"keyword"`
	Key        string    `json:"key" dt:"keyword"`
	Desc       string    `json:"desc,omitempty" dt:"text"`
	CreateDate time.Time `json:"createDate,omitempty" dt:"date"`
	Position   int
}

func (d *Doct) Add() error {
	return nil
}

type DKey struct {
}

type KGroup struct {
}
