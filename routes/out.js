// in.ts
{
  const SELF = document.querySelector("[melte-id='hi90']");
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
