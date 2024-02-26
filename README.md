# Notion.nvim

![Work In Progress](https://img.shields.io/badge/Work%20In%20Progress-orange?style=for-the-badge)
[![NeoVim](https://img.shields.io/badge/NeoVim-green.svg?style=for-the-badge&logo=neovim&logoColor=white)](https://neovim.io)
[![Lua](https://img.shields.io/badge/Lua-darkblue.svg?style=for-the-badge&logo=lua&logoColor=white)](http://www.lua.org)
[![Golang](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](http://www.go.dev)

Notion.nvim is a Neovim plugin that allows you to interact with Notion.so from Neovim.

Is is developed to simplfy adding issue numbers to git commits.

Currently it only supports searching for issues in a database and inserting the issue number in the current buffer.
Searching for issues is done using [telescope.nvim](https://github.com/nvim-telescope/telescope.nvim).


## Installation


### Requirements

- [Neovim](https://neovim.io/) 0.9.4 or later
- [plenary.nvim](https://nvim-lua/plenary.nvim)
- [telescope.nvim](https://github.com/nvim-telescope/telescope.nvim)
- notion.nvim binary, see below


#### Install the binary

The binary is not included in the repository, for now you need to install it using go.

```sh
go install github.com/nousefreak/notion.nvim/cmd/notion.nvim@latest
```


#### Using [Lazy](https://github.com/folke/lazy.nvim):

```lua
  {
    "NoUseFreak/notion.nvim",
    dependencies = {
      "nvim-lua/plenary.nvim",
    },
    config = function()
      require('telescope').load_extension 'notion'

      vim.keymap.set('n', '<leader>na', require('notion.builtin').issues_all, { noremap = true, silent = true })
      vim.keymap.set('n', '<leader>ni', require('notion.builtin').issues, { noremap = true, silent = true })
      vim.keymap.set('n', '<leader>no', require('notion.builtin').issues_owned, { noremap = true, silent = true })
    end,
  },
```

## Configuration

The plugins searches for a `.notion.json` file in the current directory and all parent directories.

This file should contain the following fields:

```json
{
  "apiKey": "your-api",
  "dbId": "your-db-id",
  "userID": "your-user-id"
}
```

### Getting an API key

The `apiKey` can be created in the Notion settings under [Integrations](https://www.notion.so/my-integrations).

### Getting the database id

Open the database you want to use and copy the id from the URL. It should look something like `https://www.notion.so/your-db-id`.

Note you may have to grant the intergration access to the database. This can be done by the hamburger menu in the top right corner of the database view.
In the menu under Connections you can `Connect to` the integration.

### Getting the user id

Not sure if a cleaner way exists. But you can get the user id by inspecting the local storage of the Notion web app.

```
JSON.parse(localStorage.getItem("gist.web.userToken")).value
```

## Keybindings

In the list of issues you can use `CTRL-X` to open the issue in the browser.

