local ncli_config = require 'notion.config'

local M = {}

M.startup_checks = function()
  if vim.fn.executable 'notion.nvim' == 0 then
    print 'Notion CLI is not installed'
    return
  end
  if ncli_config.get_api_key() == '' then
    print 'Notion API Key is not set'
    return
  end
  if ncli_config.get_db_id() == '' then
    print 'Notion Database ID is not set'
    return
  end
  return true
end

return M
