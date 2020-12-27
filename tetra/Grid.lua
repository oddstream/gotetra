-- Grid (of cells) class

local composer = require('composer')

local Cell = require 'Cell'
local Util = require 'Util'

local Grid = {
  -- prototype object
  group = nil,
  cells = nil,    -- array of Cell objects
  width = nil,      -- number of columns
  height = nil,      -- number of rows

  tapSound = nil,
  sectionSound = nil,
  lockedSound = nil,

  gameState = nil,
  levelText = nil,
}

function Grid:new(group, width, height)
  local o = {}
  self.__index = self
  setmetatable(o, self)

  o.group = group

  o.cells = {}
  o.width = width
  o.height = height

  for y = 1, height do
    for x = 1, width do
      local c = Cell:new(o, x, y)
      table.insert(o.cells, c) -- push
    end
  end

  o:linkCells2()

  if system.getInfo('environment') ~= 'simulator' then
    o.tapSound = audio.loadSound('assets/sound56.wav')
    o.sectionSound = audio.loadSound('assets/sound63.wav')
    o.lockedSound = audio.loadSound('assets/sound61.wav')
    o.completeSound = audio.loadSound('assets/complete.wav')
  end

  o.levelText = display.newText({
    parent=group,
    text='',
    x=display.contentCenterX,
    y=display.contentCenterY,
    font=native.systemFontBold,
    fontSize=512})
  o.levelText:setFillColor(0,0,0)
  o.levelText.alpha = 0.1

  return o
end

function Grid:reset()
  -- clear out the Cells
  self:iterator(function(c)
    c:reset()
  end)

  do
    local last_using = composer.getVariable('last_using')
    if not last_using then
      last_using = 0
    end
    local before = collectgarbage('count')
    collectgarbage('collect')
    local after = collectgarbage('count')
    print('collected', math.floor(before - after), 'KBytes, using', math.floor(after), 'KBytes', 'leaked', after-last_using)
    composer.setVariable('last_using', after)
  end

  self:newLevel()
end

function Grid:newLevel()

  self.colors, self.backgroundColor = Util.chooseColors()
  display.setDefault("background", unpack(self.backgroundColor))

  self:placeCoins()
  self:colorCoins()
  self:jumbleCoins()
  self:createGraphics()

  self.levelText.text = tostring(self.gameState.level)

  self:fadeIn()
end

function Grid:advanceLevel()
  assert(self.gameState)
  assert(self.gameState.level)
  self.gameState.level = self.gameState.level + 1
  self.levelText.text = tostring(self.gameState.level)
  self.gameState:write()
end

function Grid:sound(type)
  if type == 'tap' then
    if self.tapSound then audio.play(self.tapSound) end
  elseif type == 'section' then
    if self.sectionSound then audio.play(self.sectionSound) end
  elseif type == 'locked' then
    if self.lockedSound then audio.play(self.lockedSound) end
  elseif type == 'complete' then
    if self.lockedSound then audio.play(self.completeSound) end
  end
end

function Grid:linkCells2()
  for _,c in ipairs(self.cells) do
    c.n = self:findCell(c.x, c.y - 1)
    c.e = self:findCell(c.x + 1, c.y)
    c.s = self:findCell(c.x, c.y + 1)
    c.w = self:findCell(c.x - 1, c.y)
  end
end

function Grid:iterator(fn)
  for _,c in ipairs(self.cells) do
    fn(c)
  end
end

function Grid:findCell(x,y)
  for _,c in ipairs(self.cells) do
    if c.x == x and c.y == y then
      return c
    end
  end
  -- print('*** cannot find cell', x, y)
  return nil
end

function Grid:randomCell()
  return self.cells[math.random(#self.cells)]
end

function Grid:createGraphics()
  self:iterator(function(c) c:createGraphics(0.1) end)
end

function Grid:placeCoins()
  self:iterator(function(c) c:placeCoin() end)
  self:iterator(function(c) c:calcHammingWeight() end)
end

function Grid:colorCoins()
  local nColor = 1
  local section = 1
  local c = table.find(self.cells, function(d) return d.coins ~= 0 and d.color == nil end)
  while c do
    c:colorConnected(self.colors[nColor], section)
    nColor = nColor + 1
    if nColor > #self.colors then
      nColor = 1
    end
    section = section + 1
    c = table.find(self.cells, function(d) return d.coins ~= 0 and d.color == nil end)
  end
end

function Grid:jumbleCoins()
  self:iterator( function(c) c:jumbleCoin() end )
end

function Grid:isComplete()
  return table.find(self.cells, function(c) return c.coins ~= 0 end) == nil
end

function Grid:isSectionComplete(section)
  local arr = table.filter(self.cells, function(c) return c.section == section end)
  for n = 1, #arr do
    if not arr[n]:isComplete(section) then
      return false
    end
  end
  for n = 1, #arr do
    arr[n].section = 0  -- lock cell from moving
  end
  return true
end

function Grid:removeSection(section)
  -- print('remove section', section)
  self:iterator( function(c)
    if c.section == section then
      c:fadeOut()
      timer.performWithDelay(1000, function() c:reset() end, 1)
    end
  end )
end

function Grid:fadeIn()
  self:iterator( function(c) c:fadeIn() end )
end

-- function Grid:fadeOut()
--   self:iterator( function(c) c:fadeOut() end )
-- end

function Grid:destroy()
  local nStopped = audio.stop()  -- stop all channels
  print('audio stop', nStopped)

  if self.sectionSound then
    audio.dispose(self.sectionSound)
    self.sectionSound = nil
  end
  if self.tapSound then
    audio.dispose(self.tapSound)
    self.tapSound = nil
  end
  if self.lockedSound then
    audio.dispose(self.lockedSound)
    self.lockedSound = nil
  end
end

return Grid