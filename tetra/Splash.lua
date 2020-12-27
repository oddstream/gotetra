-- Splash.lua

local composer = require('composer')
local scene = composer.newScene()

local tim = nil
local logo = nil
local destination = nil

local function gotoDestination(event)
  composer.gotoScene(destination, {effect='fade'})
  return true -- we handled tap event
end

function scene:create(event)
  local sceneGroup = self.view
end

function scene:show(event)
  local sceneGroup = self.view
  local phase = event.phase

  if phase == 'will' then
    -- Code here runs when the scene is still off screen (but is about to come on screen)
    local x = display.contentCenterX
    local y = display.contentCenterY
    local cir = display.newCircle(sceneGroup, x, y, 210)
    cir:setFillColor(0.2, 0.2, 0.2)

    y = display.contentCenterY - 120
    logo = display.newImage(sceneGroup, 'oddstream logo.png', system.ResourceDirectory, x, y)

    assert(logo:addEventListener('tap', gotoDestination))

    y = display.contentCenterY
    display.newText(sceneGroup, 'oddstream', x, y, native.systemFontBold, 72)

    y = display.contentCenterY + 100
    display.newText(sceneGroup, 'games', x, y, native.systemFontBold, 72)

  elseif phase == 'did' then
    destination = event.params.scene
    -- Code here runs when the scene is entirely on screen
    tim = timer.performWithDelay(2000, gotoDestination, 1)
  end
end

function scene:hide(event)
  local sceneGroup = self.view
  local phase = event.phase

  if phase == 'will' then
    -- Code here runs when the scene is on screen (but is about to go off screen)
  elseif phase == 'did' then
    -- Code here runs immediately after the scene goes entirely off screen
    composer.removeScene('Splash')
  end
end

function scene:destroy(event)
  local sceneGroup = self.view
  -- Code here runs prior to the removal of scene's view
  if tim then
    timer.cancel(tim)
    tim = nil
  end
  logo:removeEventListener('tap', gotoDestination)
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