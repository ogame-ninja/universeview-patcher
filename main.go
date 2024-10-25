package main

import (
	ep "github.com/ogame-ninja/extension-patcher"
)

func main() {
	const (
		webstoreURL         = "https://chromewebstore.google.com/detail/universeview-extension/ipmfkhoihjbbohnfecpmhekhippaplnh"
		universeview_sha256 = "84d99167220fc7563b4ebb925f80b1e65d9420ea4a2d16ccccebafbe7d2da259"
	)

	files := []ep.FileAndProcessors{
		ep.NewFile("manifest.json", processManifest),
		ep.NewFile("background.js", processBackgroundJs),
		ep.NewFile("chrome/content/scripts/universeview.js", processUniverseviewJs),
	}

	ep.MustNew(ep.Params{
		ExpectedSha256: universeview_sha256,
		WebstoreURL:    webstoreURL,
		Files:          files,
	}).Start()
}

var replN = ep.MustReplaceN

func processManifest(by []byte) []byte {
	by = replN(by, `"name": "UniverseView Extension",`, `"name": "UniverseView Ninja Extension",`, 1)
	by = replN(by, `"*://*.ogame.gameforge.com/game/index.php*"`, `{old}, "*://*/bots/*/browser/html/*"`, 2)
	by = replN(by, `*://*.ogame.gameforge.com/*`, `<all_urls>`, 1)
	return by
}

func processBackgroundJs(by []byte) []byte {
	by = replN(by, `var a=/s(\d+)-(\w+)\.ogame\.gameforge\.com/.exec(t.url),`, `const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(t.url)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(t.url)[2];`, 1)
	by = replN(by, `a[1]+"-"+a[2];`, `universeNum + "-" + lang;`, 1)
	by = replN(by, `"https://"+o.ogameURL.replace("s{UID}","*")+"/*"`, `"*://*/*"`, 1)
	by = replN(by, `initialize(t){`, `initialize(t, u){const myUrl = new URL(u);const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(u)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(u)[2];`, 1)
	by = replN(by, `universe.initialize(s)`, `universe.initialize(s, t.url)`, 1)
	by = replN(by, `this.origin[t]="https://"+o.ogameURL.replace("{UID}",t)`, `this.origin[t]=myUrl.origin+"/api/s" + universeNum + "/" + lang + "/"`, 1)
	by = replN(by, `+"/api/"+`, `+`, 2)
	return by
}

func processUniverseviewJs(by []byte) []byte {
	by = replN(by, `const e=/s(\d+)-(\w+)\.ogame\.gameforge\.com/.exec(a.location.host);`, `const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[2];`, 1)
	by = replN(by, `u.OGAMELEGACY=/^\d/.exec(d.version)[0]<7`, `u.OGAMELEGACY = false`, 1)
	by = replN(by, `d.universeNumber=e[1],d.community=e[2]`, `d.universeNumber = universeNum, d.community = lang`, 1)
	by = replN(by, `/game/index.php`, ``, 4)
	by = replN(by, `a.location=a.location.origin+`, `a.location =`, 2)
	by = replN(by, `href="index.php?`, `href = "?`, 2)
	by = replN(by, `checkXML:function(e,t,i){`, `{old} const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[1];
const lang=/browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[2];`, 1)
	by = replN(by, `fetchXML:function(e,t,i,n){`, `{old} const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(a.location.href)[2];`, 1)
	by = replN(by, `url:(i||a.location.origin)+"/api/"+t,`, `
url: (i || a.location.origin) + "/api/s" + universeNum + "/" + lang  + "/" + t,`, 1)
	by = replN(by, `url:(n||a.location.origin)+"/api/"+t,`, `
url: (n || a.location.origin) + "/api/s" + universeNum + "/" + lang  + "/" + t,`, 1)
	return by
}
