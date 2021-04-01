# ChromeDP test

This repository tries to traverse the DOM of a given URL in order to retrieve vertical position of element with a given screen width.

# Execute

Install Go 1.16 (https://golang.org/dl/)

Execute :
```sh
$ go run main.go -u "https://www.cdiscount.com/search/10/sac+michael+kors.html"
```

# references
* [Chrome DevTools: DOM/type-node](https://chromedevtools.github.io/devtools-protocol/tot/DOM/#type-Node)
* [W3C: Nodetypes](https://dom.spec.whatwg.org/#dom-node-nodetype)
