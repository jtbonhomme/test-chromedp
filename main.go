// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	start := time.Now()
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run
	//var b1 []byte
	var ids []cdp.NodeID
	err := chromedp.Run(ctx,
		// set viewport
		chromedp.EmulateViewport(780, 651),
		chromedp.Navigate(`https://www.whatsmyua.info/?a`),
		// scroll to a xpath
		//chromedp.ScrollIntoView(`/html/body/div/div[5]`), // API div
		//chromedp.CaptureScreenshot(&b1),
		// select by JS path
	)
	if err != nil {
		log.Fatal(err)
	}

	// dom: extract full html
	err = chromedp.Run(ctx,
		chromedp.NodeIDs(`document.querySelector("body > div > div.top.block")`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			/*fmt.Println("---------------------")
			fmt.Println("GetDocument")
			fmt.Println("---------------------")
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
			spew.Dump(node)
			*/

			id := ids[0]
			//id := node.NodeID
			fmt.Println("---------------------")
			fmt.Println("GetBoxModel")
			fmt.Println("---------------------")

			b, err := dom.GetBoxModel().WithNodeID(id).Do(ctx)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
			//An array of quad vertices, x immediately followed by y for each point, points clock-wise.
			spew.Dump(b.Content)
			/*
				fmt.Println("---------------------")
				fmt.Println("DescribeNode")
				fmt.Println("---------------------")

				desc, err := dom.DescribeNode().WithNodeID(id).Do(ctx)
				if err != nil {
					fmt.Println("ERROR ", err)
					return err
				}
				spew.Dump(desc)
			*/
			// Get HTML from the page
			//fmt.Println("---------------------")
			//fmt.Println("GetOuterHTML")
			//fmt.Println("---------------------")
			//res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			//fmt.Println(res)
			return err
		}),
		chromedp.NodeIDs(`document.querySelector("body > div > div.api.block")`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			/*fmt.Println("---------------------")
			fmt.Println("GetDocument")
			fmt.Println("---------------------")
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
			spew.Dump(node)
			*/

			id := ids[0]
			//id := node.NodeID
			fmt.Println("---------------------")
			fmt.Println("GetBoxModel")
			fmt.Println("---------------------")

			b, err := dom.GetBoxModel().WithNodeID(id).Do(ctx)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
			//An array of quad vertices, x immediately followed by y for each point, points clock-wise.
			spew.Dump(b.Content)
			/*
				fmt.Println("---------------------")
				fmt.Println("DescribeNode")
				fmt.Println("---------------------")

				desc, err := dom.DescribeNode().WithNodeID(id).Do(ctx)
				if err != nil {
					fmt.Println("ERROR ", err)
					return err
				}
				spew.Dump(desc)
			*/
			// Get HTML from the page
			//fmt.Println("---------------------")
			//fmt.Println("GetOuterHTML")
			//fmt.Println("---------------------")
			//res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			//fmt.Println(res)
			return err
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	/*
		// page
		err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("---------------------")
			fmt.Println("GetLayoutMetrics")
			fmt.Println("---------------------")
			v1, v2, v3, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
			spew.Dump(v1)
			spew.Dump(v2)
			spew.Dump(v3)

			return err
		}))
		if err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile("screenshot.png", b1, 0o644); err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote screenshot1.png")
	*/
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
