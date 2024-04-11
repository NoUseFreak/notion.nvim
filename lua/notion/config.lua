local M = {
    _loaded = false,
    _api_key = '',
    _db_id = '',
    _user_id = '',
}

local function read_config_file()
    if M._loaded then
        return
    end

    local cfg_path = vim.fs.dirname(vim.fs.find({ '.notion.json' }, { upward = true })[1])
    if cfg_path == nil then
        -- print 'No .notion-cli.json found'
        return
    end

    local ok, result = pcall(vim.fn.json_decode, vim.fn.readfile(cfg_path .. '/.notion.json'))
    if not ok then
        print 'Error reading .notion.json'
        return
    end

    M._api_key = result.apiKey
    M._db_id = result.dbId
    M._user_id = result.userId

    M._loaded = true
end

M.get_api_key = function()
    read_config_file()

    return M._api_key
end

M.get_db_id = function()
    read_config_file()

    return M._db_id
end

M.get_user_id = function()
    read_config_file()

    return M._user_id
end

return M
