// in.ts
{
  const SELF = document.querySelector("[melte-id='counter0']");
  let count = 0;
  $:
    SELF.innerText = count;
}
