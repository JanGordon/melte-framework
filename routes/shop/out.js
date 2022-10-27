// components/app.js
function hello() {
  console.log("hello world");
}

// in.ts
{
  const SELF = document.querySelector("[melte-id='hi38']");
  hello();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
{
  const SELF = document.querySelector("[melte-id='w82']");
  hello();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
{
  const SELF = document.querySelector("[melte-id='hi90']");
  hello();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
{
  const SELF = document.querySelector("[melte-id='w134']");
  hello();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
