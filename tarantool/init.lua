box.cfg{
    listen = "3301",
    work_dir = "/var/lib/tarantool",
    memtx_dir = "/var/lib/tarantool",
}

box.schema.space.create('kv', {if_not_exists = true})
box.space.kv:format({
    {name = 'key', type = 'string'},
    {name = 'value', type = '*'}
})
box.space.kv:create_index('primary', {parts = {'key'}, if_not_exists = true})


local user = os.getenv("TARANTOOL_USER")
local password = os.getenv("TARANTOOL_PASSWORD")

if user and password then
    if not box.schema.user.exists(user) then
        box.schema.user.create(user, {password = password})
        box.schema.user.grant(user, 'read,write,execute', 'space', 'kv')
    end
end
