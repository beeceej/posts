<!--id: 6-->
<!--title: Dynamo DB Data Migrations -->
<!--author: Brian Jones-->
<!--postedAt: October 22nd, 2018-->
<!--updatedAt: October 22nd, 2018-->
<!--visible: true-->

Over the past month or so I've had the pleasure of working on a huge data migration on a live, customer facing application. We had chosen a primary key that a large part of our infrastructure depended upon; unfortunately for us, that key was no longer primary. Quick deadlines and changing requirements put us into a position where we needed to overhaul the identifier across all services. As dead lines loomed, we knew we had to do it quickly, safely, and accurately. To make this change as seamless as possible we needed to perform the migration:

1. Quickly (abbreviated schedule... out of our control).
2. Safely (Can't lose customer data)
3. Accurately (In transition to the new data set, we needed to ensure the integrity of our existing data either stand the same, or improved)

We went to the drawing board and came up with a plan. The plan accounted for those 3 things, quick, safe, accurate, I'll go into *how* we did that shortly, but first an overview of the architecture we were dealing with.

## Architecture At a Glance
* Services Deployed in:
  * AWS ECS
  * AWS Lambda 
  * Elastic Beanstalk
* Services Utilizing:
  * AWS DynamoDB
  * Elastic Search
  * AWS SNS
  * AWS Kinesis
  * AWS SQS
  * AWS Step Functions
  * AWS API Gateway
* Various front applications
   

Our services in ECS and beanstalk interact with all sorts of AWS Services and also amongst themselves to do work. Kinesis SQS and SNS and Lambda functions drive a large part of our event driven architecture. Of course we wanted this transition to be seamless, and the breadth of systems we were working with was a bit daunting, especially when the id that needed to change was tangled throughout our systems. We needed a quick and safe way to cut over to the new ID structure.

### Problem 1
Our system typically relies on an upstream system to seed most of it's data, in this case we couldn't rely on the upstream system to provide the data to us, so we needed a way to get that data ourself. Luckily we were able to call upstream API's to get us this data. We kept a cache of the most recent data in S3 and used that to hydrate our new tables.

### Problem 2
Remember how our primary key isn't actually primary any more, turns out you can't rely on on that kind of data to uniquely identify things. We had incorrect, orphaned, and in some cases *kind of correct* but not *actually* correct data (think weird and invisible characters from people copy and pasting). We also had operators directly interacting with the field in question. As a result of that, we had data inconsistencies across all of our data stores. When keys aren't unique and are text based it's a recipe for disaster.

### Problem 3
Event driven distributed systems tend to pass many messages around. systems which work based on passing messages are inherently decoupled, which means changes to contracts downstream mean upstream impact.  We needed to account for the producers sending the field, and also account for the consumers depending on it. All Producers Consumers had to be aware of the change. In many cases this meant the reads and writes were decoupled from table to table. So we needed the ability to make decisions about how to handle messages at run time.

### Problem 4
We made a goal to maintain next to **zero** downtime across all of our systems during the data migration. This is a tough problem in itself.

### Problem 5
Abbreviated Delivery time. sometimes that's out of our control. We ended up cutting into testing time which made the process less smooth than it could have been. Sometimes pressure is good, but it also leads to quick decisions which aren't always the best.

## The Solution
All of the data was in DynamoDB, unfortunately Dynamo doesn't perform well when you need to scan through a whole table. It takes forever. Luckily we are able to access data via a S3 Dump from our data pipeline backups. The data dump writes all of the data in new line separated data files, which means you can use the (plain-text) files and fix the data locally, then run a restore from that data into a new table. We stood new tables up, then restored into them.

With data seeded in the new tables from our point in time restore, we still needed a way to keep the data in sync with the live application. Luckily Dynamo Streams are a perfect use case for this type of work. We recycled the logic from the restore modification script, and put it into a lambda function triggered off of a Dynamo Stream. Best of all this method makes for easy clean up, once we cut over to the new tables, we were able get rid of the lambda with zero invasive code changes. 


Now that we have two Dynamo tables, one with the old ID and one with the new ID we needed a way for all of the services which depend on it to switch at runtime, without redeploying. To do this we took a simple approach with a simple Dynamo Table. The Table only has a flag, either `Feature ON` or `Feature OFF`. All services which needed access to the tables would have to delegate to the feature table. The downside of this is that every access to the old table would now require 2 reads, for us this wasn't that big of a deal. We did consider a solution with AppSync push updates but opted to go with the quick and  relatively inelegant 2 read solution.

This gave us the ability to coordinate a switch to the new table for all services in unison. No redeploys were necessary to switch over and as long as we ensured every read/write was governed by the flag,the switch over was smooth. Switching back and forth between the new and old tables was as simple as flipping the flag. We didn't need to cut back, but doing so would have been as easy as the cut over.


## What went well
Overall this exercise was a success, and so, I'll talk about what went well.

1. Our data pipeline strategy worked flawlessly. I was very happy with our ability to utilize our backups, along with Dynamo Streams to get the correct data in place. other solutions would have involved scanning the Dynamo Table, and that is painfully slow. Sometimes it's easy to complain about Dynamo, but Dynamo Streams are really cool. This provided us an out of the box solution to get our data where it needed to be; as it turns out this isn't the only place we're using Dynamo Streams, we're using them everywhere, they really are a cool feature. Further, if you have data backups sitting in s3, that can be great test data. don't let it go to waste, find a way to utilize it, experiment with new ways of storing it. We've been exploring graph data bases lately and having tons of data has allowed us to explore our options with real datasets very quickly.

2. Using Dynamo Streams and Lambda functions to transform the bad keys, and write the records to the new table worked flawlessly as well. Sometimes with Dynamo you can feel locked into the HashKey RangeKey's you've chosen. Requirements change over time and it can be easy to outgrow simple look-ups off of a primary key. Consider using Dynamo Streams to get your data into a new table, or even database; and explore if you can improve the way you query or reason about your data.

3. Our ad-hoc feature toggle table also was a success and is something we'd like to expand upon. Performance sensitive use-cases probably can't stand 2 round trips to Dynamo for every interaction with a table. One area for improvement here is AppSync (AWS's hosted GraphQL) for apps to get push notifications on changes to the Dynamo Table. More exploration should be done on this front for sure, especially given the power of flipping a configuration for 10-20+ apps in mere seconds.

## What went rocky

1. Our abbreviated timeline meant we were forced to push all of these changes, sooner than later. As a result full integration testing in our development environment didn't happen. We were able to test 90% of the critical path before pushing to prod. But, unfortunately the last 10% bit us and resulted in a couple late nights.

2. This effort touched every team within our organization, meaning everyone needed to deliver. Turns out some data was missed in the migration upstream, meaning we felt the effects. The only thing that can help here is a longer time frame and more integration testing.

The TLDR;

We...

- Created a new table
- Seeded the new table with data we fixed
- Set up a Dynamo Stream from `old -> new` so that we continued to keep the new table up to date
- Set up a configuration table so that all apps could switch back and forth at run time in unison
- Cut over
- Fix bugs
- profit
