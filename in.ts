    import hello from "Component/app.js"

{
// script for outout-hi0.js.js
 const SELF = document.querySelector("[melte-id='hi0']")
    hello()
    console.log("rerunning js ", count)
    
    SELF.querySelector("button").addEventListener("click", function() {
        console.log(SELF.querySelector("button").innerText)
        count += 1
        SELF.querySelector("button").innerText = (count).toString()
    })
    SELF.querySelector("button").innerText = (count).toString()
    

}