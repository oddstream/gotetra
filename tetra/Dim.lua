-- Dim.lua

local Dim = {
  Q = nil,

  Q50 = nil,
  Q20 = nil,
  Q10 = nil,

  numX = nil,
  numY = nil,

  marginX = nil,
  marginY = nil,

  NORTH = 1,
  EAST = 2,
  SOUTH = 4,
  WEST = 8,

  MASK = 15,

  PLACE_COIN_CHANCE = 0.2,

  cellData = nil,
  cellSquare = nil,
}

function Dim:new(Q)
  local o = {}
  self.__index = self
  setmetatable(o, self)

  o.Q = Q

  o.Q50 = math.floor(Q/2)
  o.Q20 = math.floor(Q/5)
  o.Q10 = math.floor(Q/10)

  o.numX = math.floor(display.actualContentWidth / Q)
  o.numY = math.floor(display.actualContentHeight / Q)

  o.marginX = (display.actualContentWidth - (o.numX * Q)) / 2
  o.marginY = (display.actualContentHeight - (o.numY * Q)) / 2

  o.cellData = {
    { bit=o.NORTH,  oppBit=o.SOUTH,   link='n',  c2eX=0,      c2eY=-o.Q50,  },
    { bit=o.EAST,   oppBit=o.WEST,    link='e',  c2eX=o.Q50,  c2eY=0,       },
    { bit=o.SOUTH,  oppBit=o.NORTH,   link='s',  c2eX=0,      c2eY=o.Q50,   },
    { bit=o.WEST,   oppBit=o.EAST,    link='w',  c2eX=-o.Q50, c2eY=0,       },
  }

  o.cellSquare = {
    -o.Q50, -o.Q50, -- top left
     o.Q50, -o.Q50, -- top right
     o.Q50,  o.Q50, -- bottom right
    -o.Q50,  o.Q50, -- bottom left
  }

  return o
end

return Dim