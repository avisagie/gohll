gohll
=====

A simple implementation of HyperLogLog in plain [Go](golang.org) for 
32-bit hash functions.

HyperLogLog lets you work out an approximate count of the number of
unique items in a large set of items, using a very small amount
memory. "Small" means not keeping a physical set of all the items
you've seen stream past. A 32-bit hash function on the items you count
and about 5000 bits of memory (8192 bits for this Go implementation)
should suffice to count up to 1e7 unique items to within about 2%.

E.g. you have a log of IP addresses visiting your extremely busy
website... You might want to know how many unique IPs your visitors
hail from, even though (obviously) everyone visits your site ten times
a day. After that you want to know how many uniques there were
everyday, and you want to the freedom to be able to aggregate the
unique IPs for the week, or month, or last 3 days. HyperLogLog
counters let you take the union of the 7 sets for the week, and still
get an approximate count out at the end that collapses the duplicates,
all for those 8000 bits per count.

For inspiration read [this](http://blog.aggregateknowledge.com/2012/10/25/sketch-of-the-day-hyperloglog-cornerstone-of-a-big-data-infrastructure/)
and [this](https://github.com/aggregateknowledge/postgresql-hll) and [the original paper](http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf).

Installation
------------

    go get github.com/avisagie/gohll
