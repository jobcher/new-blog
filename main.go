package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chai2010/webp"

	_ "github.com/go-sql-driver/mysql"
	gomail "gopkg.in/gomail.v2"
)

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	// 当前日期
	today := time.Now().Format("2006-01-02")
	md_name := "github_trending_" + today + ".md"

	//判断文件是否存在
	dir := "content/new/daily"
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	// 创建 Markdown 文件
	file, err := os.Create("content/new/daily/" + md_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 下载壁纸
	downloadBingWallpaper()
	// 转换壁纸格式
	tran_webp()

	// 写入 Markdown 文件头部
	_, err = file.WriteString("---\ntitle: " + today + " 打工人日报\ndate: " + today + "\ndraft: false\nauthor: 'jobcher'\nfeaturedImage: '/images/wallpaper/" + today + ".jpg.webp'\nfeaturedImagePreview: '/images/wallpaper/" + today + ".jpg.webp'\nimages: ['/images/wallpaper/" + today + ".jpg.webp']\ntags: ['日报']\ncategories: ['日报']\nseries: ['日报']\n---\n\n")
	if err != nil {
		log.Fatal(err)
	}

	// 获取微博热搜
	// get_weibo(md_name)
	// 获取github热门
	get_github(md_name)
	// 获取v2ex热门
	get_v2ex(md_name)
	// 获取DNSPOD热门
	dnsport_new(md_name)
	// 获取DIYgod热门
	DIY_god(md_name)
	// 获取abskoop热门
	abskoop(md_name)
	// // sitemap 生成
	// get_sitemap()
	// // 发送邮件
	// push_email()

	fmt.Println("成功生成文件")
}

func get_weibo(md_name string) {
	//写入标题
	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("## 微博热搜榜\n\n")

	// 发起 HTTP GET 请求
	res, err := http.Get("https://tophub.today/n/KqndgxeLl9")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("请求失败，状态码：%d", res.StatusCode)
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	var contents []string

	doc.Find(".table tbody tr").Each(func(i int, s *goquery.Selection) {
		count++
		if count > 20 {
			return
		}
		// 提取标题和url
		title := strings.TrimSpace(s.Find("td a").Text())
		url := strings.TrimSpace(s.Find("td a").Text())

		title = strings.Replace(title, "", "", -1)
		url = strings.Replace(url, "", "", -1)
		url = strings.Replace(url, " ", "", -1)

		// 将信息以 Markdown 格式写入文件
		content := fmt.Sprintf("- 排名 %d.", i+1)
		content += fmt.Sprintf("[%s]", title)
		content += fmt.Sprintf("(https://s.weibo.com/weibo?q=%s)\n", url)

		// 将 content 添加到 contents 切片中
		contents = append(contents, content)
	})

	// 将所有的 content 汇总成一个字符串
	allContent := strings.Join(contents, "\n")

	summary := AI_summary(allContent)
	fmt.Println(summary)
	fmt.Println(allContent)

	// 写入 Markdown 文件
	file, err = os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("### AI 摘要\n\n" + summary + "\n\n### 热搜链接\n\n" + allContent)
}

func get_github(md_name string) {
	//写入标题
	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("## GitHub 热门榜单\n\n")

	res, err := http.Get("https://www.github.com/trending")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("请求失败，状态码：%d", res.StatusCode)
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// count := 0
	// 查找所有的 trending repository
	doc.Find(".Box .Box-row").Each(func(i int, s *goquery.Selection) {
		// count++
		// if count > 10 {
		// 	return
		// }
		// 提取标题和作者,title 去除span标签
		title := strings.TrimSpace(s.Find("h2.h3 a").AttrOr("href", ""))
		author := strings.TrimSpace(s.Find("span.text-normal").First().Text())
		url := strings.TrimSpace(s.Find("h2.h3 a").AttrOr("href", ""))
		desc := strings.TrimSpace(s.Find("p.col-9").Text())

		// 去除斜杠
		author = strings.Replace(author, "/", "", -1)
		//翻译
		queryString := desc
		result, err := translateString(queryString)
		if err != nil {
			fmt.Println("翻译失败：", err)
			return
		}
		desc = result

		// 将信息以 Markdown 格式写入文件
		content := fmt.Sprintf("#### 排名 %d:", i+1)
		content += fmt.Sprintf("%s\n", title)
		content += fmt.Sprintf("- 简介: %s\n", desc)
		content += fmt.Sprintf("- URL: https://github.com%s\n", url)
		content += fmt.Sprintf("- 作者: %s\n\n", author)

		fmt.Println(content)

		// 写入 Markdown 文件
		file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(content)
	})
}

func get_v2ex(md_name string) {
	//写入标题
	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("## v2ex 热门帖子\n\n")

	// 发起 HTTP GET 请求
	res, err := http.Get("https://www.v2ex.com/?tab=hot")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("请求失败，状态码：%d", res.StatusCode)
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// count := 0
	// 查找所有的 trending repository
	doc.Find(".cell.item").Each(func(i int, s *goquery.Selection) {
		// count++
		// if count > 20 {
		// 	return
		// }
		// 提取标题和作者,title 去除span标签
		title := strings.TrimSpace(s.Find("span.item_title a").Text())
		url := strings.TrimSpace(s.Find("span.item_title a").AttrOr("href", ""))

		title = strings.Replace(title, " ", "", -1)
		url = strings.Replace(url, " ", "", -1)

		// 将信息以 Markdown 格式写入文件
		content := fmt.Sprintf("- %d.", i+1)
		content += fmt.Sprintf("[%s]", title)
		content += fmt.Sprintf("(https://www.v2ex.com%s)\n", url)

		fmt.Println(content)

		// 写入 Markdown 文件
		file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(content)
	})
}

func DIY_god(md_name string) {
	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("## 热点新闻\n\n")

	rssURLs := []string{
		"https://rsshub.app/telegram/channel/tnews365",
		"https://rssweb.160826.xyz/telegram/channel/tnews365",
		"https://rss.160826.xyz/telegram/channel/tnews365",
	}

	var body []byte
	var rss RSS
	var fetchSuccess bool

	for _, rssURL := range rssURLs {
		fmt.Println("尝试 RSS 源:", rssURL)
		resp, err := http.Get(rssURL)
		if err != nil {
			fmt.Println("请求失败:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 429 {
				fmt.Printf("RSS 源 %s 被限流 (429)，等待后重试...\n", rssURL)
				time.Sleep(5 * time.Second) // 等待5秒后重试
				continue
			}
			fmt.Printf("非 200 状态码: %d\n", resp.StatusCode)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			continue
		}

		if !strings.Contains(string(body), "<?xml") {
			fmt.Println("返回内容不是 XML 格式")
			continue
		}

		err = xml.Unmarshal(body, &rss)
		if err != nil {
			fmt.Println("XML 解析失败:", err)
			continue
		}

		fetchSuccess = true
		break
	}

	if !fetchSuccess {
		fmt.Println("🚫 所有 tnews365 RSS 源均不可用")
		return
	}

	// 获取当前时间，用于比较
	currentTime := time.Now().UTC()
	yesterday := currentTime.AddDate(0, 0, -1)

	// 打印调试信息
	fmt.Printf("Current time: %s\n", currentTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Yesterday: %s\n", yesterday.Format("2006-01-02 15:04:05 UTC"))

	var contents []string
	var titles []string

	for _, item := range rss.Channel.Items {
		// 解析 RSS 条目的发布时间
		itemTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// 尝试其他时间格式
			itemTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", item.PubDate)
			if err != nil {
				// 尝试 RFC822 格式
				itemTime, err = time.Parse(time.RFC822, item.PubDate)
				if err != nil {
					fmt.Printf("无法解析时间格式: %s, 错误: %v\n", item.PubDate, err)
					continue
				}
			}
		}

		// 检查是否是昨天的内容（允许一些时间误差）
		timeDiff := itemTime.Sub(yesterday)
		if timeDiff < -24*time.Hour || timeDiff > 24*time.Hour {
			fmt.Printf("跳过项目，时间不匹配: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))
			continue
		}

		fmt.Printf("匹配项目: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))

		description := strings.ReplaceAll(item.Description, "\n", "")
		content := fmt.Sprintf("#### %s\n%s\n\n", item.Title, description)
		title := fmt.Sprintf("%s\n", item.Title)

		titles = append(titles, title)
		contents = append(contents, content)
	}

	if len(contents) == 0 {
		fmt.Println("⚠️ 没有找到符合时间条件的 tnews365 项目")
		return
	}

	alltitle := strings.Join(titles, "\n")
	allContent := strings.Join(contents, "\n")
	summary := AI_summary(alltitle)

	fmt.Println(summary)
	fmt.Println(allContent)

	file.WriteString("### AI 摘要\n\n" + summary + "\n\n### 热点新闻\n\n" + allContent)
}

func abskoop(md_name string) {
	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("## 福利分享\n\n")

	rssURLs := []string{
		"https://rsshub.app/telegram/channel/abskoop",
		"https://rssweb.160826.xyz/telegram/channel/abskoop",
		"https://rss.160826.xyz/telegram/channel/abskoop",
	}

	var body []byte
	var rss RSS
	var fetchSuccess bool

	for _, rssURL := range rssURLs {
		fmt.Println("尝试 RSS 源:", rssURL)
		resp, err := http.Get(rssURL)
		if err != nil {
			fmt.Println("请求失败:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 429 {
				fmt.Printf("RSS 源 %s 被限流 (429)，等待后重试...\n", rssURL)
				time.Sleep(5 * time.Second) // 等待5秒后重试
				continue
			}
			fmt.Printf("非 200 状态码: %d\n", resp.StatusCode)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			continue
		}

		if !strings.Contains(string(body), "<?xml") {
			fmt.Println("返回内容不是 XML 格式")
			continue
		}

		err = xml.Unmarshal(body, &rss)
		if err != nil {
			fmt.Println("XML 解析失败:", err)
			continue
		}

		fetchSuccess = true
		break
	}

	if !fetchSuccess {
		fmt.Println("🚫 所有 abskoop RSS 源均不可用")
		return
	}

	// 获取当前时间，用于比较
	currentTime := time.Now().UTC()
	yesterday := currentTime.AddDate(0, 0, -1)

	// 打印调试信息
	fmt.Printf("Current time: %s\n", currentTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Yesterday: %s\n", yesterday.Format("2006-01-02 15:04:05 UTC"))

	for _, item := range rss.Channel.Items {
		// 解析 RSS 条目的发布时间
		itemTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// 尝试其他时间格式
			itemTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", item.PubDate)
			if err != nil {
				// 尝试 RFC822 格式
				itemTime, err = time.Parse(time.RFC822, item.PubDate)
				if err != nil {
					fmt.Printf("无法解析时间格式: %s, 错误: %v\n", item.PubDate, err)
					continue
				}
			}
		}

		// 检查是否是昨天的内容（允许一些时间误差）
		timeDiff := itemTime.Sub(yesterday)
		if timeDiff < -24*time.Hour || timeDiff > 24*time.Hour {
			fmt.Printf("跳过项目，时间不匹配: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))
			continue
		}

		fmt.Printf("匹配项目: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))

		description := strings.ReplaceAll(item.Description, "\n", "")
		content := fmt.Sprintf("#### %s\n%s\n\n", item.Title, description)
		fmt.Println(content)

		file.WriteString(content)
	}
}

func dnsport_new(md_name string) {
	// 多个候选 RSS 地址
	rssURLs := []string{
		"https://rsshub.app/telegram/channel/DNSPODT",
		"https://rss.160826.xyz/telegram/channel/DNSPODT",
		"https://rssweb.160826.xyz/telegram/channel/DNSPODT",
	}

	var body []byte
	var rss RSS
	var fetchSuccess bool

	for _, rssURL := range rssURLs {
		fmt.Println("尝试 RSS 源:", rssURL)
		resp, err := http.Get(rssURL)
		if err != nil {
			fmt.Println("请求失败:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 429 {
				fmt.Printf("RSS 源 %s 被限流 (429)，等待后重试...\n", rssURL)
				time.Sleep(5 * time.Second) // 等待5秒后重试
				continue
			}
			fmt.Printf("非 200 状态码: %d\n", resp.StatusCode)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			continue
		}

		if !strings.Contains(string(body), "<?xml") {
			fmt.Println("返回内容不是 XML 格式")
			continue
		}

		err = xml.Unmarshal(body, &rss)
		if err != nil {
			fmt.Println("XML 解析失败:", err)
			continue
		}

		// 成功解析
		fetchSuccess = true
		break
	}

	if !fetchSuccess {
		fmt.Println("🚫 所有 DNSPODT RSS 源均不可用")
		return
	}

	// 获取当前时间，用于比较
	currentTime := time.Now().UTC()
	yesterday := currentTime.AddDate(0, 0, -1)

	// 打印调试信息
	fmt.Printf("Current time: %s\n", currentTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Yesterday: %s\n", yesterday.Format("2006-01-02 15:04:05 UTC"))

	var contents []string
	var titles []string

	for _, item := range rss.Channel.Items {
		// 解析 RSS 条目的发布时间
		itemTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// 尝试其他时间格式
			itemTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", item.PubDate)
			if err != nil {
				// 尝试 RFC822 格式
				itemTime, err = time.Parse(time.RFC822, item.PubDate)
				if err != nil {
					fmt.Printf("无法解析时间格式: %s, 错误: %v\n", item.PubDate, err)
					continue
				}
			}
		}

		// 检查是否是昨天的内容（允许一些时间误差）
		timeDiff := itemTime.Sub(yesterday)
		if timeDiff < -24*time.Hour || timeDiff > 24*time.Hour {
			fmt.Printf("跳过项目，时间不匹配: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))
			continue
		}

		fmt.Printf("匹配项目: %s (发布时间: %s)\n", item.Title, itemTime.Format("2006-01-02 15:04:05 UTC"))

		description := strings.ReplaceAll(item.Description, "\n", "")
		content := fmt.Sprintf("#### %s\n%s\n\n", item.Title, description)
		title := fmt.Sprintf("%s\n", item.Title)

		titles = append(titles, title)
		contents = append(contents, content)
	}

	if len(contents) == 0 {
		fmt.Println("⚠️ 没有找到符合时间条件的 DNSPODT 项目")
		return
	}

	alltitle := strings.Join(titles, "\n")
	allContent := strings.Join(contents, "\n")
	summary := AI_summary(alltitle)

	fmt.Println(allContent)
	fmt.Println(summary)

	file, err := os.OpenFile("content/new/daily/"+md_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("### AI 摘要\n\n" + summary + "\n\n### 热点新闻\n\n" + allContent)
}

type BingResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func downloadBingWallpaper() {
	// 获取当前日期
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")

	// 指定保存目录
	saveDirectory := "assets/images/input/"

	// 构建保存文件路径
	savePath := filepath.Join(saveDirectory, dateString+".jpg")

	// 发起 HTTP 请求获取 Bing 每日壁纸信息
	response, err := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US")
	if err != nil {
		fmt.Println("无法获取壁纸信息:", err)
		return
	}
	defer response.Body.Close()

	// 解析 JSON 数据
	var bingResponse BingResponse
	err = json.NewDecoder(response.Body).Decode(&bingResponse)
	if err != nil {
		fmt.Println("解析壁纸信息失败:", err)
		return
	}

	if len(bingResponse.Images) == 0 {
		fmt.Println("未找到壁纸信息")
		return
	}

	// 获取壁纸 URL
	imageURL := "https://www.bing.com" + bingResponse.Images[0].URL

	// 发起 HTTP 请求下载壁纸
	imageResponse, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("无法下载壁纸:", err)
		return
	}
	defer imageResponse.Body.Close()

	// 创建保存文件
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	// 将壁纸内容保存到文件
	_, err = io.Copy(file, imageResponse.Body)
	if err != nil {
		fmt.Println("保存壁纸失败:", err)
		return
	}

	fmt.Println("壁纸已成功保存到:", savePath)

}

// 翻译
type TranslationResponse struct {
	From         string `json:"from"`
	To           string `json:"to"`
	TransResults []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func translateString(queryString string) (string, error) {
	// 使用环境变量
	apiKey := os.Getenv("BAIDU_TRANSLATE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("未设置百度翻译 API 密钥")
	}
	apiId := os.Getenv("BAIDU_TRANSLATE_API_ID")
	if apiId == "" {
		return "", fmt.Errorf("未设置百度翻译 API ID")
	}
	apiURL := "https://fanyi-api.baidu.com/api/trans/vip/translate"
	salt := "1435660288" // 随机数，这里使用固定值

	// 构建 POST 请求参数
	values := url.Values{}
	values.Set("q", queryString)
	values.Set("from", "en")
	values.Set("to", "zh")
	values.Set("appid", apiId) // 百度翻译 API 的应用ID，固定值
	sign := apiId + queryString + salt + apiKey
	fmt.Println(sign)
	values.Set("salt", salt)
	values.Set("sign", fmt.Sprintf("%x", md5.Sum([]byte(sign))))

	// 发送 POST 请求
	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return "", fmt.Errorf("请求失败：%v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败：%v", err)
	}

	// 解析 JSON 数据
	var response TranslationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析 JSON 失败：%v", err)
	}

	// 提取翻译结果
	if len(response.TransResults) > 0 {
		return response.TransResults[0].Dst, nil
	}

	return "", fmt.Errorf("未找到翻译结果")
}

func tran_webp() {
	// Specify input and output directories
	inputDir := "assets/images/input/"
	outputDir := "assets/images/wallpaper/"

	// Walk input directory to process files
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check for JPEG or PNG file
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		// Open image file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Decode image
		var img image.Image
		if ext == ".jpg" || ext == ".jpeg" {
			img, err = jpeg.Decode(file)
		} else if ext == ".png" {
			img, err = png.Decode(file)
		}
		if err != nil {
			return err
		}

		// Convert to webp
		webpName := filepath.Join(outputDir, filepath.Base(path)+".webp")
		f, _ := os.Create(webpName)
		defer f.Close()

		err = webp.Encode(f, img, &webp.Options{Quality: 50})
		if err != nil {
			return err
		}

		//关闭原文件
		file.Close()

		//清理原文件
		err = os.Remove(path)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

// 发送订阅邮件
func push_email() {
	// 环境变量
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_database := os.Getenv("DB_DATABASE")
	smtp_mail := os.Getenv("SMTP_MAIL")
	smtp_pass := os.Getenv("SMTP_PASS")

	mysql_tcp := db_user + ":" + db_pass + "@tcp(" + db_host + ":" + db_port + ")/" + db_database + "?charset=utf8"

	db, err := sql.Open("mysql", mysql_tcp)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT email FROM subscriptions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	today := time.Now().Format("2006-01-02")
	md_name := "github_trending_" + today

	// 发送邮件
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}

		m := gomail.NewMessage()
		m.SetHeader("From", smtp_mail)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "【打工人日报】 【"+today+"】")
		m.SetBody("text/html", `
		<html>
		<head>
		<style>
		body {font-family: Arial, sans-serif;}
		.container {margin: auto; width: 50%;}
		h1 {color: #333;}
		p {font-size: 16px;}
		a {color: #1a0dab; text-decoration: none;}
		.button {
		  background-color: #4CAF50; /* Green */
		  border: none;
		  color: white;
		  padding: 15px 32px;
		  text-align: center;
		  text-decoration: none;
		  display: inline-block;
		  font-size: 16px;
		  margin: 4px 2px;
		  cursor: pointer;
		}
		</style>
		</head>
		<body>
		<div class="container">
		<h2>打工人日报</h2>
		<p>【`+today+`】</p>
		<p>您订阅的打工人日报已更新，点击下方按钮查看详情。</p>
		<a href='https://www.jobcher.com/new/daily/`+md_name+`/' class='button'>点击查看</a>
		<p>为避免标记为垃圾邮件，请将此邮件地址添加到您的联系人列表。</p>
		<p>如有任何问题，请联系我们。</p>
		<P>取消订阅：<a href='https://sub.jobcher.com/unsubscribe'>https://sub.jobcher.com/unsubscribe</a></p>
		</div>
		</body>
		</html>
		`)

		d := gomail.NewDialer("smtp.qiye.aliyun.com", 25, smtp_mail, smtp_pass)

		if err := d.DialAndSend(m); err != nil {
			log.Println("Failed to send email to", email, ":", err)
		} else {
			fmt.Printf("已发送订阅邮件至 %s\n", email)
		}

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func get_sitemap() {
	rss_url := "https://www.jobcher.com/index.xml"

	// 发送 GET 请求
	resp, err := http.Get(rss_url)
	if err != nil {
		fmt.Println("Error fetching RSS feed:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// 解析 XML 数据
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return
	}

	var contents []string

	// 遍历 RSS 中的条目
	for _, item := range rss.Channel.Items {
		// 提取 URL
		url := item.Link

		cotent := fmt.Sprintf("%s\n", url)

		contents = append(contents, cotent)
	}
	fmt.Println(contents)

	//写入txt
	file, err := os.Create("sitemap.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(contents, "\n"))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Sitemap.txt file created successfully.")

}

func AI_summary(promt string) string {
	cloudflareAccountID := os.Getenv("CLOUDFLARE_ID")    // 你的 Cloudflare 账户 ID
	cloudflareAuthToken := os.Getenv("CLOUDFLARE_TOKEN") // 你的 Cloudflare 授权令牌
	ai_url := "https://api.cloudflare.com/client/v4/accounts/" + cloudflareAccountID + "/ai/run/@cf/qwen/qwen1.5-14b-chat-awq"

	messages := []Message{
		{
			Role:    "system",
			Content: "AI 生成摘要能够生成简洁、有逻辑性的文本摘要。它可以根据输入的文本内容，提取出其中的关键信息，生成易于理解的、精炼的摘要内容，方便用户快速获取文本核心信息。",
		},
		{
			Role:    "user",
			Content: promt,
		},
	}

	data := map[string]interface{}{"messages": messages}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return "摘要生成失败：JSON 序列化错误"
	}

	req, err := http.NewRequest("POST", ai_url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Request creation error: %v\n", err)
		return "摘要生成失败：请求创建错误"
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request execution error: %v\n", err)
		return "摘要生成失败：请求执行错误"
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("API returned status %d: %s\n", resp.StatusCode, string(body))
		return "摘要生成失败：API 返回错误状态码"
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "摘要生成失败：读取响应体错误"
	}

	// 打印原始响应用于调试
	fmt.Printf("Raw API response: %s\n", string(body))

	// 解析 JSON 响应
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("JSON unmarshal error: %v\n", err)
		return "摘要生成失败：JSON 解析错误"
	}

	// 检查是否有错误
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		fmt.Printf("API errors: %v\n", errors)
		return "摘要生成失败：API 返回错误"
	}

	// 尝试从不同的路径获取结果
	var response string

	// 尝试 result.response 路径
	if resultMap, ok := result["result"].(map[string]interface{}); ok {
		if resp, ok := resultMap["response"].(string); ok {
			response = resp
		}
	}

	// 如果上面没找到，尝试直接访问 response 字段
	if response == "" {
		if resp, ok := result["response"].(string); ok {
			response = resp
		}
	}

	// 如果还是没找到，尝试 messages 路径
	if response == "" {
		if messages, ok := result["messages"].([]interface{}); ok && len(messages) > 0 {
			if lastMessage, ok := messages[len(messages)-1].(map[string]interface{}); ok {
				if content, ok := lastMessage["content"].(string); ok {
					response = content
				}
			}
		}
	}

	if response == "" {
		fmt.Printf("Could not extract response from API result: %+v\n", result)
		return "摘要生成失败：无法从 API 响应中提取结果"
	}

	return response
}
