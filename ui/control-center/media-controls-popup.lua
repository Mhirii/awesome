local awful_popup = require("awful.popup")
local beautiful = require("beautiful")
local gstring = require("gears.string")
local wibox = require("wibox")

local helpers = require("helpers")
local playerctl = require("signals.playerctl")

local media_image = require("ui.widgets.media.media-image")
local media_controls = require("ui.widgets.media.media-controls")
local player_icon = require("ui.widgets.media.player-icon")

local scrolling_text = require("ui.widgets.scrolling-text")
local slider = require("ui.widgets.slider")

local media_buttons = wibox.widget {
    media_controls.prev(16),
    media_controls.play(21),
    media_controls.next(16),
    spacing = dpi(2),
    layout = wibox.layout.fixed.horizontal,
    widget = wibox.container.background
}

local media_title = scrolling_text {
    title = "Title",
    font = "Roboto 12",
    speed = 32,
    step_function = wibox.container.scroll.step_functions.waiting_nonlinear_back_and_forth
}

local artist_name = scrolling_text {
    text = "Artist",
    font = "Roboto 11",
    speed = 32,
    step_function = wibox.container.scroll.step_functions.waiting_nonlinear_back_and_forth
}

local media_info = wibox.widget {
    media_title,
    artist_name,
    layout = wibox.layout.fixed.vertical
}

local progress = slider {
    max = 1,
    bar_bg_color = beautiful.accent .. "70",
    bar_color = beautiful.accent,
    handle_width = dpi(12),
    handle_color = beautiful.accent,
    handle_border_width = 0
}

local previous_value = 0
local internal_update = false

progress:connect_signal(
    "property::value", function(_, new_value)
        if internal_update and new_value ~= previous_value then
            playerctl:set_position(new_value)
            previous_value = new_value
        end
    end
)

local interval = wibox.widget {
    text = "-/-",
    font = beautiful.font_name .. "11",
    valign = "center",
    halign = "right",
    widget = wibox.widget.textbox
}

local function secs_to_min(secs)
    local mins = math.floor(secs / 60)
    local remain_secs = secs % 60

    return string.format("%.0f", mins) .. ":" .. string.format("%02d", remain_secs)
end

local last_length = 0
playerctl:connect_signal(
    "position", function(_, interval_sec, length_sec)
        if length_sec ~= last_length then
            progress.maximum = length_sec
            last_length = length_sec
        end
        interval:set_markup_silently(secs_to_min(interval_sec) .. " / " .. secs_to_min(length_sec))
        internal_update = true
        previous_value = interval_sec
        progress.value = interval_sec
    end
)

local cover = wibox.widget {
    {
        media_image(8),
        widget = wibox.container.place
    },
    {
        player_icon(24, "popup"),
        valign = "bottom",
        halign = "right",
        widget = wibox.container.place
    },
    horizontal_offset = dpi(8),
    layout = wibox.layout.stack
}

local body_container = {
    {
        {
            {
                media_info,
                progress,
                {
                    media_buttons,
                    nil,
                    interval,
                    layout = wibox.layout.flex.horizontal
                },
                spacing = dpi(8),
                forced_width = dpi(200),
                layout = wibox.layout.fixed.vertical
            },
            {
                cover,
                widget = wibox.container.place
            },
            spacing = dpi(16),
            layout = wibox.layout.fixed.horizontal
        },
        margins = dpi(12),
        widget = wibox.container.margin
    },
    bg = beautiful.black,
    widget = wibox.container.background
}

local media_controls_popup = awful_popup {
    type = "dock",
    bg = beautiful.black,
    maximum_width = dpi(beautiful.control_center_width),
    border_width = dpi(2),
    border_color = beautiful.focus,
    ontop = true,
    visible = false,
    shape = helpers.rrect(beautiful.border_radius),
    widget = body_container
}

playerctl:connect_signal(
    "metadata", function(_, title, artist, album_path)
        media_title.text.text = gstring.xml_unescape(title)
        artist_name.text.text = gstring.xml_unescape(artist)
    end
)

playerctl:connect_signal(
    "no_players", function(_)
        media_controls_popup.visible = false
    end
)

awesome.connect_signal(
    "media::dominantcolors", function(colors)
        local bg_color, fg_color, _ = table.unpack(colors)

        -- darkening the bg color to match the dark theming
        media_controls_popup.widget.bg = bg_color .. "D0"
        media_controls_popup.border_color = bg_color
        media_controls_popup.fg = fg_color

        progress.bar_color = fg_color .. "70"
        progress.bar_active_color = fg_color
        progress.handle_color = fg_color
    end
)

return media_controls_popup
