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
)

func traverse(ctx context.Context, node *cdp.Node, name string) error {
	fmt.Printf("[%d] %d /%s/%s %s %d\n", node.NodeID, node.NodeType, name, node.NodeName, node.NodeValue, node.ChildNodeCount)
	quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("quads: %v\n", quads)
	for _, n := range node.Children {
		traverse(ctx, n, fmt.Sprintf("/%s/%s", name, node.NodeName))
	}
	return nil
}

func main() {
	start := time.Now()
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run
	err := chromedp.Run(ctx,
		// set viewport
		chromedp.EmulateViewport(780, 651),
		chromedp.Navigate(`https://www.whatsmyua.info/?a`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			var root *cdp.Node
			fmt.Println("---------------------")
			fmt.Println("GetDocument")
			fmt.Println("---------------------")
			root, err = dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			err = traverse(ctx, root, "")
			if err != nil {
				return err
			}
			// https://chromedevtools.github.io/devtools-protocol/tot/DOM/#type-Node
			// https://dom.spec.whatwg.org/#dom-node-nodetype
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	/*
		var ids []cdp.NodeID
			err = chromedp.Run(ctx,
				// select by JS path
				chromedp.NodeIDs(`document.querySelector("body > div > div.top.block")`, &ids, chromedp.ByJSPath),
				chromedp.ActionFunc(func(ctx context.Context) error {
					id := ids[0]
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
					return nil
				}),
				chromedp.NodeIDs(`document.querySelector("body > div > div.api.block")`, &ids, chromedp.ByJSPath),
				chromedp.ActionFunc(func(ctx context.Context) error {

					id := ids[0]
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
					return nil
				}),
			)
			if err != nil {
				log.Fatal(err)
			}
	*/
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
