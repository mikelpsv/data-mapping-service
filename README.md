# Data mapping service



### Backup/Restore

Backup
```
$ sudo docker exec -t mapping_db pg_dump -U dev mapping > dump.sql
```

Restore
```
$ echo "CREATE DATABASE mapping" | sudo docker exec -i mapping_db psql -U dev postgres
$ cat dump.sql | sudo docker exec -i mapping_db psql -U dev mapping
```
