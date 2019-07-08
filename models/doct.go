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

func (d *Doct) esTypeName() string {
	return "docut"
}

func (d *Doct) Add() error {
	return nil
}

type DKey struct {
}

type KGroup struct {
}

type KeyDefine struct {
	Id     string
	Code   string
	Groups []struct {
		Id       string
		Key      string
		Name     string
		Rank     int
		MetaKeys []struct {
			Id       string
			Key      string
			Name     string
			Keyword  string
			Etype    string
			Optional string
		}
	}
}

type DocVerify struct {
	kd           *KeyDefine
	doc          map[string]interface{}
	raiseOnError bool
	verifies     []interface{}
}

func (dv *DocVerify) Verify() {
	for _, g := range dv.kd.Groups {
		if d, ok := dv.doc[g.Key]; ok {
			logger.Println(d)
			if d, ok := d.([]interface{}); ok {
				logger.Println(d)
			}
		}
	}
}
