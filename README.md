# ChromeDP test

This repository tries to traverse the DOM of a given URL in order to retrieve vertical position of element with a given screen width.

# Execute

Install Go 1.16 (https://golang.org/dl/)

Run :
```sh
$ go run main.go -u "https://www.cdiscount.com/search/10/sac+michael+kors.html"
```

# Usage

```
  -max string
  -M string
        the maximum  screen width (default "1900")
  -height string
  -H string
        the default screen height (default "750")
  -inc string
  -i string
        the screen width increment (default "100")
  -min string
  -m string
        the minimum screen width (default "800")
  -url string
  -u string
        the URL to be scrapped (default "https://www.whatsmyua.info/?a")
```

# references
* [Chrome DevTools: DOM/type-node](https://chromedevtools.github.io/devtools-protocol/tot/DOM/#type-Node)
* [W3C: Nodetypes](https://dom.spec.whatwg.org/#dom-node-nodetype)
