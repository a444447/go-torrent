# Building a BitTorrent client from ground


# Introduction

BitTorrent 是一个在互联网上下载和分享文件的协议。传统的 client/server 关系中，下载者们都连接到一个中央服务器（例如：在 Netflix 上看电影，或者加载你正在阅读的网页），而在 BitTorrent 网络中，下载者被称为 peers，他们互相之间下载文件的某一个片段——这就是被称为 peer-to-peer（P2P） 的协议。




中心模式:
缺点:如果离Server很远，则转播时延很高

P2P模式:

所有参与下载的下载方都是peer。加载东西的时候可以直接通过已经有相应资源的peer中
下载该资源(人人为我,我为人人)。

争议:存在很多的版权问题

为了完成P2P模式我们需要考虑的问题:

1. 如何找到peers
    
> **找到一个中心化的站点tracker,它能提供拥有资源的peers的名单**\
> 例如以前国内流行的PT站:它作为tracker提供peers名单，同时你下载资源也会成为其中的peer。




2. 如何与peers协作完成下载

BitTorrent主要原理是需要把提供下载的文件虚拟分成大小相等的块，块大小必须为2k的整数次方（由于是虚拟分块，硬盘上并不产生各个块文件），并把每个块的索引信息和Hash验证码写入种子文件中

现在我们总结一些构造一个BitTorrent客服端的步骤:
+ step1 Bencode库 -> 序列化与反序列化
+ step2 torrent解析文件 -> 获得tracker 和 info
+ step3 tracker模块  -> peers信息
+ step4 download模块 -> pieces与校验
+ step5 assember -> 把pieces拼装为file


# Finding Peers



## .torrent file

.torrent 文件描述了可下载文件的内容和连接到特定 tracker 的信息，而这就是开始下载一个 torrent 所需要的全部。Debian 的 .torrent 文件长这样：

```
d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:"Debian CD from cdimage.debian.org"13:creation datei1573903810e9:httpseedsl145:https://cdimage.debian.org/cdimage/release/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso145:https://cdimage.debian.org/cdimage/archive/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.isoe4:infod6:lengthi351272960e4:name31:debian-10.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces26800:�����PS�^�� (binary blob of the hashes of each piece)ee
```



这一堆乱码被用一种称为 Bencode（发音为 bee-encode）的方法进行了编码，我们需要想办法解码它。
Bencode包含四种类型:string,int,slice,dictionary

Bencode虽然不如JSON一样可读，但是有一个优势就是很容易从流中转换，
>适合网络传输(一个字符一个字符的读取)\
> 比如给出`8:announce`我们仅根据第一个字符就可以知道长度和类型
> 因此我们就可以知道接下来走什么逻辑

因为string前面的编码代表了他们的长度:

例如,数字`7`表示为`i7e`,`4:spam`表示字符串spam,`l4:spami7e` 表示 `['spam',7]`，`d4:spami7ee` 表示 `{spam:7}`。
> 注意:`list`,`int`,`dict`都是以l,i,d开头，e结尾

