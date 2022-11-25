const melteSite = "melte"
const assets = [
  "/",
  "/out.js",
  "/thing",
  "/thing/out.js"
]

self.addEventListener("install", installEvent => {
  installEvent.waitUntil(
    caches.open(melteSite).then(cache => {
      cache.addAll(assets)
    })
  )
})

self.addEventListener("fetch", fetchEvent => {
    fetchEvent.respondWith(
      caches.match(fetchEvent.request).then(res => {
        return res || fetch(fetchEvent.request)
      })
    )
  })