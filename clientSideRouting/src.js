var caches = {};



async function cacheAllLinks() {
    var links = Array(document.querySelectorAll("a"))
    for (let link of links) {
        console.log(links, link)
        var l = link[0]
        let url = new URL(l.href)
        await fetch(url)
        .then(function(response) {
            return response.text()
        }).then(function(data) {
            caches[url] = data
        })
    }
}

// add vent lister for all links and prevent default even if lniks are still loading
Array(document.querySelectorAll("a")).forEach(function(link, index) {
    
    link[index].addEventListener("click", async function(e) {
        e.preventDefault();
        var l = link[0]
        if (caches.hasOwnProperty(l)) {
            let url = new URL(l.href)
            var response = caches[url]
            const parse = Range.prototype.createContextualFragment.bind(document.createRange());
            var doc = document.implementation.createHTMLDocument();
            doc.body.innerHTML = "Hello my nane is jan"
            console.log("Loading cached page")
            document.body.innerHTML = doc.body.innerHTML
        } else {
            console.log("Page hasn't been cached, loading...")
            await fetch(url)
            .then(function(response) {
                return response.text()
            }).then(function(data) {
                caches[url] = data
                var doc = document.implementation.createHTMLDocument();
                doc.querySelector("html").innerHTML = response
        
            })
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

        history.replaceState( {} , doc.title, l.href );
        document.querySelector("body").innerHTML = doc.body.innerHTML

    })
})
//listen for chnages in dom and see if links have been modified or added
cacheAllLinks()