# Cassandra unique storage

The module is golang library for storing unique data in Cassandra database, by maintaining checksums of the objects that are stored.
It uses LWT to maintain consistency across all data.

## Blog post

For more details see the [blog post](https://npenkov.copm/2020-11-11-golang-cassandra-lwt/)