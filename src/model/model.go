package model

type Anime struct {
	Id           string `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Chapter      string `db:"chapter" json:"chapter"`
	Total        int    `db:"total" json:"total"`
	Update       string `db:"update" json:"update"`
	Index        int    `db:"index" json:"index"`
	Created      string `db:"created" json:"created"`
	Year         string `db:"year" json:"year"`
	Area         string `db:"area" json:"area"`
	Picture      string `db:"picture" json:"picture"`
	Introduction string `db:"introduction" json:"introduction"`
	Form         string `db:"form" json:"form"`
	Flag         string `db:"flag" json:"flag"`
}
type Chapter struct {
	Id      string `db:"id" json:"id"`
	Pid     string `db:"pid" json:"pid"`
	Name    string `db:"name" json:"name"`
	Path    string `db:"path" json:"path"`
	Source  string `db:"source" json:"source"`
	JX      string `db:"j_x" json:"j_x"`
	Created string `db:"created" json:"created"`
}
type Director struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Pid        string `db:"pid" json:"pid"`
	CreateTime string `db:"create_time" json:"create_time"`
}
type Star struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Pid        string `db:"pid" json:"pid"`
	CreateTime string `db:"create_time" json:"create_time"`
}
