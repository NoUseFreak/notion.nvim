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
        open = ncli_builtins.issues,
        all = ncli_builtins.issues_all,
        owned = ncli_builtins.issues_owned,
    },
    health = function()
        local health = vim.health or require 'health'
        local ok = health.ok or health.report_ok
        local error = health.error or health.report_error

        local good = true
        if not pcall(require, 'notion.builtin') then
            good = false
            error 'Notion.nvim is not installed'
        end

        if vim.fn.executable 'notion.nvim' == 0 then
            good = false
            error 'Notion CLI is not installed'
        end

        if good then
            ok 'Notion.nvim is healhty'
        end
    end,
}
