package main

import (
	"bufio"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type FileInfoData struct {
	FileInfos []FileInfo
}

type FileInfo struct {
	Uri      string
	FileName string
	Encoding string
}

var (
	inFile   *string
	outDir   *string
	showHelp *bool
	inURL    *string
)

func init() {
	inFile = flag.String("file", "", "file to consume expects image tags")
	inURL = flag.String("url", "", "url to consume expects image tags")
	outDir = flag.String("outdir", "./", "directory to output")
	showHelp = flag.Bool("help", false, "show all feature flags")
}

func findMojis(url string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

}

func getResultAs(url *string) (base64 string, err error) {
	resp, err := http.Get(*url)

	if err != nil {
		fmt.Printf("can't get url %v\n", url)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading request body")
		return
	}
	base64 = b64.StdEncoding.EncodeToString(body)
	return base64, nil
}

func processFile(filename *string) (files []FileInfo, err error) {

	file, err := os.Open(*filename)
	defer file.Close()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	re := regexp.MustCompile(`^https://.*\/(.*)\?.*`)

	for _, line := range lines {
		re.FindStringSubmatch(line)
		found := re.FindStringSubmatch(line)

		if len(found) <= 1 {
			continue
		}

		files = append(files, FileInfo{
			Uri:      line,
			FileName: found[1],
		})

		fmt.Printf("%s, %s\n", line, found[1])
	}

	return

}

func main() {

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		return
	}

	if *inFile != "" {
		// fmt.Println(*inFile)
		files, _ := processFile(inFile)

		for i, file := range files {
			files[i].Encoding, _ = getResultAs(&file.Uri)
			fmt.Println(i)
		}
		// files := []FileInfo{
		// 	FileInfo{
		// 		FileName: "turd",
		// 		Encoding: "encodelol",
		// 	},
		// 	FileInfo{
		// 		FileName: "moo",
		// 		Encoding: "wat",
		// 	},
		// }

		tmpl := template.Must(template.ParseFiles("./asset.css.tmpl"))
		f, err := os.Create("out.scss")
		if err != nil {
			panic(err)
		}

		// w := bufio.NewWriter(f)

		err = tmpl.Execute(f, FileInfoData{files})

		if err != nil {
			panic(err)
		}
	}

	if *inURL != "" {
		findMojis(*inURL)
	}

}
