var caches = {};



async function cacheAllLinks() {
    var links = Array.from(document.querySelectorAll("a"))
    console.log("caching new links...", links)
    for (let link of links) {
        console.log(link)
        var l = link
        let url = new URL(l.href)
        if( window.location.origin == url.origin) {
        } else {
            console.log("passing this link: not on same origin")
            continue
        }
        await fetch(url)
        .then(function(response) {
            return response.text()
        }).then(function(data) {
            caches[url] = data
        }).catch((error)=>{console.log("Failed to fetch link")})


        link.addEventListener("click", async function(e) {
            e.preventDefault();
            if (caches.hasOwnProperty(l)) {
                let url = new URL(l.href)
                var response = caches[url]
                const parse = Range.prototype.createContextualFragment.bind(document.createRange());
                var doc = document.implementation.createHTMLDocument();
                doc.documentElement.innerHTML = response
                console.log("Loading cached page", response)
                document.body.innerHTML = doc.body.innerHTML
                console.log(doc.head.innerHTML)
                document.head.querySelectorAll("*:not(script)").forEach(function(child) {
                    child.remove()
                })
                doc.head.querySelectorAll("*:not(script").forEach(function (child) {
                    document.head.appendChild(child)
                })
                // document.head.innerHTML = doc.head.innerHTML
                history.replaceState( {} , doc.title, l.href );
                var ssrScripts = ""
                console.log()
                doc.querySelectorAll("script[ssr='']").forEach(function (child) {
                    if (child.innerText.charAt(child.innerText.length-1) == ";") {
                        ssrScripts+=child.innerText
                    } else {
                        ssrScripts+=child.innerText+"\n"
                    }
                })
                
                
                // refill <melte-reload>
                doc.body.querySelectorAll("melte-reload").forEach(function(child){
                    if (child.getAttribute("js") != "") {
                        ssrScripts+= ""// add proxy and modifier
                        // child.innerText = Function("'use strict';" + child.getAttribute("js"))()
                        console.log("Preloaded responsive html")
                    } else {
                        console.log("melte-reload js field empty")
                    }
                })
                Function(ssrScripts + "function handler(){}")()
                
                document.body.querySelectorAll("script").forEach(function (script) {
                    if (script.src.includes("out.js")) {
                        var newSrc = new URL(script.src.slice(0, script.src.indexOf("out.js")))
                        var newScript = document.createElement("script")
                        script.parentNode.appendChild(newScript)
                        script.remove()
                        newScript.src = newSrc.pathname + "out.js?cachebuster="+ new Date().getTime()
                        console.log(newSrc.pathname + "out.js?cachebuster="+ new Date().getTime())
                    }
                })
                cacheAllLinks()
            } else {
                setTimeout(()=>{console.log("Page hasn't been cached, loading...")}, 1000)
                await fetch(url)
                .then(function(response) {
                    return response.text()
                }).then(function(data) {
                    caches[url] = data
                    var doc = document.implementation.createHTMLDocument();
                    doc.documentElement.innerHTML = response
                    document.body.innerHTML = doc.body.innerHTML
                    history.replaceState( {} , doc.title, l.href );
                    cacheAllLinks()
                })
                .catch(()=>{console.log("Failed to fetch link")})
            }
    
            {var hello = "hello"}
            {console.log(hello)}
            // should html frags be replaced server side with every request and updated wiht js or 
            // js do everything
            //both
            //make state kept variables opt in to sevrer hydration
            // all compoennt scripts should be removed for new page
            // scripts can be rerun on every state using the popstate event
            //defien custom variable like this
            
            //@melte-custom: var abs, global, server
            var hello = "hello"
            
            //define on var change:
    
            //@melte-custom: function change, hello
            function onChange () {
    
            }
    
            
        })
    }
}

// add vent lister for all links and prevent default even if lniks are still loading
Array(document.querySelectorAll("a")).forEach(function(link, index) {
    
    
})
function makeEvalContext(declarations) {
    eval(declarations);
    return function (str) { eval(str); }
}
//listen for chnages in dom and see if links have been modified or added
window.addEventListener("popstate", cacheAllLinks)
cacheAllLinks()