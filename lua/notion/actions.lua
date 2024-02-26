local actions = require 'telescope.actions'
local action_state = require 'telescope.actions.state'

--- actions
local function close_telescope_prompt(prompt_bufnr)
  local selection = action_state.get_selected_entry()
  actions.close(prompt_bufnr)
  return selection
end

local M = {}

M.issue_insert = function(prompt_bufnr)
  local selection = close_telescope_prompt(prompt_bufnr)
  if vim.api.nvim_buf_get_option(vim.api.nvim_get_current_buf(), 'modifiable') then
    vim.api.nvim_put({ '#' .. (selection.value or "") }, 'b', true, true)
  end
end

M.webview_issue = function(prompt_bufnr)
  -- Open issue in webview
  local selected = close_telescope_prompt(prompt_bufnr)
  if vim.fn.has('mac') == 1 then
    vim.fn.system('open ' .. selected._url)
  elseif vim.fn.has('unix') == 1 then
    vim.fn.system('xdg-open ' .. selected._url)
  else
    print('Unsupported OS')
  end

  print('Open issue #' .. selected.value .. ' - ' .. selected._url .. ' in webview')
end

return M
