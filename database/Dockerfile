FROM tarantool/tarantool:1.10.2
COPY app.lua /opt/tarantool
CMD ["tarantool", "/opt/tarantool/app.lua"]
# docker build -t tarantooldb .
# docker run --rm -d -p 3301:3301 --name=db tarantooldb
