create table if not exists logs
(
    type       String,
    module     String,
    message    String,
    stacktrace String,
    time       String,
    timestamp  UInt64
) engine = MergeTree() partition by type order by (type, module)