# sov2ex
一个第三方的 V2EX 站内搜索

# 使用
直接访问 https://www.sov2ex.com

# API

### 搜索
URL: https://www.sov2ex.com/api/search

Method: GET

Params:

| 参数名称 | 类型 | 必须 | 描述 |
| --- | --- | --- | --- |
| q | string | true | 查询关键词 |

```
https://www.sov2ex.com/api/search?q=大规模集群环境
```

Reponse:

| 参数名称 | 类型 | 必须 | 描述 |
| --- | --- | --- | --- |
| took | int | true | 搜索过程耗时(ms) |
| timed_out | bool | true | 是否超时 |
| total | int | true | 命中主题总数 |
| hits | array | true | 主题列表 |
| &nbsp;&nbsp;_source | object | true | 主题信息 |
| &nbsp;&nbsp;&nbsp;&nbsp;node | int | true | 结点 id |
| &nbsp;&nbsp;&nbsp;&nbsp;replies | int | true | 回复数量 |
| &nbsp;&nbsp;&nbsp;&nbsp;created | string | true | 创建时间(UTC) |
| &nbsp;&nbsp;&nbsp;&nbsp;member | string | true | 主题作者 |
| &nbsp;&nbsp;&nbsp;&nbsp;id | int | true | 主题 id |
| &nbsp;&nbsp;&nbsp;&nbsp;title | string | true | 主题标题 |
| &nbsp;&nbsp;&nbsp;&nbsp;content | string | true | 主题内容 |
| &nbsp;&nbsp;highlight | object | false | 高亮 |
| &nbsp;&nbsp;&nbsp;&nbsp;title | array | false | 标题高亮（最多 1 条） |
| &nbsp;&nbsp;&nbsp;&nbsp;content | array | false | 主题内容高亮（最多 1 条） |
| &nbsp;&nbsp;&nbsp;&nbsp;postscript_list.content | array | false | 附言高亮（最多 1 条） |

```
{
    "took": 34,
    "timed_out": false,
    "total": 53591,
    "hits": [
        {
            "_index": "topic_v1",
            "_type": "topic",
            "_id": "303776",
            "_score": 91.11005,
            "_source": {
                "node": 11,
                "replies": 13,
                "created": "2016-09-04T01:37:41",
                "member": "jasonailu",
                "id": 303776,
                "title": "怎样在公共集群上隔离出自己的空间？",
                "content": "公共集群就是很多人具有集群的 root 用户密码，\r\n\r\n我想隔离出自己的空间，并且防止依赖的基础库被删除，\r\n\r\n另外，请教有什么其他在公共集群的管理使用方法？\r\n\r\nPs. Redhat."
            },
            "highlight": {
                "title": [
                    "怎样在公共<em>集群</em>上隔离出自己的空间？"
                ],
                "postscript_list.content": [
                    "Hadoop 在 docker 的生产<em>环境</em>和真机器上的性能什么的有区别吗？\n\n我想在<em>集群</em>搭建 docker + hadoop <em>集群</em>不知道可行不？与真机相比性能如何？这样备份 image 也好点。"
                ],
                "content": [
                    "公共<em>集群</em>就是很多人具有<em>集群</em>的 root 用户密码，\r\n\r\n我想隔离出自己的空间，并且防止依赖的基础库被删除，\r\n\r\n另外，请教有什么其他在公共<em>集群</em>的管理使用方法？\r\n\r\nPs. Redhat."
                ]
            }
        }
    ]
}
```
