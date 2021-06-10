#!/usr/bin/env tarantool

box.cfg {
    listen = 3301
}

s = box.schema.space.create('storage')
s:format({{name = 'key', type = 'string'}, {name = 'value', type = 'string'}})
s:create_index('primary', {type = 'HASH', parts = {'key'}})

