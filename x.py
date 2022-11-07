t = """

<h1>Hello this is the shop page</h1>
<a href="/">Home page</a>
<hi></hi>
<w></w>
{{for (let i of [1,2,3,4])
    <h1></h1><p></p>
    <w></w>
}}
<!-- If i have counter then it should be able to use globally state kept varibales between routes -->
<!-- <counter val={count}></counter> -->
"""

while True:
    f = int(input("Char : "))
    print(t[f-2],t[f-1],t[f],t[f+1], t[f+2])