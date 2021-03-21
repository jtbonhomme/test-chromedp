// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"crypto/sha1"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

const (
	defaultURL = `https://www.whatsmyua.info/?a`
	usageURL   = "the URL to be scrapped"
)

type ElementAnalysis struct {
	filename     string
	screenWidht  int
	screenHeigth int
	headerNodes  []*cdp.Node
}

func (e *ElementAnalysis) getElementHeights(ctx context.Context) error {
	var res string
	for _, node := range e.headerNodes {
		res += fmt.Sprintf("-------------------- %s ------------------ \n", node.NodeName)
		res += fmt.Sprintf("Node: %d %s %s %s %s %+v %s %s\n",
			node.NodeID,
			node.Name,
			node.Value,
			node.NodeName,
			node.NodeValue,
			node.Attributes,
			node.PartialXPath(),
			node.FullXPath())
		res += fmt.Sprintf("\tChildren (%d): %+v \n", node.ChildNodeCount, node.Children)
		if node.ChildNodeCount == 0 {
			html, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				res += fmt.Sprintf("\t%s\n", err.Error())
			}
			res += fmt.Sprintf("%s\n", html)
		}
		quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
		if err != nil {
			res += fmt.Sprintf("\t%s\n", err.Error())
			continue
		}
		res += fmt.Sprintf("\tPosition top/bottom vertical element position: [%d px.. %d px]\n", int(quads[0][1]), int(quads[0][5]))
	}
	err := ioutil.WriteFile(e.filename, []byte(res), 0o644)
	if err != nil {
		return errors.New("write file " + e.filename + ": " + err.Error())
	}
	return nil
}

func main() {
	start := time.Now()
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var url string
	/*e := ElementAnalysis{
		headerNodes: []*cdp.Node{},
	}*/

	flag.StringVar(&url, "url", defaultURL, usageURL)
	flag.StringVar(&url, "u", defaultURL, usageURL+" (shorthand)")
	flag.Parse()
	// run
	fmt.Printf("Open URL %s\n", url)
	var ids []cdp.NodeID
	var bodys []*cdp.Node
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("JT-WebScraper-1.0"),
		// set viewport
		chromedp.EmulateViewport(780, 651),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		/*chromedp.ActionFunc(func(c context.Context) error {
			e.filename = "780x651.dom"
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			//return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
			return nil
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &e.headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(e.getElementHeights),
		// set viewport
		chromedp.EmulateViewport(1280, 720),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			e.filename = "1280x720.dom"
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			//return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
			return nil
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &e.headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(e.getElementHeights),
		// set viewport
		chromedp.EmulateViewport(1900, 1280),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			e.filename = "1900x1280.dom"
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			//return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
			return nil
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &e.headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(e.getElementHeights),
		*/chromedp.ActionFunc(func(ctx context.Context) error {
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
		chromedp.Nodes(`body`, &bodys, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context /*h cdptypes.Handler*/) error {
			var err error

			fmt.Println("---------------------")
			fmt.Println("Get body outerHTML after modification")
			fmt.Println("---------------------")
			err = traverse(ctx, bodys[0])
			if err != nil {
				return errors.New("failed traverse dom: " + err.Error())
			}

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

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
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
