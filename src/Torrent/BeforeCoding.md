我们在Bencode库中已经实现了如下功能:
+ 将bencode格式的文本转换为Bobject
+ 将Boject转换为bencode格式

但是我们想一个更轻松的结构，可以直接根据Bencode内容转换为一个Struct,能够直接
获得一系列信息

```
d
  8:announce
    41:http://bttracker.debian.org:6969/announce
  7:comment
    35:"Debian CD from cdimage.debian.org"
  13:creation date
    i1573903810e
  4:info
    d
      6:length
        i351272960e
      4:name
        31:debian-10.2.0-amd64-netinst.iso
      12:piece length
        i262144e
      6:pieces
        26800:�����PS�^�� (binary blob of the hashes of each piece)
    e
e
```

处理JSON格式的时候也有类似的步骤，我们称之为marshal和unmarshal
+ marshal 是将struct结构转换为bencode
+ unmarshal 是将bencode转换为struct结构