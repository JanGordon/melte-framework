var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __commonJS = (cb, mod) => function __require() {
  return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));

// components/app.js
var require_app = __commonJS({
  "components/app.js"() {
    var result = (() => {
      var hello2 = "hellol";
      return hello2;
    })();
    console.log(result);
  }
});

// in.ts
var import_app = __toESM(require_app());
{
  const SELF = document.querySelector("[melte-id='hi90']");
  (0, import_app.default)();
  SELF.querySelector("button").addEventListener("click", function() {
    console.log(SELF.querySelector("button").innerText);
    count += 1;
    SELF.querySelector("button").innerText = count.toString();
  });
  SELF.querySelector("button").innerText = count.toString();
}
