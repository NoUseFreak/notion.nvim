local ncli_builtins = require 'notion.builtin'

return require('telescope').register_extension {
  setup = function()
    require('telescope').setup {
      defaults = {
        set_env = {
          NOTION_INTEGRATION_TOKEN = require('notion.config').get_api_key(),
        },
      },
    }
  end,
  exports = {
    issues = ncli_builtins.issue_static,
    issues_dynamic = ncli_builtins.issue_dynamic,
  },
  health = function()
    local health = vim.health or require "health"
    local ok = health.ok or health.report_ok
    local error = health.error or health.report_error

    if not pcall(require, 'notion.builtin') then
      error('Notion.nvim is not installed')
    end

    if vim.fn.executable('notion.nvim') == 0 then
      error('Notion CLI is not installed')
    end

    ok('Notion.nvim is healhty')
  end,
}
