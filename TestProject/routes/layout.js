export default function layout () {
    return "<h1>Layout</h1><div><slot></slot></div>"
}

// THis layout will be used by every page until a new layout file is specified.
// The new layout must have a differnt name
// layout-itemlist.js
// A file can specify which layout it needs to use.
// A layout can even contain another layout by importing it see layout.melte
