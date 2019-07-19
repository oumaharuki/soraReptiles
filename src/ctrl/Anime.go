package ctrl

import (
	"encoding/json"
	"fmt"
	"model"
	"regexp"
	"strconv"
	"tools"
)

func AnimeHandle(rs string) {
	reg := regexp.MustCompile(`<a class="v-playBtn" href="(?s:(.*?))"`)
	allUrl := reg.FindAllStringSubmatch(rs, -1)

	anime := make(chan int)
	for i, item := range allUrl {
		fmt.Println("第" + strconv.Itoa(i+1) + "个：" + item[1])
		//处理单个动漫
		go AnimeItemHandle(i+1, item[1], anime)
	}
	for i := 0; i <= len(allUrl); i++ {
		fmt.Println(<-anime)
	}
}

//处理单个动漫
func AnimeItemHandle(i int, url string, anime chan int) {
	url = "http://pilipali.cc" + url
	rs, err := HttpGet(url)
	if err != nil {
		return
	}
	AnimeData := model.AnimeData{}
	AnimeData = AnimeItemExtractHandle(rs)

	Picture := SaveImg(AnimeData.Picture)
	go Save2Mysql(AnimeData, Picture)

	anime <- i
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
