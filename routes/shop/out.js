// components/app.js
function hello() {
  console.log("hello world");
}

// in.ts
{
  const SELF = document.querySelector("[melte-id='hi0']");
  hello();
  console.log("rerunning js ", count);
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
