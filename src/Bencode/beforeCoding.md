# Bencode库

我们将四种类型称为Bobject.

**String**

bencode编码中,string的表示如下

`5:hello`,需要注意的是`5`并不是真正的数值5，而是ascii码值`'5'`

**int**

int的表示如下:

`i2200e`,也就是以`i`开头，`e`结尾，中间是ascii形式的`0-9`和`-`

**List**

List是一个复合类型，前面的int和string是符合类型,我们表示为`l[][]e`,
其中`[]`可以是任意的类型(string,int,list,dict)。

比如`l[1:a][li2200e3:abc]e`,这就表示为['a',[2200,'abc']].可以使用递归下降
的方法解析。

**Dict**

dict与list比较类似，形式如:`d[key][Boject]e`

比如`d[1:a][li2200e3:abc]e`，表示`'a':[2200,'abc']`

----------------
因此我们需要完成的任务:

io.Reader (读取bytes --> 解析为) Bobject
(-->序列化) io.Writer