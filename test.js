var obj = {"hello"}
  
  var proxied = new Proxy(obj, {
    get: function(target, prop) {
      console.log({ type: 'get', target, prop });
      return Reflect.get(target, prop);
    },
    set: function(target, prop, value) {
      console.log({ type: 'set', target, prop, value });
      return Reflect.set(target, prop, value);
    }
  });
  
  proxied.bar = 2;
  // {type: 'set', target: <obj>, prop: 'bar', value: 2}
  
  proxied.foo;
  // {type: 'get', target: <obj>, prop: 'bar'}