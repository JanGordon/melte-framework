// in.ts
{
  const SELF = document.querySelector("[melte-id='hi124']");
  let count = counthi124;
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("#s").innerText = count.toString();
  });
}
