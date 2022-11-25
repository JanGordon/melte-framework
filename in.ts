    import hello from "Component/app.js"

{
// script for out-hi90.js
 const SELF = document.querySelector("[melte-id='hi90']")
    hello()

    SELF.querySelector("button").addEventListener("click", function() {
        console.log(SELF.querySelector("button").innerText)
        count += 1
        SELF.querySelector("button").innerText = (count).toString()
    })
    SELF.querySelector("button").innerText = (count).toString()
    //Using js to update very page load is too slow and causes an update which is very ugly

}