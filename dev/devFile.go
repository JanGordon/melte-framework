package dev

func router() string {
	return `(() => {
		// clientSideRouting/src.js
		var caches = {};
		async function cacheAllLinks() {
		  var links = Array.from(document.querySelectorAll("a"));
		  console.log("caching new links...", links);
		  for (let link of links) {
			console.log(link);
			var l = link;
			let url = new URL(l.href);
			if (window.location.origin == url.origin) {
			} else {
			  console.log("passing this link: not on same origin");
			  continue;
			}
			await fetch(url).then(function(response) {
			  return response.text();
			}).then(function(data) {
			  caches[url] = data;
			}).catch((error) => {
			  console.log("Failed to fetch link");
			});
			link.addEventListener("click", async function(e) {
			  e.preventDefault();
			  if (caches.hasOwnProperty(l)) {
				let url2 = new URL(l.href);
				var response = caches[url2];
				const parse = Range.prototype.createContextualFragment.bind(document.createRange());
				var doc = document.implementation.createHTMLDocument();
				doc.documentElement.innerHTML = response;
				console.log("Loading cached page", response);
				document.body.innerHTML = doc.body.innerHTML;
				history.replaceState({}, doc.title, l.href);
				document.body.querySelectorAll("script").forEach(function(script) {
				  if (script.src.includes("out.js")) {
					var newSrc = new URL(script.src.slice(0, script.src.indexOf("out.js")));
					var newScript = document.createElement("script");
					script.parentNode.appendChild(newScript);
					script.remove();
					newScript.src = newSrc.pathname + "out.js?cachebuster=" + new Date().getTime();
					console.log(newSrc.pathname + "out.js?cachebuster=" + new Date().getTime());
				  }
				});
				cacheAllLinks();
			  } else {
				setTimeout(() => {
				  console.log("Page hasn't been cached, loading...");
				}, 1e3);
				await fetch(url).then(function(response2) {
				  return response2.text();
				}).then(function(data) {
				  caches[url] = data;
				  var doc2 = document.implementation.createHTMLDocument();
				  doc2.documentElement.innerHTML = response;
				  document.body.innerHTML = doc2.body.innerHTML;
				  history.replaceState({}, doc2.title, l.href);
				  cacheAllLinks();
				}).catch(() => {
				  console.log("Failed to fetch link");
				});
			  }
			  {
				var hello = "hello";
			  }
			  {
				console.log(hello);
			  }
			  var hello = "hello";
			  function onChange() {
			  }
			});
		  }
		}
		Array(document.querySelectorAll("a")).forEach(function(link, index) {
		});
		window.addEventListener("popstate", cacheAllLinks);
		cacheAllLinks();
	  })();
	  `
}

func hotReload() string {
	return `(()=>{try{o=new WebSocket("ws://"+location.host+"/hotReloadWS"),o.onopen=n=>{o.send("Succesfully connected")},o.onmessage=n=>{setTimeout(()=>{window.location.reload()},1e3),console.log("REloading"),o.close(),o.send("reloaded")},console.log("connected")}catch{}var o;})();`
}
