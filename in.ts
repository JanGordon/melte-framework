    import hello from "Component/app.js"

{
// script for outout-counter0.js.js
 const SELF = document.querySelector("[melte-id='counter0']")
    let count = 0;
    $: SELF.innerText = count

}