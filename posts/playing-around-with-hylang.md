<!--id: 8-->
<!--title: Playing Around With Hylang -->
<!--author: Brian Jones-->
<!--visible: true-->

Lately I've been writing python, playing with different lisps (mostly common lisp), and doing alot of work on AWS. So I thought I'd write a quick post on a cool python based lisp I've played around with for about 45 minutes.

### [Hylang](http://docs.hylang.org/en/stable/)

Here's how I've set my environment up for some quick Hylang development (the official docs are also pretty good. I got up and running in 15 min, or less.)

```
$ pyenv local 3.72 # this is just the latest I had installed already
```
```
$ pyenv virtualenv tmp-hy
```
```
$ pyenv activate 3.7.2/envs/tmp-hy
```
```
$ pip install git+https://github.com/hylang/hy.git
```

Those 3 incantations will give you a nice sandboxed installation of Hylang. You can run the interpreter by typing just `hy` or you can give it a file to munch on by saying `hy some-code.hy`.


Since I've been using DynamoDB a lot lately I thought I'd try out some boto3. Turns out it's pretty easy.

if you're following along make sure you've already `$ pip install git+https://github.com/hylang/hy.git` and also make sure to run:
`$ pip install boto3`  

```lisp
(import boto3)
(setv dynamo-resource (boto3.resource "dynamodb"))
(setv blog-post-table (dynamo-resource.Table "blog-posts"))
(setv key-expr {
  "id" "7"
  "md5" "ef7582d3ccac18418063bd19715614af" })
(print key-expr)
(setv item (get  (blog-post-table.get-item :Key key-expr ) "Item"))
(print item)
```

if you enter: `hy lispy-boto.hy` it will print out the blog post.

Even better, we could define a function like:

```lisp
(setv dynamo-resource (boto3.resource "dynamodb"))

(defn lispy-get-item [hash_key_name hash_key_val range_key_name range_key_val table_name]
  (setv dynamo-table (dynamo-resource.Table table_name))
  (setv key-expr {
    hash_key_name hash_key_val
    range_key_name range_key_val })
  (get (dynamo-table.get-item :Key key-expr) "Item"))
```

and we can import that function into plain 'ol python like:

```python
import hy
lispy_boto = __import__("lispy-boto") # hack because I named the file with a `-`

item = lispy_boto.lispy_get_item(
  "id",
  "7",
  "md5",
  "ef7582d3ccac18418063bd19715614af",
  "blog-posts") # You could use kwargs to make this more pythonic

print(item)
```

This is just scratching the surface, but I was surprised how easy this was and how little fiddling it took.

Hylang can be found [on github](https://github.com/hylang/hy) and the [docs](http://docs.hylang.org/en/stable/quickstart.html) aren't half bad. Really cool project