package dev

func router() string {
	return `(()=>{var r={};async function l(){var i=Array.from(document.querySelectorAll("a"));console.log("caching new links...",i);for(let u of i){console.log(u);var e=u;let a=new URL(e.href);await fetch(a).then(function(n){return n.text()}).then(function(n){r[a]=n}),u.addEventListener("click",async function(n){if(n.preventDefault(),r.hasOwnProperty(e)){let t=new URL(e.href);var d=r[t];let c=Range.prototype.createContextualFragment.bind(document.createRange());var h=document.implementation.createHTMLDocument();h.documentElement.innerHTML=d,console.log("Loading cached page",d),document.body.innerHTML=h.body.innerHTML,history.replaceState({},h.title,e.href),document.body.querySelectorAll("script").forEach(function(o){if(o.src.includes("out.js")){var m=new URL(o.src.slice(0,o.src.indexOf("out.js"))),f=document.createElement("script");o.parentNode.appendChild(f),o.remove(),f.src=m.pathname+"out.js?cachebuster="+new Date().getTime(),console.log(m.pathname+"out.js?cachebuster="+new Date().getTime())}}),l()}else setTimeout(()=>{console.log("Page hasn't been cached, loading...")},1e3),await fetch(a).then(function(t){return t.text()}).then(function(t){r[a]=t;var c=document.implementation.createHTMLDocument();c.documentElement.innerHTML=d,document.body.innerHTML=c.body.innerHTML,history.replaceState({},c.title,e.href),l()});var s="hello";console.log(s);var s="hello";function g(){}})}}Array(document.querySelectorAll("a")).forEach(function(i,e){});window.addEventListener("popstate",l);l();})();
	`
}

func hotReload() string {
	return `(()=>{try{o=new WebSocket("ws://"+location.host+"/hotReloadWS"),o.onopen=n=>{o.send("Succesfully connected")},o.onmessage=n=>{setTimeout(()=>{window.location.reload()},1e3),console.log("REloading"),o.close(),o.send("reloaded")},console.log("connected")}catch{}var o;})();`
}
