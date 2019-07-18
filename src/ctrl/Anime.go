package ctrl

import (
	"fmt"
	"model"
	"regexp"
	"strconv"
)

func AnimeHandle(rs string) {
	reg := regexp.MustCompile(`<a class="v-playBtn" href="(?s:(.*?))"`)
	allUrl := reg.FindAllStringSubmatch(rs, 1)
	page := make(chan int)
	for i, item := range allUrl {
		fmt.Println("第" + strconv.Itoa(i) + "个：" + item[1])
		//处理单个动漫
		go AnimeItemHandle(i, item[1], page)
	}
	for i := 0; i <= len(allUrl); i++ {
		fmt.Printf("第%d个结束\n", <-page)
	}
}

//处理单个动漫
func AnimeItemHandle(i int, url string, page chan int) {
	url = "http://pilipali.cc" + url
	rs, err := HttpGet(url)
	if err != nil {
		return
	}
	AnimeItemExtractHandle(rs)
	page <- i
}

//单个动漫数据提取
func AnimeItemExtractHandle(rs string) {
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
	dramas := extractHandle(rs, `<li><a href="(?s:(.*?))</a></li>`, -1)
	//introduction 简介
	introduction := extractHandle(rs, `<div class="tv-bd"><div class="infor_intro">(?s:(.*?))</div>`, 1)

	starStr := ""
	if len(stars) > 0 {
		starStr = stars[0]
	}
	star := extractHandle(starStr, `target="_blank">(?s:(.*?))</a>`, -1)
	drama := ExtractDramaHandle(dramas)
	arr := LinkHandle(drama)
	fmt.Println(title)
	fmt.Println(em_num)
	fmt.Println(year)
	fmt.Println(area)
	fmt.Println(star)
	fmt.Println(director)
	fmt.Println(picture)
	fmt.Println(drama)
	fmt.Println(introduction)
	fmt.Println(arr)
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
func ExtractDramaHandle(rs []string) (content []model.Drama) {
	for _, item := range rs {
		url := extractHandle(item, `(?s:(.*?))">`, 1)
		name := extractHandle(item, `">(?s:(.*?))`, 1)
		var drama model.Drama
		drama.Url = url[0]
		drama.Name = name[0]
		content = append(content, drama)
	}
	return
}

//获取视频链接
func LinkHandle(content []model.Drama) (arr []model.Drama) {
	for i, item := range content {
		fmt.Println(item.Name)
		//处理单个动漫
		url := ExtractLinkHandle(i, item.Url)
		var drama model.Drama
		if len(url) > 0 {
			drama.Url = url[0]
		}
		drama.Name = item.Name
		arr = append(content, drama)
	}
	return
}
func ExtractLinkHandle(index int, url string) []string {
	url = "http://pilipali.cc" + url
	rs, err := HttpGet(url)
	if err != nil {
		return []string{}
	}
	//drama 剧集，多个，不同源
	dramas := extractHandle(rs, `<script type="text/javascript">var player_data=(?s:(.*?))</script>`, 1)
	fmt.Println("dramas:", dramas)
	return dramas
}
