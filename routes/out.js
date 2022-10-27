// components/app.js
function hello() {
  console.log("hello world");
}

// in.ts
{
  const SELF = document.querySelector("[melte-id='hi54']");
  hello();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
