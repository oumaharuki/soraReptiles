package ctrl

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func workChan(index int,page chan int)  {
	url:="http://pilipali.cc/vod/show/id/4/page/"+strconv.Itoa(index)+".html"
	//得到每一页动漫数据
	rs,err:=HttpGet(url)
	if err!=nil{
		return
	}

	AnimeHandle(rs)

	//AnimeHandle(rs)

	//f,err:=os.Create("第 "+strconv.Itoa(i)+" 页.html")
	//if err!=nil{
	//	return
	//}
	//
	//f.WriteString(rs)
	//f.Close()
	page<-index
}
func Work(start,end int)  {
	page:=make(chan int)
	for i:=start;i<=end;i++{
		go workChan(i,page)
	}
	for i:=start;i<=end;i++{
		fmt.Printf("第%d页结束\n",<-page)
	}
}
func HttpGet(url string)(rs string,err error)  {
	resp,err:=http.Get(url)
	if err!=nil{
		return
	}
	defer resp.Body.Close()

	buf:=make([]byte,4096)
	for {
		n,err1:=resp.Body.Read(buf)
		if n==0{
			break
		}
		if err1!=nil&&err1!=io.EOF{
			err=err1
			return
		}

		rs +=string(buf[:n])
	}
	return
}