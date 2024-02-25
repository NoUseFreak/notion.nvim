local pickers = require 'telescope.pickers'
local actions = require 'telescope.actions'
local ncli_actions = require 'notion.actions'
local ncli_config = require 'notion.config'
local ncli_previewers = require 'notion.previewers'
local Job = require 'plenary.job'
local startup_checks = require('notion.checks').startup_checks

local M = {}

local exec_cmd = function(cmd, cwd)
    local stderr = {}
    local stdout, ret = Job:new({
        command = 'notion.nvim',
        args = cmd,
        cwd = cwd,
        env = {
            'NOTION_INTEGRATION_TOKEN=' .. ncli_config.get_api_key(),
            'PATH=' .. os.getenv 'PATH',
        },
        on_stderr = function(_, data)
            table.insert(stderr, data)
        end,
    }):sync(20000)

    return stdout, ret, stderr
end

local notion_api = function(action, args)
    local json, ret, err = exec_cmd { action, '--db-id=' .. ncli_config.get_db_id(), args }
    if ret ~= 0 then
        print('Error: ' .. table.concat(err, '\n'))
        return {}
    end
    local results = {}
    for _, v in ipairs(json) do
        table.insert(results, vim.json.decode(v))
    end
    return results
end

local notion_entity_maker = function(issue)
    return {
        value = issue.id,
        ordinal = issue.id .. ' ' .. issue.title .. ' ' .. table.concat(issue.assignees or {}, ' '),
        display = issue.id .. string.rep(' ', 10 - #issue.id) .. ' ' .. issue.title,
    }
end

local function showLoading(msg, fetch_fn, complete_fn)
    local row = math.floor((vim.o.lines - 5) / 2)
    local width = math.floor(vim.o.columns / 1.5)
    local col = math.floor((vim.o.columns - width) / 2)
    for _ = 1, (width - #msg) / 2, 1 do
        msg = ' ' .. msg
    end
    local prompt_win, prompt_opts = require('plenary.popup').create(msg, {
        border = {},
        borderchars = require('telescope.config').values.borderchars,
        height = 5,
        col = col,
        line = row,
        width = width,
    })
    vim.api.nvim_win_set_option(prompt_win, 'winhl', 'Normal:TelescopeNormal')
    vim.api.nvim_win_set_option(prompt_win, 'winblend', 0)
    local prompt_border_win = prompt_opts.border and prompt_opts.border.win_id
    if prompt_border_win then
        vim.api.nvim_win_set_option(prompt_border_win, 'winhl', 'Normal:TelescopePromptBorder')
    end
    vim.defer_fn(
        vim.schedule_wrap(function()
            local results = fetch_fn()
            if not pcall(vim.api.nvim_win_close, prompt_win, true) then
                print('Unable to close window: ', 'notion', '/', prompt_win)
            end
            complete_fn(results)
        end),
        10
    )
end

M.issue_dynamic = function(opts)
    opts = opts or {}

    if not startup_checks() then
        return
    end

    pickers
        .new({
            opts,
        }, {
            prompt_title = 'Notion Issues',
            debounce = 50,
            finder = require('telescope.finders').new_dynamic {
                fn = function(prompt)
                    if prompt == '' or #prompt < 3 then
                        return {}
                    end
                    local response = notion_api('db-issue', prompt)
                    local result = {}
                    for _, line in ipairs(response) do
                        for _, issue in ipairs(line) do
                            table.insert(result, issue)
                        end
                    end
                    return result
                end,
                entry_maker = notion_entity_maker,
            },
            sorter = require('telescope.config').values.generic_sorter {},
            previewer = ncli_previewers.issue_previewer,
            attach_mappings = function(_, map)
                actions.select_default:replace(ncli_actions.issue_insert)
                map('i', '<C-x>', ncli_actions.webview_issue)
                map('n', '<C-x>', ncli_actions.webview_issue)

                return true
            end,
        })
        :find()
end

M.issue_static = function(opts)
    opts = opts or {}

    if not startup_checks() then
        return
    end

    showLoading('Loading Issues', function()
        local response = notion_api('db-issue', '')
        local result = {}
        for _, line in ipairs(response) do
            for _, issue in ipairs(line) do
                table.insert(result, issue)
            end
        end
        return result
    end, function(result)
        pickers
            .new(opts, {
                prompt_title = 'Notion Issues',
                debounce = 50,
                finder = require('telescope.finders').new_table {
                    results = result,
                    entry_maker = notion_entity_maker,
                },
                sorter = require('telescope.config').values.generic_sorter {},
                previewer = ncli_previewers.issue_previewer,
                attach_mappings = function(_, map)
                    actions.select_default:replace(ncli_actions.issue_insert)
                    map('i', '<C-x>', ncli_actions.webview_issue)
                    map('n', '<C-x>', ncli_actions.webview_issue)

                    return true
                end,
            })
            :find()
    end)
end

return M
