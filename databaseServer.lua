#!/usr/bin/env tarantool

box.cfg {
    listen = 3301;
    log_format='json';
    log='logFile.txt';
 }

log = require('log')
json = require('json')
fun = require('fun')

box.once(
    "init", function()
        box.schema.space.create('storage')
        box.space.storage:format({
            { name = 'id', type = 'integer' },
            { name = 'name', type = 'string' },
            { name = 'category', type = 'string'}
        })
        box.space.storage:create_index(
            'primary',
            {
                type = 'hash',
                parts = {'id'}
            })
    end)

local function keyGen(path)
    local key = string.reverse(path)
    local i = string.find(key, '/')
    -- if path contains / at the end
    if (i == 1) then
        i = string.find(key, '/', 2)
        key = string.sub(key, 2, i - 1)
    else
        key = string.sub(key, 1, i - 1)
    end
    key = string.reverse(key)
    return key
end

local function get_products(req)
    local method = req.method
    log.info('Method: '..method)
    local all_products = fun.totable(box.space.storage:pairs())
    return { status = 200, body = json.encode(all_products) }
end

local function get_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = tonumber(keyGen(req.path))
    local product = box.space.storage:get(id)
    if product == nil then
        return { status = 400 }
    end
    return { status = 200, body = json.encode({name = product.name, id = id, category = product.category}) }
end

local function add_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = 0
    local all_products = fun.totable(box.space.storage:pairs())
    if #all_products > 0 then
        for _, tuple in box.space.storage.index.primary:pairs(nil, {iterator = box.index.ALL}) do
            if tuple.id > id then
                id = tuple.id
            end
        end
        id = id + 1
    end
    local name = req:json().name
    local category = req:json().category
    if name == nil or category == nil then
        log.info('Status: 400')
        log.error('Body isn\'t correct')
        return { status = 400 }
    end
    local exists = box.space.storage:count(id)
    -- return code 409 if key already exists, insert if not
    if exists > 0 then
        log.info('Status: 409')
        log.error('Product already exists in database')
        return { status = 409 }
    else
        box.space.storage:insert{id, name, category }
        log.info('Status: 200')
        log.info('Inserted: { id: '..id..', name: '..name..', category: '..category..' }')
        return { status = 200, body = json.encode(id)}
    end
end

local function edit_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = tonumber(req:json().id)
    local name = req:json().name
    local category = req:json().category
    -- return status 400 if body is incorrect
    if id == nil or (name == nil and category == nil) then
        log.info('Status: 400')
        log.error('Body isn\'t correct')
        return { status = 400 }
    end
    local exists = box.space.storage:count(id)
    -- update value if product exists, return status 404 if not
    if exists > 0 then
        local product = box.space.storage:get(id)
        local logStr = 'Update: { id: '..id
        if name == nil then
            box.space.storage:put{id, product.name, category}
            logStr = logStr..', category: '..category..' }'
        else
            if category == nil then
                box.space.storage:put{id, name, product.category }
                logStr = logStr..', name'..name..' }'
            else
                box.space.storage:put{id, name, category }
                logStr = logStr..', name'..name..', category'..category..' }'
            end
        end
        log.info('Status: 200')
        log.info(logStr)
        return { status = 200 }
    else
        log.info('Status: 404')
        log.error('Product doesn\'t exists in database')
        return { status = 404 }
    end
end

local function delete_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = tonumber(req:json().id)
    local exists = box.space.storage:count(id)
    -- update value if product exists, return status 404 if not
    if exists > 0 then
        log.info('Status: 200')
        log.info('Delete: { id: '..id..' }')
        box.space.storage:delete(id)
        return { status = 200 }
    else
        log.info('Status: 404')
        log.error('Product doesn\'t exists in database')
        return { status = 404 }
    end
end

local function delete_all(_)
    box.space.storage:drop()
    box.space._schema:delete('onceinit')
end


local server = require('http.server').new(nil, 8888)

server:route({ path = '/get_products', method = 'GET' }, get_products) -- for get all products
server:route({ path = '/add_product', method = 'POST' }, add_product) -- for add new product
server:route({ path = '/edit_product', method = 'POST' }, edit_product) -- for edit product
server:route({ path = '/delete_product', method = 'POST' }, delete_product) -- for delete product
server:route({ path = '/:id', method = 'GET' }, get_product) -- for get product
server:route({ path = '/delete_all'}, delete_all) -- not business logic


server:start()
