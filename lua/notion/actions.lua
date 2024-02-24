local actions = require 'telescope.actions'
local action_state = require 'telescope.actions.state'

--- actions
local function close_telescope_prompt(prompt_bufnr)
  local selection = action_state.get_selected_entry()
  actions.close(prompt_bufnr)
  local tmp_table = vim.split(selection.value, '\t')
  if vim.tbl_isempty(tmp_table) then
    return
  end
  return tmp_table[1]
end

local M = {}

M.issue_insert = function(prompt_bufnr)
  local issue_number = close_telescope_prompt(prompt_bufnr)
  if vim.api.nvim_buf_get_option(vim.api.nvim_get_current_buf(), 'modifiable') then
    vim.api.nvim_put({ '#' .. issue_number }, 'b', true, true)
  end
end

M.webview_issue = function(prompt_bufnr)
  -- Open issue in webview
  local issue_number = close_telescope_prompt(prompt_bufnr)
  print('Open issue #' .. issue_number .. ' in webview')
end

return M
