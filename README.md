# melte-framework

This is my first project made with go so excuse the very badly organised code and commented out code

I made this framework because I think that modern frameworks rely too heavily on JS and should mostly serve static html.

If you can read my code then pull requests are welcome but don't expect me to very quick as I have school. 


# Using it (why?)
To run start the dev server you firstly have to install it with ./install.sh (this works on my chromebook but I can't guarantte that it will work on anything else) then run 
    melte dev <port number>
There will be a lot of junk printed out but you should be able to connect to localhost:<port number>

# Features
- For loops:
  {{for (let i of [1,2,3,4])
    <h1>${i}</h1>
  }}
- Client side routing
- A feeble attempt at hot reload (it is very temperemental)
- Preserve state across routes (hopefully soon with indexedDB as well):
  Put this before variable declarations:
    //=keep state: js
    var count = 10;
- components:
  place .melte file in rootofproject/components and write basically svelte
- File based router
- I think its quite quick but the need for v8 slows it down a lot

It should have very good client side performance and got 100 lighthouse score in my tests

# What is still needs:
- A lot
- layouts
- desparately a .gitignore

You can (will be) see an example of most of the features [here](https://www.github.com/JanGordon/melte-demo)
    
Don't know I called it melte proabaly because you're a [melte](https://www.urbandictionary.com/define.php?term=melt)
