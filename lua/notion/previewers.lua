local previewers = require 'telescope.previewers'
local ncli_config = require 'notion.config'

local M = {}

M.issue_previewer = previewers.new_termopen_previewer {
    title = 'Issue Preview',
    get_command = function(entry)
        local tmp_table = vim.split(entry.value, '\t')
        if vim.tbl_isempty(tmp_table) then
            return { 'echo', '' }
        end
        return {
            os.getenv 'SHELL',
            '-c',
            'notion.nvim db-issue-detail --db-id=' .. ncli_config.get_db_id() .. ' ' .. tmp_table[1] .. ' --render-content | less -RS',
        }
        -- return { 'notion.nvim', 'db-issue-detail', '--db-id=' .. ncli_config.get_db_id(), tmp_table[1], '--render-content', '|', 'less', '-RS' }
    end,
}

return M
