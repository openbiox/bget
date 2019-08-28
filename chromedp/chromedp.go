package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openbiox/butils/log"

	cdp "github.com/chromedp/chromedp"
)

func Chrome2URLs(url string) []string {
	// create context
	ctx, cancel := cdp.NewContext(context.Background())
	defer cancel()
	var err error
	var attbs []map[string]string
	urls := []string{}
	// run task list
	//err := cdp.Run(ctx, visibleNejm("https://www.nejm.org/doi/full/10.1056/NEJMoa1902226", &attbs))
	if strings.Contains(url, "www.nejm.org") {
		err = cdp.Run(ctx, visibleNejm(url, &attbs))
	} else if strings.Contains(url, "sciencedirect.com") {
		err = cdp.Run(ctx, visibleScienceDirect(url, &attbs))
	} else if strings.Contains(url, "www.ncbi.nlm.nih.gov/Traces/study") {
		err = cdp.Run(ctx, visibleSraRunSelect(url, &attbs, ctx))
	}
	if err != nil {
		log.Fatal(err)
	}
	for i := range attbs {
		for k, v := range attbs[i] {
			if k == "href" {
				urls = append(urls, v)
			} else if strings.Contains(k, "http") {
				urls = append(urls, v)
			}
		}
	}
	return urls
}
func visibleScienceDirect(host string, attbs *[]map[string]string) cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(host),
		cdp.WaitVisible(`.show-toc-button`, cdp.ByQuery),
		cdp.Click(`.show-toc-button`, cdp.ByQuery),
		cdp.WaitVisible(`a[href="#app2"]`, cdp.ByQuery),
		cdp.Click(`a[href="#app2"]`, cdp.ByQuery),
		cdp.WaitVisible(`#app2`, cdp.ByQuery),
		cdp.AttributesAll(".Appendices a.icon-link[href]", attbs, cdp.ByQueryAll),
	}
}
func visibleNejm(host string, attbs *[]map[string]string) cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(host),
		cdp.WaitVisible(`#article_supplementary_material`, cdp.ByID),
		cdp.Click("#article_supplementary_material", cdp.ByID),
		cdp.ActionFunc(func(context.Context) error {
			return nil
		}),
		cdp.WaitVisible(`td a[data-download-type='Supplementary Protocol']`, cdp.ByQuery),
		cdp.ActionFunc(func(context.Context) error {
			return nil
		}),
		cdp.AttributesAll(".o-article-body__collapsible-content td a[data-interactiontype=multimedia_download]", attbs, cdp.ByQueryAll),
		cdp.ActionFunc(func(context.Context) error {
			return nil
		}),
	}
}

func visibleSraRunSelect(url string, attbs *[]map[string]string, ctx context.Context) cdp.Tasks {
	(*attbs) = append((*attbs), make(map[string]string))
	(*attbs) = append((*attbs), make(map[string]string))
	return cdp.Tasks{
		cdp.Navigate(url),
		cdp.ActionFunc(func(context.Context) error {
			log.Infof("Visiting %s", url)
			return nil
		}),
		cdp.WaitVisible(`#t-rit-all`, cdp.ByQuery),
		cdp.Click(`#t-rit-all`, cdp.ByQuery),
		cdp.Sleep(time.Second * 5),
		cdp.ActionFunc(func(context.Context) error {
			tars, _ := cdp.Targets(ctx)
			for _, t := range tars {
				if strings.Contains(t.URL, "study/backends") && strings.Contains(t.URL, "assay_type") {
					(*attbs)[0][t.URL] = t.URL
				}
			}
			return nil
		}),
		cdp.WaitVisible(`#t-acclist-all`, cdp.ByQuery),
		cdp.Click(`#t-acclist-all`, cdp.ByQuery),
		cdp.Sleep(time.Second * 5),
		cdp.ActionFunc(func(context.Context) error {
			tars, _ := cdp.Targets(ctx)
			for _, t := range tars {
				if strings.Contains(t.URL, "study/backends") && strings.Contains(t.URL, "assay_type") {
					(*attbs)[1][t.URL] = t.URL
				}
			}
			return nil
		}),
		cdp.Sleep(time.Second * 2),
	}
}

func main() {
	fmt.Println(Chrome2URLs("https://www.ncbi.nlm.nih.gov/Traces/study/?acc=PRJNA133369"))
}
