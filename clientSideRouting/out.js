// in.ts
var caches = {};
async function cacheAllLinks() {
  var links = Array(document.querySelector("a"));
  for (let link of links) {
    console.log("Caching: ", link.href);
    caches[link.href] = await fetch(link.href);
  }
}
cacheAllLinks();
