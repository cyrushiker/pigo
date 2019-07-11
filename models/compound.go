package models

import (
	"time"
)

const (
	// the rules of naming
	codeRP = `^[a-z]{2,}$`
	keyRP  = `^[a-z]+([a-zA-Z]+)$`
	grpRP  = `^[a-z]+(_?[a-z]+)+$`
)

type BaseModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createTime"`
	UpdatedAt time.Time  `json:"updateTime"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type Atom struct {
	BaseModel
	Key       string `json:"key" gorm:"type:varchar(100);unique_index;not null"`
	Name      string `json:"name"`
	NameEn    string `json:"name_en"`
	EnAbbr    string `json:"en_abbr"`
	Level     int    `json:"level" gorm:"default:1"`
	Keyword   string `json:"keyword"`
	Remark    string `json:"remark" gorm:"type:varchar(300)"`
	Options   string `json:"options" gorm:"type:varchar(500)"`
	Etype     string `json:"etype" gorm:"type:varchar(50)"`
	Eanalyzer string `json:"eanalyzer"`
}

type Molecule struct {
	BaseModel
	Key     string `json:"key" gorm:"type:varchar(100);unique_index;not null"`
	Name    string `json:"name"`
	NameEn  string `json:"name_en"`
	Keyword string `json:"keyword"`
	Remark  string `json:"remark" gorm:"type:varchar(300)"`
}

type Compound struct {
	BaseModel
	Code        string `json:"code" gorm:"type:varchar(100);unique_index;not null"`
	Name        string `json:"name"`
	Description string `json:"description" gorm:"type:varchar(500)"`
	Remark      string `json:"remark" gorm:"type:varchar(300)"`
	DbName      string `json:"db_name"`
}

type CompoundElem struct {
	BaseModel
	CompoundID uint `json:"compound_id"`
	MoleculeID uint `json:"molecule_id"`
	Position   int  `json:"position" gorm:"default:1"`
}

func (a *Atom) Save() error {
	// todo: 去重
	// todo: 事务
	// todo: 验证mk各个字段是否符合规则
	return db.Create(a).Error
}
