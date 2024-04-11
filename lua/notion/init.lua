local M = {}

M.config = {
    prefix = '#',
}

M.setup = function(config)
    M.config = vim.tbl_extend('force', M.config, config or {})
end

return M
