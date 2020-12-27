-- Cell class

local bit = require('plugin.bit')
local bezier = require('Bezier')

local SHOW_SQUARE = false

local Cell = {
  -- prototype object
  grid = nil,      -- grid we belong to
  x = nil,        -- column 1 .. width
  y = nil,        -- row 1 .. height
  center = nil,   -- point table, screen coords

  n, e, s, w = nil, nil, nil, nil,

  coins = 0,
  bitCount = 0,     -- hammingWeight
  color = nil,      -- e.g. {0,1,0}
  section = 0,      -- number of section (0 if locked)

  square = nil,     -- ShapeObject for outline
  grp = nil,        -- display group to put shapes in
  grpObjects = nil, -- list of ShapeObject for fading
}

function Cell:new(grid, x, y)
  local dim = _G.dimensions

  local o = {}
  self.__index = self
  setmetatable(o, self)

  o.grid = grid
  o.x = x
  o.y = y

  -- calculate where the screen coords center point will be
  o.center = {x=(x*dim.Q) - dim.Q + dim.Q50, y=(y*dim.Q) - dim.Q + dim.Q50}
  o.center.x = o.center.x + dim.marginX
  o.center.y = o.center.y + dim.marginY

  -- "These coordinates will automatically be re-centered about the center of the polygon."
  o.square = display.newPolygon(o.grid.group, o.center.x, o.center.y, dim.cellSquare)
  o.square:setFillColor(0,0,0) -- if alpha == 0, we don't get tap events
  if SHOW_SQUARE then
    o.square:setStrokeColor(0.1)
    o.square.strokeWidth = 4
  end

  o.square:addEventListener('tap', o) -- table listener

  return o
end

function Cell:reset()
  self.coins = 0
  self.bitCount = 0
  self.color = nil
  self.section = 0

  if self.grp then
    self.grp:removeSelf()
    self.grp = nil
  end
  if self.grpObjects then
    self.grpObjects = nil
  end
end

function Cell:calcHammingWeight()
  local function hammingWeight(coin)
    local w = 0
    for dir = 1, 4 do
      if bit.band(coin, 1) == 1 then
        w = w + 1
      end
      coin = bit.rshift(coin, 1)
    end
    return w
  end

  self.bitCount = hammingWeight(self.coins)
end

function Cell:shiftBits(num)
  local dim = _G.dimensions

  num = num or 1
  while num > 0 do
    if bit.band(self.coins, dim.WEST) == dim.WEST then
      -- high bit is set
      self.coins = bit.lshift(self.coins, 1)
      self.coins = bit.band(self.coins, dim.MASK)
      self.coins = bit.bor(self.coins, 1)
    else
      self.coins = bit.lshift(self.coins, 1)
    end
    num = num - 1
  end
end

function Cell:unshiftBits(num)
  local dim = _G.dimensions

  local function unshift(n)
    if bit.band(n, 1) == 1 then
      n = bit.rshift(n, 1)
      n = bit.bor(n, dim.WEST)
    else
      n = bit.rshift(n, 1)
    end
    return n
  end

  -- assert(unshift(1) == 8)
  -- assert(unshift(2) == 1)
  -- assert(unshift(4) == 2)
  -- assert(unshift(8) == 4)
  -- assert(unshift(15) == 15)

  num = num or 1
  while num > 0 do
    self.coins = unshift(self.coins)
    num = num - 1
  end
end

function Cell:isComplete(section)
  local dim = _G.dimensions

  if section and self.section ~= section then
    return false
  end
  for _, cd in ipairs(dim.cellData) do
    if bit.band(self.coins, cd.bit) == cd.bit then
      local cn = self[cd.link]
      if not cn then
        return false
      end
      if section and cn.section ~= section then
        return false
      end
      if bit.band(cn.coins, cd.oppBit) == 0 then
        return false
      end
    end
  end
  return true
end

function Cell:placeCoin()
  local dim = _G.dimensions

  for _,cd in ipairs(dim.cellData) do
    if math.random() < dim.PLACE_COIN_CHANCE then
      if self[cd.link] then
        self.coins = bit.bor(self.coins, cd.bit)
        self[cd.link].coins = bit.bor(self[cd.link].coins, cd.oppBit)
      end
    end
  end
end

function Cell:jumbleCoin()
  self:unshiftBits(math.random(5))
end

function Cell:colorConnected(color, section)
  local dim = _G.dimensions

  self.color = color
  self.section = section

  for _, cd in ipairs(dim.cellData) do
    if bit.band(self.coins, cd.bit) == cd.bit then
      local cn = self[cd.link]
      if cn and cn.coins ~= 0 and cn.color == nil then
        cn:colorConnected(color, section)
      end
    end
  end
end

function Cell:rotate(dir)

  local function afterRotate()
    self:createGraphics(1)
    if self.grid:isSectionComplete(self.section) then
      self.grid:sound('section')
      self.grid:removeSection(self.section)
    end
  end

  dir = dir or 'clockwise'

  if self.section == 0 then
    self.grid:sound('locked')
  elseif self.grp then
    self.grid:sound('tap')
    -- shift bits now (rather than in afterRotate) in case another tap happens while animating
    local degrees
    if dir == 'clockwise' then
      self:shiftBits()
      degrees = 90
    elseif dir == 'anticlockwise' then
      self:unshiftBits()
      degrees = -90
    end

    transition.to(self.grp, {
      time = 100,
      rotation = degrees,
      onComplete = afterRotate,
    })
  end
end

function Cell:tap(event)
  -- implement table listener for tap events
  if self.grid:isComplete() then
    -- print('completed', event.name, event.numTaps, self.x, self.y, self.coins, self.bitCount)
    self.grid:sound('complete')
    self.grid:reset()
    self.grid:advanceLevel()
  else
    self:rotate('clockwise')
  end
  return true
end

function Cell:createGraphics(scale)
  local dim = _G.dimensions

  scale = scale or 1

  assert(self.grid.backgroundColor)
  self.square:setFillColor(unpack(self.grid.backgroundColor))

  if 0 == self.coins then
    return
  end

  -- gotcha the 4th argument to set color function ~= the .alpha property
  -- blue={0,0,1}
  -- print(table.unpack(blue), 3)
  --> 0 6

  if self.grp then
    self.grp:removeSelf()
    self.grp = nil
    self.grpObjects = nil
  end

  self.grp = display.newGroup()
  -- center the group on the center of the hexagon, otherwise it's at 0,0
  self.grp.x = self.center.x
  self.grp.y = self.center.y
  self.grid.group:insert(self.grp)

  self.grpObjects = {}

  local sWidth = dim.Q / 7 --dim.Q10
  local capRadius = sWidth / 2

  if self.bitCount == 1 then

    local cd = table.find(dim.cellData, function(b) return self.coins == b.bit end)
    assert(cd)
    local line = display.newLine(self.grp,
      cd.c2eX / 2.5,
      cd.c2eY / 2.5,
      cd.c2eX,
      cd.c2eY)
    line.strokeWidth = sWidth
    line:setStrokeColor(unpack(self.color))
    line.alpha = 1
    line.xScale, line.yScale = scale, scale
    table.insert(self.grpObjects, line)

    local endcap = display.newCircle(self.grp, cd.c2eX, cd.c2eY, capRadius)
    endcap:setFillColor(unpack(self.color))
    endcap.alpha = 1
    endcap.xScale, endcap.yScale = scale, scale
    table.insert(self.grpObjects, endcap)

    local circle = display.newCircle(self.grp, 0, 0, dim.Q20)
    circle.strokeWidth = sWidth
    circle:setStrokeColor(unpack(self.color))
    circle:setFillColor(unpack(self.grid.backgroundColor))
    circle.alpha = 1
    circle.xScale, circle.yScale = scale, scale

    table.insert(self.grpObjects, circle)

  else
--[[
  with self.bitCount > 3
  three consective bits should produce same pattern (rotated) no matter where they occur in coins:
    000111
    001110
    011100
    111000
    110001 - ugly
    100011 - ugly
  hence self.bitCount > 2
]]
    -- make a list of edge coords we need to visit
    local arr = {}
    for _,cd in ipairs(dim.cellData) do
      if bit.band(self.coins, cd.bit) == cd.bit then
        table.insert(arr, {x=cd.c2eX, y=cd.c2eY})
      end
    end

    -- close path for better aesthetics
    if self.bitCount > 2 then
      table.insert(arr, arr[1])
    end

    for n = 1, #arr-1 do
    -- use (off-)center and (off-)center as control points
      -- local av = 1.8  -- make the three-edges circles round
      local av = 3  -- make the three-edges circles triangular
      local cp1 = {x=(arr[n].x)/av, y=(arr[n].y)/av}
      local cp2 = {x=(arr[n+1].x)/av, y=(arr[n+1].y)/av}
      local curve = bezier.new(
        arr[n].x, arr[n].y,
        cp1.x, cp1.y,
        cp2.x, cp2.y,
        arr[n+1].x, arr[n+1].y)
      local curveDisplayObject = curve:get()
      curveDisplayObject.strokeWidth = sWidth
      curveDisplayObject:setStrokeColor(unpack(self.color))
      curveDisplayObject.alpha = 1
      curveDisplayObject.xScale,curveDisplayObject.yScale = scale, scale
      self.grp:insert(curveDisplayObject)
      table.insert(self.grpObjects, curveDisplayObject)
    end

    for n = 1, #arr do
      local endcap = display.newCircle(self.grp, arr[n].x, arr[n].y, capRadius)
      endcap:setFillColor(unpack(self.color))
      endcap.alpha = 1
      endcap.xScale, endcap.yScale = scale, scale
      table.insert(self.grpObjects, endcap)
    end

  end
end

function Cell:fadeIn()
  if self.grpObjects then
    for _, o in ipairs(self.grpObjects) do
      transition.scaleTo(o, {xScale=1, yScale=1, time=500});
    end
  end
end

function Cell:fadeOut()
  if self.grpObjects then
    for _, o in ipairs(self.grpObjects) do
      transition.scaleTo(o, {xScale=0.1, yScale=0.1, time=500});
    end
  end
end

--[[
function Cell:destroy()
  self.grid = nil
  -- TODO
end
]]

return Cell