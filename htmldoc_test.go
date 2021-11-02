package exhtml

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/html"
)

var u, err = url.Parse("https://cn.nikkei.com/industry/itelectric-appliance/46280-2021-10-09-01-47-33.html?tmpl=component&print=1&page=")

func TestExtractRssGuids(t *testing.T) {
	src := "https://zh.vietnamplus.vn/rss/news.rss"
	ls, err := ExtractRssGuids(src)
	if err != nil {
		t.Error(err)
	}
	for _, e := range ls {
		fmt.Println(e)
	}
}

func TestExtractRss(t *testing.T) {
	src := "https://china.kyodonews.net/rss/news.xml"
	ls, err := ExtractRss(src)
	if err != nil {
		t.Error(err)
	}
	for _, e := range ls {
		fmt.Println(e)
	}
}

func TestGetRawAndDoc(t *testing.T) {
	raw, doc, err := GetRawAndDoc(u, 2*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(raw))
	fmt.Println(doc)
}

func TestTagWithAttr(t *testing.T) {
	// u, err := url.Parse("https://www.bbc.com/zhongwen/simp/world-55655858")
	u, err := url.Parse("https://zh.vietnamplus.vn/Utilities/Print.aspx?contentid=143985")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := TagWithAttr(doc, "time", "datetime")
	fmt.Println(tc)
	for _, v := range tc {
		for _, a := range v.Attr {
			fmt.Println(a.Key, ":", a.Val)
		}
	}
}

func TestDivWithAttr(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := DivWithAttr(doc, "data-desc", "內容頁")
	plist := ElementsByTag(tc[0], "p")
	for _, v := range plist {
		fmt.Println(v.FirstChild.Data)
	}
}
func TestDivWithAttr2(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := DivWithAttr2(raw, "data-desc", "內容頁")
	fmt.Println(string(tc))
}

func TestElementsByTag(t *testing.T) {
	u, err = url.Parse("https://www.bbc.com/zhongwen/simp/world-55655858")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTag(doc, "main")
	a := ElementsByTag(tc[0], "p")
	for _, v := range a {
		if v.FirstChild != nil {
			fmt.Println(v.FirstChild.Data)
		}
	}
}

func TestElementsByTag2(t *testing.T) {
	// u, err = url.Parse("https://www.dw.com/zh/%E6%8B%9C%E7%99%BB%E5%85%A8%E4%BD%93%E7%BE%8E%E5%9B%BD%E5%85%AC%E6%B0%91%E6%9C%AC%E6%9C%88%E5%B0%B1%E9%83%BD%E8%83%BD%E6%8E%A5%E7%A7%8D%E7%96%AB%E8%8B%97/a-57119062")
	// if err != nil {
	//         t.Errorf("url Parse err: %v", err)
	// }
	// raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	// if err != nil {
	//         t.Errorf("GetRawAndDoc err: %v", err)
	// }
	// a := ElementsByTagAndClass2(raw, "div", "longText")
	// b := ElementsByTag2(a, "p")
	// fmt.Printf("a: %s\n\nb: %s", a, b)

	u, err = url.Parse("https://www.voachinese.com/a/Washington-Post-Chinese-missile-silos-icbm-20210630/5948766.html")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	a := ElementsByTagAndClass2(raw, "div", "wsw")
	b := ElementsByTag2(a, "p")
	fmt.Printf("a: %s\n\nb: %s", a, b)
}

func TestElementsByTagAndClass(t *testing.T) {
	s, err := ioutil.ReadFile("./test.html")
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndClass(doc, "div", "paragraph")
	a := ElementsByTag(tc[0], "h2", "p")
	for _, v := range a {
		if v.FirstChild != nil {
			fmt.Println(v.FirstChild.Data)
		}
	}
}

func TestElementsByTagAndClass2(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndClass2(raw, "div", "wsw")
	fmt.Println(tc)
}

func TestElementsByTagAndId(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndId(doc, "div", "storytext")
	plist := ElementsByTag(tc[0], "p")
	for _, v := range plist {
		if v.FirstChild != nil {
			if v.FirstChild.Data == "b" {
				blist := ElementsByTag(v, "b")
				fmt.Print("**")
				for _, b := range blist {
					fmt.Print(b.FirstChild.Data)
				}
				fmt.Print("**\n")
				// fmt.Println("**" + v.FirstChild.FirstChild.Data + "**")
			} else {
				fmt.Println(v.FirstChild.Data)
			}
		}
	}
}

func TestMetaByProperty(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByProperty(doc, "article:modified_time")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020-08-25T09:42:32+08:00"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestMetaByItemprop(t *testing.T) {
	u, err = url.Parse("https://www.cna.com.tw/news/aopl/202009290075.aspx")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByItemprop(doc, "dateModified")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020/09/29 11:49"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestMetaByName(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByName(doc, "parsely-pub-date")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020-07-09T18:04:00+08:00"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestElementsRmByTag(t *testing.T) {
	s, err := ioutil.ReadFile("./test.html")
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	ElementsRmByTag(doc, "br")
	n := ElementsByTag(doc, "br")
	if len(n) > 0 {
		t.Errorf("want 0, got: %v", len(n))
		fmt.Println(doc)
	}
}

func TestElementsRmByTagClass(t *testing.T) {
	s := []byte(testHtml)
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	ElementsRmByTagClass(doc, "div", "article-photo")
	n := ElementsByTagAndClass(doc, "div", "article-photo")
	if len(n) > 0 {
		t.Errorf("want 0, got: %v", len(n))
		fmt.Println(doc)
	}
}

func TestElementsByTagAttr(t *testing.T) {
	s := []byte(testHtml)
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	ns := ElementsByTagAttr(doc, "span", "data-reactroot", "")
	for _, n := range ns {
		if n.FirstChild == nil {
			t.Errorf("n.FirstChild is nil")
		}
		if n.FirstChild.Data != "Test ElementsByTagAttr" {
			t.Errorf("want: %s, got: %s", "Test ElementsByTagAttr", n.FirstChild.Data)
		}
	}
}

var testHtml = `
<span data-reactroot>Test ElementsByTagAttr</span>
<div class="content article-body" style="width: 100%">
<div class="ExternalClassD4FD9D4E94174110AEAB3E1D8765742B">
<div class="article-photo">
<div class="article-photo"><img alt="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hinh anh 1" src="https://cdnimgzh.vietnamplus.vn/t1000/uploaded/wbxx/2021_08_10/tbt.jpg" title="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hình ảnh 1" class="cms-photo" data-photo-original-src="https://cdnimgzh.vietnamplus.vnhttps://cdnimgzh.vietnamplus.vn/t1000/uploaded/afbb/2021_08_10/tbt.jpg" cms-photo-caption="越共中央总书记阮富仲。图自越通社"/><span>越共中央总书记阮富仲。图自越通社</span></div>
</div>
<br/>
越通社河内——接受越通社驻德国记者采访时，德国马恩（Marx-Engels）基金会主席、德国共产党国际委员会成员斯蒂芬·库纳（Stefan Kühner）认为，越南共产党中央委员会总书记阮富仲署名文章“社会主义理论与实践若干问题和越南走向社会主义的道路”已指明越南走向社会主义的道路。<br/>
<br/>
斯蒂芬·库纳表示，越共中央总书记阮富仲已具体地阐述了越南在特殊条件和特点中走向社会主义的道路。阮富仲在文章中强调：“跨过资本主义制度是跨过资本主义的压迫、不公平和剥削制度，放弃不符合社会主义制度的不良风气和政治体制，而不是放弃人类在资本主义发展时期所取得的文明成就与价值。当然，继承这些成就需要本着科学和发展观点去选择。” 斯蒂芬·库纳认为，对全球共产主义者来说，该论点已阐明，越南一直在寻找的是走向社会主义道路。<br/>
<div class="article-photo"><img alt="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hinh anh 2" src="https://cdnimgzh.vietnamplus.vn/t1000/uploaded/wbxx/2021_08_10/chuyengia.jpg" title="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hình ảnh 2" class="cms-photo" data-photo-original-src="https://cdnimgzh.vietnamplus.vnhttps://cdnimgzh.vietnamplus.vn/t1000/uploaded/afbb/2021_08_10/chuyengia.jpg" cms-photo-caption="德国马恩基金会主席斯蒂芬·库纳。图自越通社"/><span>德国马恩基金会主席斯蒂芬·库纳。图自越通社</span></div>
<br/>
斯蒂芬·库纳指出，阮富仲总书记在文章中已突出了越南战争后困难，在遭受破败、禁运等的条件下所做出的努力以及越南当前的发展情况。阮富仲指出社会主义建设不能教条化，社会主义制度不仅是通过选举建设政府，同时需要人民的力量，为人民服务。此外，阮富仲还指明越南“社会主义定向市场经济体制”这一概念。<br/>
<br/>
阮富仲总书记署名文章进一步指明越南走向社会主义的道路，同时明确生产力的发展对公正、平等的经济体建设注入重要动力。<br/>
<br/>
谈及越南共产党的作用，斯蒂芬·库纳强调，这是越南革命取得胜利的决定性因素。他同时表示，在德国共产党第二十三次代表大会上，德国共产党已强调加强各国共产党之间以及国际工人之间的合作有助于巩固全球范围内的革命运动。（完）</div>
<div style="text-align: right" class="cms-author"><strong> 越通社</strong></div>
</div>
<div class="content article-body" style="width: 100%">
<div class="ExternalClassD4FD9D4E94174110AEAB3E1D8765742B">
<div class="article-photo">
<div class="article-photo"><img alt="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hinh anh 1" src="https://cdnimgzh.vietnamplus.vn/t1000/uploaded/wbxx/2021_08_10/tbt.jpg" title="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hình ảnh 1" class="cms-photo" data-photo-original-src="https://cdnimgzh.vietnamplus.vnhttps://cdnimgzh.vietnamplus.vn/t1000/uploaded/afbb/2021_08_10/tbt.jpg" cms-photo-caption="越共中央总书记阮富仲。图自越通社"/><span>越共中央总书记阮富仲。图自越通社</span></div>
</div>
<br/>
越通社河内——接受越通社驻德国记者采访时，德国马恩（Marx-Engels）基金会主席、德国共产党国际委员会成员斯蒂芬·库纳（Stefan Kühner）认为，越南共产党中央委员会总书记阮富仲署名文章“社会主义理论与实践若干问题和越南走向社会主义的道路”已指明越南走向社会主义的道路。<br/>
<br/>
斯蒂芬·库纳表示，越共中央总书记阮富仲已具体地阐述了越南在特殊条件和特点中走向社会主义的道路。阮富仲在文章中强调：“跨过资本主义制度是跨过资本主义的压迫、不公平和剥削制度，放弃不符合社会主义制度的不良风气和政治体制，而不是放弃人类在资本主义发展时期所取得的文明成就与价值。当然，继承这些成就需要本着科学和发展观点去选择。” 斯蒂芬·库纳认为，对全球共产主义者来说，该论点已阐明，越南一直在寻找的是走向社会主义道路。<br/>
<div class="article-photo"><img alt="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hinh anh 2" src="https://cdnimgzh.vietnamplus.vn/t1000/uploaded/wbxx/2021_08_10/chuyengia.jpg" title="德国马恩基金会主席：越共中央总书记阮富仲署名文章指明越南走向社会主义的道路 hình ảnh 2" class="cms-photo" data-photo-original-src="https://cdnimgzh.vietnamplus.vnhttps://cdnimgzh.vietnamplus.vn/t1000/uploaded/afbb/2021_08_10/chuyengia.jpg" cms-photo-caption="德国马恩基金会主席斯蒂芬·库纳。图自越通社"/><span>德国马恩基金会主席斯蒂芬·库纳。图自越通社</span></div>
<br/>
斯蒂芬·库纳指出，阮富仲总书记在文章中已突出了越南战争后困难，在遭受破败、禁运等的条件下所做出的努力以及越南当前的发展情况。阮富仲指出社会主义建设不能教条化，社会主义制度不仅是通过选举建设政府，同时需要人民的力量，为人民服务。此外，阮富仲还指明越南“社会主义定向市场经济体制”这一概念。<br/>
<br/>
阮富仲总书记署名文章进一步指明越南走向社会主义的道路，同时明确生产力的发展对公正、平等的经济体建设注入重要动力。<br/>
<br/>
谈及越南共产党的作用，斯蒂芬·库纳强调，这是越南革命取得胜利的决定性因素。他同时表示，在德国共产党第二十三次代表大会上，德国共产党已强调加强各国共产党之间以及国际工人之间的合作有助于巩固全球范围内的革命运动。（完）</div>
<div style="text-align: right" class="cms-author"><strong> 越通社</strong></div>
</div>
`
