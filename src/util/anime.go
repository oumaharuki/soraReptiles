package util

import (
	"encoding/json"
	"fmt"
	"model"
	"regexp"
	"strconv"
	"sync"
	"time"
	"tools"
)

var (
	wgGetPage sync.WaitGroup
	wgSave    sync.WaitGroup
	dataChan  chan model.AnimeData
	doChan    = make(chan int, 5)
)

func Anime(start, end string) {
	dataChan = make(chan model.AnimeData, 4000)
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)
	for i := startInt; i <= endInt; i++ {
		fmt.Printf("第%d页开始\n", i)
		url := "http://pilipali.cc/vod/show/id/4/page/" + strconv.Itoa(i) + ".html"

		wgGetPage.Add(1)
		go func(pageUrl string) {
			GetPageData(pageUrl)
			wgGetPage.Done()
		}(url)

	}
	go func() {
		wgGetPage.Wait()
		close(dataChan)
		fmt.Println("chan close")
	}()

	for item := range dataChan {

		wgSave.Add(1)
		go func(AnimeData model.AnimeData) {
			doChan <- 1
			str := extractHandle(AnimeData.Picture, `/([0-9a-z]+\.[a-z]+)`, 1)
			nowTime := int(time.Now().Unix())
			timestr := strconv.Itoa(nowTime)
			path := "./public/upload/anime"
			imgPath := path + "/" + timestr + ".jpg"
			img := "/upload/anime/" + timestr + ".jpg"
			if len(str) > 0 {
				imgPath = path + "/" + str[0]
				img = "/upload/anime/" + str[0]
			}

			SaveImg(AnimeData.Picture, imgPath, path, str)
			Save2Mysql(AnimeData, img)
			<-doChan
			wgSave.Done()
		}(item)

	}
	wgSave.Wait()
}
func GetPageData(url string) {
	//得到每一页动漫数据
	rs, err := HttpGet(url)
	if err != nil {
		return
	}

	reg := regexp.MustCompile(`<a class="v-playBtn" href="(?s:(.*?))"`)
	allUrl := reg.FindAllStringSubmatch(rs, -1)

	for i, item := range allUrl {
		fmt.Println("第" + strconv.Itoa(i+1) + "个：" + item[1])
		//处理单个动漫
		//if item[1] == "/vod/detail/id/2337.html" {
		//	continue
		//}
		AnimeItemHandle(i+1, item[1])
	}

}

//处理单个动漫
func AnimeItemHandle(i int, url string) {
	url = "http://pilipali.cc" + url
	rs, err := HttpGet(url)
	if err != nil {
		return
	}
	AnimeData := model.AnimeData{}
	AnimeData = AnimeItemExtractHandle(rs)
	dataChan <- AnimeData

}

//单个动漫数据提取
func AnimeItemExtractHandle(rs string) model.AnimeData {
	//title
	title := extractHandle(rs, `<div class="tit"><h1 class="clearfix">(?s:(.*?))</h1>`, 1)
	//em_num
	em_num := extractHandle(rs, `<p class="p_txt"><em class="em_num">(?s:(.*?))</em>`, 1)
	//year
	year := extractHandle(rs, `target="_blank">(?s:(.*?))</a>&nbsp;&nbsp;&nbsp;<em class="em_tit">`, 1)
	//area
	area := extractHandle(rs, `地区：.* target="_blank">(?s:(.*?))</a>&nbsp;</em>`, 1)
	//star 多个
	stars := extractHandle(rs, `<em class="em_tit">主演：</em>(?s:(.*?))&nbsp;</li>`, -1)
	//director
	director := extractHandle(rs, `导演：.*target="_blank">(?s:(.*?))</a>`, 1)
	//picture
	picture := extractHandle(rs, `class="v-pic"><img src="(?s:(.*?))"`, 1)
	//drama 剧集，多个，不同源
	dramas := extractHandle(rs, `<li><a href="(?s:(.*?))</li>`, -1)
	//introduction 简介
	introduction := extractHandle(rs, `<div class="tv-bd"><div class="infor_intro">(?s:(.*?))</div>`, 1)

	starStr := ""
	if len(stars) > 0 {
		starStr = stars[0]
	}
	star := extractHandle(starStr, `target="_blank">(?s:(.*?))</a>`, -1)
	drama := ExtractDramaHandle(dramas)
	arr := LinkHandle(drama)

	AnimeData := model.AnimeData{}
	if len(title) > 0 {
		AnimeData.Title = title[0]
	}
	if len(em_num) > 0 {
		AnimeData.EmNum = em_num[0]
	}
	if len(year) > 0 {
		AnimeData.Year = year[0]
	}
	if len(area) > 0 {
		AnimeData.Area = area[0]
	}
	AnimeData.Star = star
	AnimeData.Director = director
	if len(picture) > 0 {
		AnimeData.Picture = picture[0]
	}
	if len(introduction) > 0 {
		AnimeData.Introduction = introduction[0]
	}
	fmt.Println(arr)
	AnimeData.Drama = arr

	return AnimeData
}

//提取数据函数
func extractHandle(rs, regStr string, num int) (content []string) {
	reg := regexp.MustCompile(regStr)
	allUrl := reg.FindAllStringSubmatch(rs, num)
	for _, item := range allUrl {
		//content=item[1]
		content = append(content, item[1])
	}
	return
}

//提取每一集
func ExtractDramaHandle(rs []string) (content []model.Dramas) {
	for _, item := range rs {
		url := extractHandle(item, `(?s:(.*?))">`, 1)
		name := extractHandle(item, `">(?s:(.*?))</a>`, 1)
		var drama model.Dramas
		drama.Url = url[0]
		drama.Name = name[0]
		content = append(content, drama)
	}
	return
}

//获取视频链接
func LinkHandle(content []model.Dramas) (DramaData []model.DramaData) {
	for i, item := range content {
		obj := model.DramaData{}
		//处理单个动漫
		obj = ExtractLinkHandle(i, item.Url, item.Name)
		DramaData = append(DramaData, obj)
	}
	return DramaData
}
func ExtractLinkHandle(index int, url, name string) (arr model.DramaData) {
	url = "http://pilipali.cc" + url
	rs, err := HttpGet(url)
	if err != nil {
		return model.DramaData{}
	}
	//drama 剧集，多个，不同源
	dramas := extractHandle(rs, `<script type="text/javascript">var player_data=(?s:(.*?))</script>`, 1)
	obj := model.DramaData{}
	if len(dramas) > 0 {
		var b []byte = []byte(dramas[0])
		var data model.DramaPlay
		err = json.Unmarshal(b, &data)
		tools.CheckErr(err)

		obj.Url = data.Url
		obj.Name = name
		obj.From = data.From
	}
	return obj
}
