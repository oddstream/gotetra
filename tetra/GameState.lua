-- GameState.lua

local json = require('json')

local GameState = {
  level = 0,
}

local filePath = system.pathForFile('gameState.json', system.DocumentsDirectory)
print(filePath)

function GameState:new()
  local o = {}
  self.__index = GameState
  setmetatable(o, self)

  o.level = 0
  o:read()

  return o
end

function GameState:read()
  local file = io.open(filePath, 'r')
  if file then
    local contents = file:read('*a')
    io.close(file)
    print('state loaded from file', contents)
    local state = json.decode(contents)
    if state and state.level then
      self.level = state.level
    end
  end
end

function GameState:write()
  local file = io.open(filePath, 'w')
  if file then
    file:write(json.encode(self))
    print('state written to', filePath)
    io.close(file )
  end
end

return GameState
