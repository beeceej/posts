<!--id: 2-->
<!--title: Dynamo DB and Elastic Search! -->
<!--author: Brian Jones-->
<!--visible: false-->

# [DynamoDB](https://aws.amazon.com/dynamodb/) Is Awesome!

it has a free tier and it's incredibly easy to stand up a table withour worrying about server instances. It does have warts though, things like [Not supporting empty strings](https://forums.aws.amazon.com/thread.jspa?threadID=90137) and iffy queryability (which I'll show a strategy for coping). DynamoDB is able to handle heavy write workloads like a champ and queries on well known hash/range keys are efficient. But, you'll quickly run into issues when you need to issue complex queries on multiple data facets. The Dynamo based solution is to create multiple indexes (_equivalent to copying the whole table_). Multiple indexes become expensive, and don't scale well. If an application needs to scan through a Dynamo Table, it will be _**slow**_ and _**expensive**_. A much better pattern would be to create a Materialized View of your Dynamo Data in a search index like Solr, or in this case [ElasticSearch](https://www.elastic.co/products/elasticsearch). Luckily Dynamo is awesome and provides a feature called [Dynamo Streams](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.html). These make it possible to stream data in real time. In practice, Dynamo Streams, lambda functions and Elasticsearch have proven to be a great pattern for indexing Dynamo Data, while avoiding some its pain points. This pattern has proven useful and repeatable enough to abstract this concept out into a Serverless project. Let's dig into how it works.

## Elasticsearch Primer
