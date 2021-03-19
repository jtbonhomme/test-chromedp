// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func traverse(ctx context.Context, node *cdp.Node, name string) error {
	if node == nil {
		return nil
	}
	fmt.Printf("[%d] %d /%s/%s %s %d\n", node.NodeID, node.NodeType, name, node.NodeName, node.NodeValue, node.ChildNodeCount)
	quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("quads: %v\n", quads)
	nodeParams, err := dom.DescribeNode().WithNodeID(node.NodeID).WithDepth(1).WithPierce(true).Do(ctx)
	err = dom.RequestChildNodes(nodeParams.NodeID).WithDepth(-1).Do(ctx)
	if err != nil {
		return err
	}
	for _, n := range nodeParams.Children {
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
	var ids []cdp.NodeID
	var headerNodes []*cdp.Node
	err := chromedp.Run(ctx,
		// set viewport
		chromedp.EmulateViewport(780, 651),
		chromedp.Navigate(`https://www.whatsmyua.info/?a`),
		chromedp.NodeIDs("body", &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(ids[0]).WithDepth(-1).Do(c)
		}),
		chromedp.Nodes(":is(div, a, form, img)", &headerNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range headerNodes {
				quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("-------------------- %s ------------------ \n", node.NodeName)
				fmt.Printf("Node: %d %s %s %s %s %+v %s\n",
					node.NodeID,
					node.Name,
					node.Value,
					node.NodeName,
					node.NodeValue,
					node.Attributes,
					node.FullXPath())
				fmt.Printf("\tChildren: %+v \n", node.Children)
				fmt.Printf("\tPosition top/bottom vertical element position: [%d px.. %d px]\n", int(quads[0][1]), int(quads[0][5]))
			}
			return nil
		}),
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
			// export html for debug
			res, err := dom.GetOuterHTML().WithNodeID(root.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile("layout.html", []byte(res), 0o644)
			if err != nil {
				return err
			}

			/*			err = traverse(ctx, root, "")
						if err != nil {
							return err
						}*/
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
