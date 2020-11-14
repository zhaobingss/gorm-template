package gorm_template

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/CloudyKit/jet/v6"
	"gorm.io/gorm"
)

type gormEngine struct {
	GDB   *gorm.DB
	set   *jet.Set
	vars  jet.VarMap
	cache map[string]*jet.Template
}

func NewGormEngine(dialector gorm.Dialector, config *gorm.Config, sqlDir string) (*gormEngine, error) {
	eg := &gormEngine{}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	eg.GDB = db
	eg.set = jet.NewSet(jet.NewInMemLoader())
	eg.vars = make(jet.VarMap)
	eg.cache = map[string]*jet.Template{}
	err = eg.loadTemplate(sqlDir)
	return eg, err
}

func (e *gormEngine) clone() *gormEngine {
	n := &gormEngine{}
	n.GDB = e.GDB
	n.cache = e.cache
	n.vars = e.vars
	n.set = e.set
	return n
}

func (e *gormEngine) cloneWithDB(db *gorm.DB) *gormEngine {
	n := e.clone()
	n.GDB = db
	return n
}

func (e *gormEngine) DB() (*sql.DB, error) {
	return e.GDB.DB()
}

func (e *gormEngine) QueryTpl(name string, param interface{}, dest interface{}) error {

	tpl, ok := e.cache[name]
	if !ok {
		return errors.New("can't find sql name:" + name)
	}

	var bts bytes.Buffer

	err := tpl.Execute(&bts, e.vars, param)
	if err != nil {
		return err
	}

	return e.GDB.Raw(bts.String()).Find(dest).Error

}

func (e *gormEngine) ExecTpl(name string, param interface{}) (int64, error) {

	tpl, ok := e.cache[name]
	if !ok {
		return 0, errors.New("can't find sql name:" + name)
	}

	var bts bytes.Buffer

	err := tpl.Execute(&bts, e.vars, param)
	if err != nil {
		return 0, err
	}

	db := e.GDB.Exec(bts.String())

	return db.RowsAffected, db.Error
}

func (e *gormEngine) Transcation(f func(e *gormEngine) error, opts ...*sql.TxOptions) error {
	return e.GDB.Transaction(func(tx *gorm.DB) error {
		return f(e.cloneWithDB(tx))
	}, opts...)
}

func (e *gormEngine) loadTemplate(sqlDir string) error {
	files, err := GetFiles(sqlDir)
	if err != nil {
		return err
	}

	if len(files) < 1 {
		return nil
	}

	for _, file := range files {
		err = e.loadTemplateToCache(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *gormEngine) loadTemplateToCache(filePath string) error {

	sqlmap, err := parseSqlmap(filePath)
	if err != nil {
		return err
	}

	ns := sqlmap.Namespace
	sqls := sqlmap.Sqls
	if len(sqls) < 1 {
		return nil
	}

	for _, s := range sqls {
		err = e.loadTemplateSqlToCache(ns, s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *gormEngine) loadTemplateSqlToCache(ns string, sql *Sql) error {
	name := ns + "." + sql.Id
	tpl, err := e.set.Parse(name, sql.Content)
	if err != nil {
		return err
	}

	e.cache[name] = tpl

	return nil
}
