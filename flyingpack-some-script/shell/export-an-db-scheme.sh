docker exec -it an_pg_db bash \
pg_dump --dbname=test_db -f out.sql -s -x -O --no-tablespaces \
exit \
docker cp an_pg_db:out.sql out.sql