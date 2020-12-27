-- Util.lua

local Util = {}
Util.__index = Util

--[[
  Linear interpolation. The idea is very simple, you have 2 values, and you want to “walk” between those values by a factor.
  If you pass a factor of 0, you are pointing to the beginning of the walk so the value is equal to start.
  If you pass a factor of 1, you are pointing to the end of the walk so the value is equal to end.
  Any factor between 0 and 1 will add a (1-factor) of start argument and a factor of end argument.
  (e.g with start 0 and end 10 with a factor 0.5 you will have a 5, so the half of the path)
]]
function Util.lerp(start, finish, factor)
  -- return start*(1-factor) + finish*factor
  -- Precise method, which guarantees v = v1 when t = 1.
  -- https://en.wikipedia.org/wiki/Linear_interpolation
  return (1 - factor) * start + factor * finish;
end

--[[
  The opposite of lerp. Instead of a range and a factor, we give a range and a value to find out the factor.
]]
function Util.normalize(start, finish, value)
  return (value - start) / (finish - start)
end

--[[
  converts a value from the scale [fromMin, fromMax] to a value from the scale[toMin, toMax].
  It’s just the normalize and lerp functions working together.
]]
function Util.mapValue(value, fromMin, fromMax, toMin, toMax)
  return Util.lerp(toMin, toMax, Util.normalize(fromMin, fromMax, value))
end

function Util.clamp(value, min, max)
  return math.min(math.max(value, min), max)
end

function Util.isPointInTriangle(px,py,ax,ay,bx,by,cx,cy)
  -- assert(type(ax)=='number')
  -- assert(type(ay)=='number')
  -- assert(type(bx)=='number')
  -- assert(type(by)=='number')
  -- assert(type(cx)=='number')
  -- assert(type(cy)=='number')
  local v0 = {cx-ax,cy-ay}
  local v1 = {bx-ax,by-ay}
  local v2 = {px-ax,py-ay}

  local dot00 = (v0[1]*v0[1]) + (v0[2]*v0[2]);
  local dot01 = (v0[1]*v1[1]) + (v0[2]*v1[2]);
  local dot02 = (v0[1]*v2[1]) + (v0[2]*v2[2]);
  local dot11 = (v1[1]*v1[1]) + (v1[2]*v1[2]);
  local dot12 = (v1[1]*v2[1]) + (v1[2]*v2[2]);

  local invDenom = 1 / (dot00 * dot11 - dot01 * dot01)

  local u = (dot11 * dot02 - dot01 * dot12) * invDenom
  local v = (dot00 * dot12 - dot01 * dot02) * invDenom

  return ((u >= 0) and (v >= 0) and (u + v < 1))
end

function Util.chooseColors()
  local function dec2Float(r,g,b)
    return {r/255,g/255,b/255}
  end

  -- https://en.wikipedia.org/wiki/Web_colors
  local colorsGreen = {
    dec2Float(  0, 100,   0), -- DarkGreen
    dec2Float( 85, 107,  47), -- DarkOliveGreen
    dec2Float(107, 142,  35), -- OliveDrab
    dec2Float(139,  69,  19), -- SaddleBrown
    dec2Float( 80,  80,   0), -- Olive
    dec2Float(154, 205,  50), -- YellowGreen
    dec2Float( 46, 139,  87), -- SeaGreen
    -- dec2Float(128,128,128),
  }
  local colorsPink = {
    dec2Float(255,192,203), -- Pink
    dec2Float(255,105,180), -- HotPink
    dec2Float(219,112,147), -- PaleVioletRed
    dec2Float(255, 20,147), -- DeepPink
    dec2Float(199, 21,133), -- MediumVioletRed
    dec2Float(238,130,238), -- Violet
  }
  local colorsBlue = {
    dec2Float( 25, 25,112), -- MidnightBlue
    dec2Float( 65,105,225), -- RoyalBlue
    dec2Float( 30,144,255), -- DodgerBlue
    dec2Float(135,206,250), -- LightSkyBlue
    dec2Float( 70,130,180), -- SteelBlue
  }
  local colorsOrange = {
    dec2Float(255,165,  0),
    dec2Float(255, 69,  0),
    dec2Float(255,127, 80),
    dec2Float(255,140,  0),
    dec2Float(255, 99, 71),
  }
  local colorsGray = {
    dec2Float(128,128,128),
    dec2Float(192,192,192),
    dec2Float(112,128,144),
    dec2Float( 50, 50, 50),
    dec2Float(220,220,220),
  }

  local colorsYellow = {
    dec2Float(189,183,107),  -- DarkKhaki
    dec2Float(240,230,140),  -- Khaki
    dec2Float(255,218,185),  -- PeachPuff
    dec2Float(255,228,181),  -- Moccasin
    dec2Float(255,239,213),  -- PapayaWhip
    dec2Float(255,250,205),  -- LemonChiffon
  }

  local colorsAll = {
    colorsGreen,
    colorsBlue,
    colorsOrange,
    colorsYellow,
    colorsGray,
    colorsPink,
  }

  local colors = colorsAll[math.random(#colorsAll)]

  local sum = {0,0,0}
  for _,row in ipairs(colors) do
    for i = 1,3 do
      sum[i] = sum[i] + row[i]
    end
  end
  local avg = {0,0,0}
  for i = 1, 3 do
    avg[i] = sum[i] / #colors
  end

  -- trace('average color', avg[1], avg[2], avg[3])

  return colors, {avg[1]*0.5, avg[2]*0.5, avg[3]*0.5}
end

return Util
