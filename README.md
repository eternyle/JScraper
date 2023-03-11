# JScraper

JScraper is a Go tool that collects links to JavaScript files from stdin.


### Installation

```
go get github.com/eternyle/JScraper
```

### Usage

```
cat subdomain.txt | JScraper
```

#### Output

```
https://cdn.cookielaw.org/scripttemplates/otSDKStub.js
https://cdn.azizhakim.com/static/js/vendor.883ed9a0.js
https://cdn.azizhakim.com/static/js/retailer.0431142c.en.js
https://cdn.cookielaw.org/consent/12e743c7-ebdb-4c47-b445-d37bcd22ba62/OtAutoBlock.js
https://cdn.cookielaw.org/scripttemplates/otSDKStub.js
https://cdn.azizhakim.com/static/js/vendor.883ed9a0.js
https://cdn.azizhakim.com/static/js/retailer.0431142c.en.js
https://cdn.cookielaw.org/consent/12e743c7-ebdb-4c47-b445-d37bcd22ba62/OtAutoBlock.js
```

* subdomain.txt

```
https://abc.azizhakim.com
https://xyz.azizhakim.com
https://hello.azizhakim.com
```

JScraper + Linkfinder

```
cat subdomain.txt | JScraper >> jsfile.txt
```

```
while read line; do python3 ~/tools/linkfinder.py -i $line -o cli ; done < jsfile.txt
```
