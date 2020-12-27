-- MyNet.lua

local Dim = require 'Dim'
local Grid = require 'Grid'
local GameState = require 'GameState'

local physics = require 'physics'
physics.start()
physics.setGravity(0, 0)  -- 9.8
print(physics.engineVersion)

local composer = require('composer')
local scene = composer.newScene()
local widget = require('widget')

widget.setTheme('widget_theme_android_holo_dark')

local asteroidsTable = {}

-- local scrollView = nil
local gridGroup = nil
local backGroup = nil

local astroLoopTimer = nil

local grid = nil

local function createAsteroid2(x, y, color)
  local newAsteroid = display.newCircle(backGroup, x, y, math.random(6))
  table.insert(asteroidsTable, newAsteroid)
  physics.addBody(newAsteroid, 'dynamic', { density=0.3, radius=10, bounce=0.9 } )
  newAsteroid:setLinearVelocity( math.random( -25,25 ), math.random( -25,25 ) )
  newAsteroid:setFillColor(unpack(color))
end

local function astroLoop()

  grid:iterator(function(c)
    if c.bitCount == 1 then
      createAsteroid2(c.center.x, c.center.y, c.color)
    end
  end)

  -- Remove asteroids which have drifted off screen
  for i = #asteroidsTable, 1, -1 do
    local thisAsteroid = asteroidsTable[i]

    if ( thisAsteroid.x < 0 or
       thisAsteroid.x > display.contentWidth or
       thisAsteroid.y < 0 or
       thisAsteroid.y > display.contentHeight )
    then
      display.remove( thisAsteroid )
      table.remove( asteroidsTable, i )
    end
  end
end

function scene:create(event)
  local sceneGroup = self.view

  gridGroup = display.newGroup()
  sceneGroup:insert(gridGroup)

  backGroup = display.newGroup()
  sceneGroup:insert(backGroup)

--[[
  scrollView = widget.newScrollView({
    top = 0,
    left = 0,
    width = display.actualContentWidth,
    height = display.actualContentHeight,
    -- isBounceEnabled = false,
    backgroundColor = {0,0,0},
  })
  -- scrollView.x = 0
  -- scrollView.y = 0
  -- scrollView.anchorX = 0
  -- scrollView.anchorY = 0
  sceneGroup:insert(scrollView)
]]
  if system.getInfo('platform') == 'win32' then
    _G.dimensions = Dim:new(100)
  else
    _G.dimensions = Dim:new(200)
  end

  -- for debugging the gaps between cells problem
  -- display.setDefault('background', 0.5,0.5,0.5)

  grid = Grid:new(gridGroup, _G.dimensions.numX, _G.dimensions.numY)
  grid.gameState = GameState:new()
  grid:newLevel()

end

function scene:show(event)
  local sceneGroup = self.view
  local phase = event.phase

  if phase == 'will' then
    -- Code here runs when the scene is still off screen (but is about to come on screen)
    astroLoop()
  elseif phase == 'did' then
    -- Code here runs when the scene is entirely on screen
    astroLoopTimer = timer.performWithDelay(10000, astroLoop, 0)
  end
end

function scene:hide(event)
  local sceneGroup = self.view
  local phase = event.phase

  if phase == 'will' then
    -- Code here runs when the scene is on screen (but is about to go off screen)
    if astroLoopTimer then timer.cancel(astroLoopTimer) end
  elseif phase == 'did' then
    -- Code here runs immediately after the scene goes entirely off screen
    composer.removeScene('Tetra')
  end
end

function scene:destroy(event)
  local sceneGroup = self.view

  grid:destroy()

  -- Code here runs prior to the removal of scene's view
end

-- -----------------------------------------------------------------------------------
-- Scene event function listeners
-- -----------------------------------------------------------------------------------
scene:addEventListener('create', scene)
scene:addEventListener('show', scene)
scene:addEventListener('hide', scene)
scene:addEventListener('destroy', scene)
-- -----------------------------------------------------------------------------------

return scene
