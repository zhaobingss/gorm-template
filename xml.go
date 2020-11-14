package template

import (
	"errors"
	"github.com/beevik/etree"
	"io/ioutil"
)

type SqlMap struct {
	Namespace string
	Sqls      []*Sql
}

type Sql struct {
	Id      string
	Content string
}

func parseSqlmap(path string) (*SqlMap, error) {

	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}


	doc := etree.NewDocument()
	err = doc.ReadFromBytes(bts)
	if err != nil {
		return nil, err
	}

	sqlmapEle := doc.SelectElement("sqlmap")
	if sqlmapEle == nil {
		return nil, errors.New("can not find sqlmap tag in file:" + path)
	}

	sqlmap := &SqlMap{}
	sqlmap.Namespace = sqlmapEle.SelectAttrValue("namespace", "")
	if sqlmap.Namespace == "" {
		return nil, errors.New("sqlmap attr namespace is empty in file:" + path)
	}

	sqls := make([]*Sql, 0)
	sqlmap.Sqls = sqls

	sqlEles := sqlmapEle.SelectElements("sql")
	if len(sqlEles) < 1 {
		return sqlmap, nil
	}

	for _, sqlEle := range sqlEles {
		s, err := parseSql(sqlEle)
		if err != nil {
			return nil, errors.New(err.Error() + " in file:" + path)
		}
		if s != nil {
			sqlmap.Sqls = append(sqlmap.Sqls, s)
		}
	}


	return sqlmap, nil
}

func parseSql(sqlEle *etree.Element) (*Sql, error) {

	id := sqlEle.SelectAttrValue("id", "")
	if id == "" {
		return nil, errors.New("id is empty")
	}

	sql := &Sql{}
	sql.Id = id
	sql.Content = sqlEle.Text()

	return sql, nil
}