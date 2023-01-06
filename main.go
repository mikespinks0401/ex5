package main

/*
1.parse page
2.function shoudl take in (*htmlNode, links map, originalDomainName, linkToParse)
2.get all links[filter out links that don't containt domain name or begin with "/"]
3.add links to Links map
4.parse links by concatonating to initial link


*/
import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"golang.org/x/net/html"
)


type SiteMap struct{
	Domain string
	VisitedPaths []string
	UnvisitedPaths []string
}


func main(){
	urlFlag := flag.String("url", "", "Url to be crawled")
	flag.Parse()
	if *urlFlag == "" {
	panic("must provide a url to crawl")
	}

	url := *urlFlag
	fmt.Println("URL Provided -",url)
	siteMap := SiteMap{}
	siteMap.addDomain(url)
	htmlNode := getHtmlNodeFromDomain(siteMap.Domain)
	getLinksFromNode(htmlNode, &siteMap)
	keepGoing := true
	for keepGoing {
		val := siteMap.UnvisitedPaths[0]
		link := fmt.Sprintf("https://%s/%s",siteMap.Domain, val)
		if siteMap.VisitPath(val){
			continue
		}
		htmlNode := getHtmlNodeFromDomain(link)
		getLinksFromNode(htmlNode, &siteMap)
		if len(siteMap.UnvisitedPaths) == 0{
			keepGoing = false
		}
	}
	for _,val := range siteMap.VisitedPaths{
		fmt.Println("visited-",val)
	}
}

func getHtmlNodeFromDomain(domain string)*html.Node{
	parseLink := ""
	if !strings.Contains(domain, "http") || !strings.Contains(domain, "https"){
		parseLink = fmt.Sprintf("https://%s", domain)
	}
	if parseLink == ""{
		parseLink = domain
	}
	resp, err := http.Get(parseLink)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	htmlNode, err := html.Parse(resp.Body)
	if err != nil {
			fmt.Println(nil)
		}
		return htmlNode
	}

func getLinksFromNode(n *html.Node,siteMap *SiteMap){
	if	n == nil{
		return
	}
	if n.Type != html.ElementNode && n.Type != html.DocumentNode {
		return
	}
	if n.Data == "a"{
		for _,val := range n.Attr{ 
			if val.Key == "href"{
				fmt.Printf("%s:%s\n",val.Key,val.Val)
				if strings.Contains(val.Val, siteMap.Domain) || val.Val[0] == '/'{
					path := siteMap.getPath(val.Val) 
					err := removeTrailingSlash(&path)
					if err != nil {
						return	
					}
					if path == ""{
						return
					}
					if siteMap.isUnvisited(path) && !siteMap.isVisited(path){
						if path != ""{
						siteMap.addPath(path)
						}
					}
				}  		
			}
		}
	}
	c := n.FirstChild
	for c != nil {
		getLinksFromNode(c, siteMap)
		c = c.NextSibling
	}

	
}


func (m *SiteMap)addPath(path string){
	if !m.isUnvisited(path) {
		fmt.Println(path, "Already exist")
		return
	}
	newPath := ""
	if path[len(path) - 1] == '/'{
		newPath = path[:len(path) - 1]
	}else {
		newPath = path
	}

	m.UnvisitedPaths = append(m.UnvisitedPaths, newPath)
}

func (m *SiteMap)isVisited(path string)bool{
	//return a bool if the path already exist in visied
	for _,val := range m.VisitedPaths{
		if val == path{
			return true
		}
	}
	return false
}

func (m *SiteMap)isUnvisited(path string)bool{
	for _,val := range m.UnvisitedPaths{
		if val == path || val == path + "/"{
			return false
		}
	}
	return true
}

func (m *SiteMap)VisitPath(url string) bool{

	//splits url by domain name -> checks if the path exist in visited paths -> add path to array if it does not exist
	//remove path from unvisited paths array
	
	splitDomain := strings.Split(url, m.Domain) 
	path := splitDomain[len(splitDomain) - 1]
	

	if m.isVisited(path){
		return true
	}
		
	m.VisitedPaths = append(m.VisitedPaths, path)
	unvisitedIndex := -1
	for ind,val := range m.UnvisitedPaths{
		if val == path{
			unvisitedIndex = ind
		}
	}
	if unvisitedIndex == -1 {
		panic("error")
	}
	m.UnvisitedPaths = append(m.UnvisitedPaths[:unvisitedIndex], m.UnvisitedPaths[unvisitedIndex + 1:]...)
	return false
}
		


func (m *SiteMap)addDomain(domain string){
	if strings.Contains(domain, "//"){
		splitDomain := strings.Split(domain, "//")
		domain := splitDomain[len(splitDomain) - 1]
		m.Domain = domain
		return
	}
	m.Domain = domain
}

func (m *SiteMap)getPath(href string) string{
	path := ""
	if strings.Contains(href, "http") && strings.Contains(href, m.Domain) || strings.Contains(href, "https") && strings.Contains(href, m.Domain){
		splitHref := strings.Split(href, m.Domain)
		path = splitHref[len(splitHref) - 1]
		path = path[1:]
	}else if href[0] == '/'{
		path = href[1:]
	}
	return path 	
}

func removeTrailingSlash(path *string)error{
	if len(strings.TrimSpace(*path)) == 0{
		return fmt.Errorf("cannot format empty string")
	}
	newPath := strings.TrimSpace(*path)
	length := len(newPath)
	if string(newPath[length - 1]) == "/"{
		newPath = newPath[:length-1]
	}
	*path = newPath
	return nil
}