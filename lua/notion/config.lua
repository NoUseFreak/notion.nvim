local M = {
  _loaded = false,
  _api_key = '',
  _db_id = '',
}

local function read_config_file()
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
  M._loaded = true
end

M.get_api_key = function()
  if not M._loaded then
    read_config_file()
  end

  return M._api_key
end

M.get_db_id = function()
  if not M._loaded then
    read_config_file()
  end

  return M._db_id
end

return M
