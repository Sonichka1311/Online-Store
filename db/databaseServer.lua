#!/usr/bin/env tarantool

-- config --------------------------------------------------------------------------------------------------------------
box.cfg {
    listen = 3302;
    log_format='json';
    log='logs/logFile.txt';
    memtx_dir='logs';
    wal_dir='logs';
 }

log = require('log')
json = require('json')
fun = require('fun')

-- init ----------------------------------------------------------------------------------------------------------------
box.once(
    "init", function()
        box.schema.space.create('products')
        box.space.products:format({
            { name = 'id', type = 'integer' },
            { name = 'name', type = 'string' },
            { name = 'category', type = 'string'},
        })
        box.space.products:create_index(
            'primary',
            {
                type = 'hash',
                parts = {'id'}
            }
        )

        box.schema.space.create('users')
        box.space.users:format({
            { name = 'email', type = 'string' },
            { name = 'password', type = 'string' },
        })
        box.space.users:create_index(
            'primary',
            {
                type = 'hash',
                parts = {'email'}
            }
        )

        box.schema.space.create('sessions')
        box.space.sessions:format({
            { name = 'email', type = 'string' },
            { name = 'refresh_token', type = 'string' },
            { name = 'expire', type = 'string' }
        })
        box.space.sessions:create_index(
            'primary',
            {
                type = 'hash',
                parts = {'email'}
            }
        )
        box.space.sessions:create_index(
            'token',
            {
                type = 'hash',
                parts = {'refresh_token'}
            }
        )
    end)

-- additional functions ------------------------------------------------------------------------------------------------
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

-- products handlers ---------------------------------------------------------------------------------------------------
local function get_products(req)
    local method = req.method
    log.info('Method: '..method)
    local all_products = fun.totable(box.space.products:pairs())
    return { status = 200, body = json.encode(all_products) }
end

local function get_product(req)
    local method = req.method
    log.info('Method: '..method)
    --local id = req.
    local id = tonumber(keyGen(req.path))
    local exists = box.space.products:count(id)
    if exists > 0 then
        local product = box.space.products:get(id)
        return { status = 200, body = json.encode({name = product.name, id = id, category = product.category})}
    else
        return { status = 400 }
    end
end

local function add_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = 0
    local all_products = fun.totable(box.space.products:pairs())
    if #all_products > 0 then
        for _, tuple in box.space.products.index.primary:pairs(nil, {iterator = box.index.ALL}) do
            if tuple.id > id then
                id = tuple.id
            end
        end
        id = id + 1
    end
    local name = req:json().name
    local category = req:json().category
    if name == nil or category == nil then
        log.info('Status: 400, Message: Body isn\'t correct')
        return { status = 400 }
    end
    local exists = box.space.products:count(id)
    -- return code 409 if key already exists, insert if not
    if exists > 0 then
        log.info('Status: 409, Message: Product already exists in database')
        return { status = 409 }
    else
        box.space.products:insert{id, name, category}
        log.info('Status: 200')
        log.info('Inserted: { id: '..id..', name: '..name..', category: '..category..' }')
        local product = box.space.products:get(id)
        return { status = 200, body = json.encode({name = product.name, id = id, category = product.category}) }
    end
end

local function edit_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = tonumber(req:json().id)
    local name = req:json().name
    local category = req:json().category
    -- return status 400 if body is incorrect
    if id == nil then
        log.info('Status: 400, Message: Body isn\'t correct')
        return { status = 400 }
    end
    local exists = box.space.products:count(id)
    -- update value if product exists, return status 404 if not
    if exists > 0 then
        local product = box.space.products:get(id)
        local logStr = 'Update: { id: '..id
        if (name == nil or string.len(name) == 0) and (category == nil or string.len(category) == 0) then
            log.info('Status: 400, Message: Body isn\'t correct')
            return { status = 400 }
        end
        if name == nil or string.len(name) == 0 then
            box.space.products:put{id, product.name, category}
            logStr = logStr..', category: '..category..' }'
        else
            if category == nil or string.len(category) == 0 then
                box.space.products:put{id, name, product.category }
                logStr = logStr..', name'..name..' }'
            else
                box.space.products:put{id, name, category }
                logStr = logStr..', name'..name..', category'..category..' }'
            end
        end
        log.info('Status: 200')
        log.info(logStr)
        local product = box.space.products:get(id)
        return { status = 200, body = json.encode({name = product.name, id = id, category = product.category}) }
    else
        log.info('Status: 404, Message: Product doesn\'t exists in database')
        return { status = 404 }
    end
end

local function delete_product(req)
    local method = req.method
    log.info('Method: '..method)
    local id = tonumber(req:json().id)
    local exists = box.space.products:count(id)
    -- update value if product exists, return status 404 if not
    if exists > 0 then
        local product = box.space.products:get(id)
        log.info('Status: 200')
        log.info('Delete: { id: '..id..' }')
        box.space.products:delete(id)
        return { status = 200, body = json.encode({name = product.name, id = id, category = product.category}) }
    else
        log.info('Status: 404, Message: Product doesn\'t exists in database')
        return { status = 404 }
    end
end

local function delete_all(_)
    box.space.products:drop()
    box.space._schema:delete('onceinit')
end

-- users handlers ------------------------------------------------------------------------------------------------------
local function add_user(req)
    local method = req.method
    log.info('Method: '..method)
    local email = req:json().email
    local password = req:json().password
    if email == nil or password == nil then
        log.info('Status: 400, Message: Body isn\'t correct')
        return { status = 400 }
    end
    local exists = box.space.users:count(email)
    -- return code 409 if key already exists, insert if not
    if exists > 0 then
        log.info('Status: 409, Message: User already exists in database')
        return { status = 409 }
    else
        box.space.users:insert{email, password}
        log.info('Status: 200')
        log.info('Inserted: { email: '..email..', password: '..password)
        return { status = 200 }
    end
end

-- sessions handlers ---------------------------------------------------------------------------------------------------
local function add_session(req)
    local method = req.method
    log.info('Method: '..method)
    local email = req:json().email
    local password = req:json().password
    local refresh_token = req:json().refresh_token
    local expire = req:json().expire
    local exists = box.space.users:count(email)
    if exists > 0 then
        local user = box.space.users:get(email)
        if user.password == password then
            local session_exists = box.space.sessions:count(email)
            if session_exists > 0 then
                box.space.sessions:put{email, refresh_token, expire}
                log.info('Status: 200')
                log.info('Updated: { email: '..email..', refresh_token: '..refresh_token..", expire: "..expire..' }')
                return { status = 200 }
            else
                box.space.sessions:insert{email, refresh_token, expire}
                log.info('Status: 200')
                log.info('Inserted: { email: '..email..', refresh_token: '..refresh_token..", expire: "..expire..' }')
                return { status = 200 }
            end
        else
            log.info('Status: 400, Message: Bad password '..password..' for user '..email)
            return { status = 400 }
        end
    else
        log.info('Status: 404, Message: There is no user '..email..' in database')
        return { status = 404 }
    end
end

local function refresh_session(req)
    local method = req.method
    log.info('Method: '..method)
    local old_refresh_token = req:json().old_refresh_token
    local refresh_token = req:json().refresh_token
    local expire = req:json().expire
    local now = req:json().time
    local exists = box.space.sessions.index.token:count(old_refresh_token)
    if exists > 0 then
        local session = box.space.sessions.index.token:get(old_refresh_token)
        if session.expire > now then
            box.space.sessions:put{session.email, refresh_token, expire }
            log.info('Status: 200')
            log.info('Updated: { email: '..session.email..', refresh_token: '..refresh_token..", expire: "..expire..' }')
            local user = box.space.users:get(session.email)
            return { status = 200, body = json.encode({email = user.email, password = user.password}) }
        else
            log.info('Status: 403, Message: Expired refresh_token')
            return { status = 403 }
        end
    else
        log.info('Status: 404, Message: Invavid refresh_token')
        return { status = 404 }
    end
end

-- server's config -----------------------------------------------------------------------------------------------------
local server = require('http.server').new(nil, 3301)

server:route({ path = '/', method = 'GET' }, get_products) -- for get all products
server:route({ path = '/product/:id', method = 'GET' }, get_product) -- for get product
server:route({ path = '/product', method = 'POST' }, add_product) -- for add new product
server:route({ path = '/product', method = 'PUT' }, edit_product) -- for edit product
server:route({ path = '/product', method = 'DELETE' }, delete_product) -- for delete product
----------------------------
server:route({ path = '/delete_all'}, delete_all) -- not business logic, delete all from products database
----------------------------
server:route({ path = '/user', method = 'POST' }, add_user) -- for add new user
----------------------------
server:route({ path = '/session', method = 'POST' }, add_session) -- for add new session
server:route({ path = '/session', method = 'PUT' }, refresh_session) -- for refresh session
----------------------------

server:start()
