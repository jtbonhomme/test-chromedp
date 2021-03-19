// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"flag"
	"fmt"

	"crypto/sha1"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

const (
	defaultURL = `https://www.whatsmyua.info/?a`
	usageURL   = "the URL to be scrapped"
)

func main() {
	start := time.Now()
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var url string

	flag.StringVar(&url, "url", defaultURL, usageURL)
	flag.StringVar(&url, "u", defaultURL, usageURL+" (shorthand)")
	flag.Parse()
	// run
	fmt.Printf("Open URL %s\n", url)
	var ids []cdp.NodeID
	var headerNodes []*cdp.Node
	err := chromedp.Run(ctx,
		// set viewport
		chromedp.EmulateViewport(780, 651),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range headerNodes {
				/*				fmt.Printf("-------------------- %s ------------------ \n", node.NodeName)
								fmt.Printf("Node: %d %s %s %s %s %+v %s %s\n",
									node.NodeID,
									node.Name,
									node.Value,
									node.NodeName,
									node.NodeValue,
									node.Attributes,
									node.PartialXPath(),
									node.FullXPath())
								fmt.Printf("\tChildren: %+v \n", node.Children)*/
				quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					//					fmt.Printf("\t%s\n", err.Error())
					continue
				}
				//fmt.Printf("\tPosition top/bottom vertical element position: [%d px.. %d px]\n", int(quads[0][1]), int(quads[0][5]))
				if quads[0][1] > quads[0][5] {
					continue
				}
			}
			return nil
		}),
		// set viewport
		chromedp.EmulateViewport(1280, 720),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range headerNodes {
				/*				fmt.Printf("-------------------- %s ------------------ \n", node.NodeName)
								fmt.Printf("Node: %d %s %s %s %s %+v %s %s\n",
									node.NodeID,
									node.Name,
									node.Value,
									node.NodeName,
									node.NodeValue,
									node.Attributes,
									node.PartialXPath(),
									node.FullXPath())
								fmt.Printf("\tChildren: %+v \n", node.Children)*/
				quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					//					fmt.Printf("\t%s\n", err.Error())
					continue
				}
				//fmt.Printf("\tPosition top/bottom vertical element position: [%d px.. %d px]\n", int(quads[0][1]), int(quads[0][5]))
				if quads[0][1] > quads[0][5] {
					continue
				}
			}
			return nil
		}),
		// set viewport
		chromedp.EmulateViewport(1900, 1280),
		chromedp.Navigate(url),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
		}),
		chromedp.Nodes(":is(div, a, form, img, li, h1, h2, h3)", &headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range headerNodes {
				/*				fmt.Printf("-------------------- %s ------------------ \n", node.NodeName)
								fmt.Printf("Node: %d %s %s %s %s %+v %s %s\n",
									node.NodeID,
									node.Name,
									node.Value,
									node.NodeName,
									node.NodeValue,
									node.Attributes,
									node.PartialXPath(),
									node.FullXPath())
								fmt.Printf("\tChildren: %+v \n", node.Children)*/
				quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					//					fmt.Printf("\t%s\n", err.Error())
					continue
				}
				//fmt.Printf("\tPosition top/bottom vertical element position: [%d px.. %d px]\n", int(quads[0][1]), int(quads[0][5]))
				if quads[0][1] > quads[0][5] {
					continue
				}
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			//			var root *cdp.Node
			fmt.Println("---------------------")
			fmt.Println("GetDocument")
			fmt.Println("---------------------")
			/*			root, err = dom.GetDocument().Do(ctx)
						if err != nil {
							return err
						}*/
			// export html for debug
			res, err := dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			if err != nil {
				return err
			}
			h := sha1.New()
			h.Write([]byte(res))
			bs := h.Sum(nil)
			fmt.Printf("%x\n", bs)
			err = ioutil.WriteFile("layout.html", []byte(res), 0o644)
			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
