package ctrl

import (
	"fmt"
	"io"
	"model"
	"net/http"
	"os"
	"time"
	"tools"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func SaveImg(url, imgPath, path string, str []string) {
	///upload/vod/20190712-1/37028a8a314e23ed79ef7e4c31dd14b4.jpg

	fmt.Println(str)
	fmt.Println(len(str))

	url = "http://pilipali.cc" + url
	bol := Exists(path)

	if !bol {
		err1 := os.Mkdir(path, os.ModePerm) //创建文件夹
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		f, err := os.Create(imgPath)

		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		buf := make([]byte, 4096)
		for {
			n, err1 := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			if err1 != nil && err1 != io.EOF {
				err = err1
				return
			}

			f.Write(buf[:n])
		}
	} else {
		f, err := os.Create(imgPath)
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		buf := make([]byte, 4096)
		for {
			n, err1 := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			if err1 != nil && err1 != io.EOF {
				err = err1
				return
			}

			f.Write(buf[:n])
		}
	}
	time.Sleep(time.Second)
}
func Save2Mysql(AnimeData model.AnimeData, Picture string) {
	dbConn := tools.GetDefDb()

	AnimeDatas := []model.Anime{}
	_, err := dbConn.DbMap.Select(&AnimeDatas, "select * from anime where name=?", AnimeData.Title)
	tools.CheckErr(err)

	if len(AnimeDatas) == 0 {
		AnimeDatasRs, err := dbConn.DbMap.Exec("insert into `anime` (`name`,`em_num`,`year`,`area`,`picture`,`introduction`,`form`"+
			") values(?,?,?,?,?,?,?)", AnimeData.Title, AnimeData.EmNum, AnimeData.Year, AnimeData.Area, Picture, AnimeData.Introduction,
			"pilipili")
		tools.CheckErr(err)
		lastId, err := AnimeDatasRs.LastInsertId()
		tools.CheckErr(err)

		for _, item := range AnimeData.Star {
			stars := []model.Star{}
			_, err := dbConn.DbMap.Select(&stars, "select * from star where name=?", item)
			tools.CheckErr(err)

			if len(stars) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `star` (`name`,`anime_id`"+
					") values(?,?)", item, lastId)
				tools.CheckErr(err)
			}
		}

		for _, item := range AnimeData.Director {
			directors := []model.Director{}
			_, err = dbConn.DbMap.Select(&directors, "select * from director where name=?", item)
			tools.CheckErr(err)
			if len(directors) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `director` (`name`,`anime_id`"+
					") values(?,?)", item, lastId)
				tools.CheckErr(err)
			}
		}

		for _, item := range AnimeData.Drama {
			dramas := []model.Drama{}
			_, err = dbConn.DbMap.Select(&dramas, "select * from drama where name=? and play_url=? and source=? ",
				item.Name, item.Url, item.From)
			tools.CheckErr(err)
			if len(dramas) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `drama` (`name`,`anime_id`,`play_url`,`source`"+
					") values(?,?,?,?)", item.Name, lastId, item.Url, item.From)
				tools.CheckErr(err)
			}
		}
	} else {
		for _, item := range AnimeData.Drama {
			dramas := []model.Drama{}
			_, err = dbConn.DbMap.Select(&dramas, "select * from drama where name=? and play_url=? and source=? and anime_id=?",
				item.Name, item.Url, item.From, AnimeDatas[0].Id)
			tools.CheckErr(err)
			if len(dramas) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `drama` (`name`,`anime_id`,`play_url`,`source`"+
					") values(?,?,?,?)", item.Name, AnimeDatas[0].Id, item.Url, item.From)
				tools.CheckErr(err)
			}
		}
	}
}
