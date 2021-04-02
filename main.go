// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"strconv"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

const (
	defaultURL          string = `https://www.whatsmyua.info/?a`
	usageURL            string = "the URL to be scrapped"
	usageMin            string = "the minimum screen width"
	usageMax            string = "the maximum  screen width"
	usageHeight         string = "the default screen height"
	usageInc            string = "the screen width increment"
	DOMNodesSelector    string = ":is(div, a, form, img, li, h1, h2, h3)"
	minimumScreenWidth  string = "800"
	maximumScreenWidth  string = "1900"
	defaultScreenHeight string = "750"
	defaultWidthInc     string = "100"
)

type ViewPort struct {
	screenWidth  int64
	screenHeight int64
}

type WidthVariation []int
type FullReport map[string]WidthVariation

type ElementAnalysis struct {
	ctx context.Context
	ViewPort
	headerNodes []*cdp.Node
	ids         []cdp.NodeID
	report      FullReport
}

func getMidZoneHeight(bottom, top int) int {
	if bottom < top {
		return 0
	}
	return int(top + (bottom - top)/2)
}

func (e *ElementAnalysis) getElementHeights(ctx context.Context) error {
	var res, output string
LOOP:
	for _, node := range e.headerNodes {
		res = ""
		res += fmt.Sprintf("-------------------- %s ------------------ \n", node.NodeName)
		res += fmt.Sprintf("Node: %d %s %s %s %s %+v %s\n",
			node.NodeID,
			node.Name,
			node.Value,
			node.NodeName,
			node.NodeValue,
			node.Attributes,
			node.FullXPath())

		res += fmt.Sprintf("\tChildren (%d): %+v \n", node.ChildNodeCount, node.Children)
		if node.ChildNodeCount == 0 {
			html, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				continue LOOP
			}
			res += fmt.Sprintf("%s\n", html)
		}
		quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
		if err != nil {
			continue LOOP
		}
		mid := getMidZoneHeight(int(quads[0][5]), int(quads[0][1]))
		e.report[node.FullXPath()] = append(e.report[node.FullXPath()], mid)
		res += fmt.Sprintf("\tPosition top/bottom vertical element position: %d [%d px.. %d px]\n", mid, int(quads[0][1]), int(quads[0][5]))
		output += res
	}

	filename := fmt.Sprintf("%dx%d.dom", e.screenWidth, e.screenHeight)
	err := ioutil.WriteFile(filename, []byte(output), 0o644)
	if err != nil {
		return errors.New("write file " + filename + ": " + err.Error())
	}

	return nil
}

func (e *ElementAnalysis) analyzeElementsHeights() error {
	var err error
	fmt.Printf("analyzeElementsHeights screen dimension : %d x %d\n", e.screenWidth, e.screenHeight)
	err = chromedp.Run(e.ctx,
		chromedp.EmulateViewport(e.screenWidth, e.screenHeight),
		chromedp.Nodes(DOMNodesSelector, &e.headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(e.getElementHeights),
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *ElementAnalysis) printCSV() error {
	var err error
	var output string
	fmt.Printf("\n%+v\n", e.report)
	filename := "output.csv"
	for xpath, dim := range e.report {
		output += xpath
		for _, y := range dim {
			output += ", " + strconv.Itoa(y)
		}
		output += "\n"
	}
	err = ioutil.WriteFile(filename, []byte(output), 0o644)
	if err != nil {
		return errors.New("write file " + filename + ": " + err.Error())
	}
	return nil
}

func main() {
	var err error
	var nWidth int
	start := time.Now()
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var minScreenWidth, maxScreenWidth, screenHeight, widthInc, url string
	var minWidth, maxWidth, height, inc int

	e := ElementAnalysis{
		headerNodes: []*cdp.Node{},
		ctx:         ctx,
		report:      make(FullReport),
	}
	flag.StringVar(&minScreenWidth, "min", minimumScreenWidth, usageMin)
	flag.StringVar(&minScreenWidth, "m", minimumScreenWidth, usageMin+" (shorthand)")
	flag.StringVar(&maxScreenWidth, "max", maximumScreenWidth, usageMax)
	flag.StringVar(&maxScreenWidth, "M", maximumScreenWidth, usageMax+" (shorthand)")
	flag.StringVar(&screenHeight, "height", defaultScreenHeight, usageHeight)
	flag.StringVar(&screenHeight, "H", defaultScreenHeight, usageHeight+" (shorthand)")
	flag.StringVar(&widthInc, "inc", defaultWidthInc, usageInc)
	flag.StringVar(&widthInc, "i", defaultWidthInc, usageInc+" (shorthand)")

	flag.StringVar(&url, "url", defaultURL, usageURL)
	flag.StringVar(&url, "u", defaultURL, usageURL+" (shorthand)")
	flag.Parse()

	fmt.Printf("Open URL %s\n", url)
	err = chromedp.Run(ctx,
		emulation.SetUserAgentOverride("JT-WebScraper-1.0"),
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}

	minWidth, err = strconv.Atoi(minScreenWidth)
	if err != nil {
		log.Fatal(err)
	}
	maxWidth, err = strconv.Atoi(maxScreenWidth)
	if err != nil {
		log.Fatal(err)
	}
	inc, err = strconv.Atoi(widthInc)
	if err != nil {
		log.Fatal(err)
	}
	height, err = strconv.Atoi(screenHeight)
	if err != nil {
		log.Fatal(err)
	}

	for m := minWidth; m < maxWidth; m += inc {
		nWidth++
		e.screenWidth = int64(m)
		e.screenHeight = int64(height)
		err = e.analyzeElementsHeights()
		if err != nil {
			log.Fatal(err)
		}
	}

	var ids []cdp.NodeID
	err = chromedp.Run(ctx,
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			fmt.Println("---------------------")
			fmt.Println("Get body outerHTML before modification")
			fmt.Println("---------------------")

			// export html for debug
			res, err := dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			if err != nil {
				return errors.New("get outer html from node: " + err.Error())
			}
			h := sha1.New()
			h.Write([]byte(res))
			bs := h.Sum(nil)
			sha := fmt.Sprintf("%x", bs)
			fmt.Println("Output " + sha + ".html")
			err = ioutil.WriteFile(sha+".html", []byte(res), 0o644)
			if err != nil {
				return errors.New("save layout " + sha + ".html: " + err.Error())
			}
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}


	err = e.printCSV()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTook: %f secs for %d different width\n", time.Since(start).Seconds(), nWidth)
}

func traverse(ctx context.Context, node *cdp.Node) error {
	var err error
	err = dom.RequestChildNodes(cdp.NodeID(node.BackendNodeID)).WithDepth(-1).Do(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("error while requesting child nodes from node id %d: %s", node.NodeID, err.Error()))
	}
	fmt.Printf("[%d] %d %s %s %d %+v\n", node.BackendNodeID, node.NodeType, node.NodeName, node.NodeValue, node.ChildNodeCount, node.Children)
	nodeParams, err := dom.DescribeNode().WithBackendNodeID(node.BackendNodeID).WithDepth(1).WithPierce(true).Do(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("error while describing node id %d: %s", node.BackendNodeID, err.Error()))
	}
	for _, n := range nodeParams.Children {
		//for _, n := range node.Children {
		err = traverse(ctx, n)
		if err != nil {
			return errors.New(fmt.Sprintf("error while traversing child node %d from node %d: %s", n.BackendNodeID, node.BackendNodeID, err.Error()))
		}
	}
	return nil
}
